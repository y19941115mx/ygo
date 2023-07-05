package command

import (
	"context"
	"errors"
	"fmt"

	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/y19941115mx/ygo/framework"
	"github.com/y19941115mx/ygo/framework/cobra"
	"github.com/y19941115mx/ygo/framework/contract"
	"github.com/y19941115mx/ygo/framework/provider/ssh"
	"github.com/y19941115mx/ygo/framework/util"
)

// initDeployCommand 为自动化部署的命令
func initDeployCommand() *cobra.Command {
	deployCommand.AddCommand(deployRollbackCommand)
	return deployCommand
}

// deployCommand 一级命令，显示帮助信息
var deployCommand = &cobra.Command{
	Use:   "deploy",
	Short: "部署项目的命令，支持回滚操作",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		// 创建部署文件夹
		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}

		// 编译后端到部署文件夹
		if err := deployBuildBackend(c, deployFolder); err != nil {
			return err
		}

		// 上传部署文件夹并执行对应的shell
		return deployUploadAction(deployFolder, container)
	},
}

// deployRollbackCommand 部署回滚
var deployRollbackCommand = &cobra.Command{
	Use:   "rollback",
	Short: "部署回滚",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		if len(args) != 1 {
			return errors.New("参数错误,请按照参数进行回滚 ./ygo deploy rollback [version]")
		}

		version := args[0]

		// 获取版本信息
		appService := container.MustMake(contract.AppKey).(contract.App)
		deployFolder := filepath.Join(appService.DeployFolder(), version)

		// 上传部署文件夹并执行对应的shell
		return deployUploadAction(deployFolder, container)
	},
}

func deployBuildBackend(c *cobra.Command, deployFolder string) error {
	container := c.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	appService := container.MustMake(contract.AppKey).(contract.App)
	envService := container.MustMake(contract.EnvKey).(contract.Env)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	env := envService.AppEnv()
	binFile := "ygo"

	if configService.GetString("deploy.backend.binFile") != "" {
		binFile = configService.GetString("deploy.backend.binFile")
	}

	// 编译前端
	path, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln("请在Path路径中先安装go")
	}

	// 组装命令
	deployBinFile := filepath.Join(deployFolder, binFile)
	cmd := exec.Command(path, "build", "-o", deployBinFile, "./")
	cmd.Env = os.Environ()
	// 设置GOOS和GOARCH
	if configService.GetString("deploy.backend.goos") != "" {
		cmd.Env = append(cmd.Env, "GOOS="+configService.GetString("deploy.backend.goos"))
	}
	if configService.GetString("deploy.backend.goarch") != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+configService.GetString("deploy.backend.goarch"))
	}

	// 执行命令
	ctx := context.Background()
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(ctx, "go build err", map[string]interface{}{
			"err": err,
			"out": string(out),
		})
		return err
	}
	logger.Info(ctx, "编译成功", nil)

	// 复制.env文件
	if util.Exists(filepath.Join(appService.BaseFolder(), ".env")) {
		if err := util.CopyFile(filepath.Join(appService.BaseFolder(), ".env"), filepath.Join(deployFolder, ".env")); err != nil {
			return err
		}
	}

	// 复制config文件
	deployConfigFolder := filepath.Join(deployFolder, "config", env)
	if !util.Exists(deployConfigFolder) {
		if err := os.MkdirAll(deployConfigFolder, os.ModePerm); err != nil {
			return err
		}
	}
	if err := util.CopyFolder(filepath.Join(appService.ConfigFolder(), env), deployConfigFolder); err != nil {
		return err
	}

	logger.Info(ctx, "build local ok", nil)
	return nil
}

// 上传部署文件夹，并且执行对应的前置和后置的shell
func deployUploadAction(deployFolder string, container framework.Container) error {
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	sshService := container.MustMake(contract.SSHKey).(contract.SSHService)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	// 遍历所有deploy的服务器
	deployNodes := configService.GetStringSlice("deploy.connections")
	if len(deployNodes) == 0 {
		return errors.New("deploy connections len is zero")
	}
	remoteFolder := configService.GetString("deploy.remote_folder")
	if remoteFolder == "" {
		return errors.New("remote folder is empty")
	}

	preActions := make([]string, 0, 1)
	postActions := make([]string, 0, 1)

	preActions = append(preActions, configService.GetStringSlice("deploy.backend.pre_action")...)
	postActions = append(postActions, configService.GetStringSlice("deploy.backend.post_action")...)

	// 对每个远端服务
	for _, node := range deployNodes {
		sshClient, err := sshService.GetClient(ssh.WithConfigPath(node))
		if err != nil {
			return err
		}
		client, err := sftp.NewClient(sshClient)
		if err != nil {
			return err
		}

		// 执行所有的前置命令
		for _, action := range preActions {
			// 创建session
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}
			logger.Info(context.Background(), "execute pre action start", map[string]interface{}{
				"cmd":        action,
				"connection": node,
			})
			// 执行命令，并且等待返回
			bts, err := session.CombinedOutput(action)
			if err != nil {
				session.Close()
				return err
			}
			session.Close()
			// 执行前置命令成功
			logger.Info(context.Background(), "execute pre action", map[string]interface{}{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
		}

		if err := uploadFolderToSFTP(container, deployFolder, remoteFolder, client); err != nil {
			logger.Info(context.Background(), "upload folder failed", map[string]interface{}{
				"err": err,
			})
			return err
		}
		logger.Info(context.Background(), "upload folder success", nil)

		for _, action := range postActions {
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}
			logger.Info(context.Background(), "execute post action start", map[string]interface{}{
				"cmd":        action,
				"connection": node,
			})
			bts, err := session.CombinedOutput(action)
			if err != nil {
				session.Close()
				return err
			}
			logger.Info(context.Background(), "execute post action finish", map[string]interface{}{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
			session.Close()
		}
	}
	return nil
}

// 上传部署文件夹
func uploadFolderToSFTP(container framework.Container, localFolder, remoteFolder string, client *sftp.Client) error {
	logger := container.MustMake(contract.LogKey).(contract.Log)
	// 遍历 deploy 文件夹
	return filepath.Walk(localFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 相对 deploy 文件夹的相对路径
		relPath := strings.Replace(path, localFolder, "", 1)
		if relPath == "" {
			return nil
		}
		// 如果是遍历到了一个目录
		if info.IsDir() {
			dirPath := toLinuxPath(filepath.Join(remoteFolder, relPath))
			logger.Info(context.Background(), "mkdir: "+dirPath, nil)

			// 创建这个目录
			return client.MkdirAll(dirPath)
		}

		// 打开本地的文件
		rf, readErr := os.Open(filepath.Join(localFolder, relPath))
		if readErr != nil {
			return errors.New("read file " + filepath.Join(localFolder, relPath) + " error:" + err.Error())
		}
		defer rf.Close()
		// 检查文件大小
		rfStat, _ := rf.Stat()
		// 打开/创建远端文件
		remoteFilePath := toLinuxPath(filepath.Join(remoteFolder, relPath))
		f, err := client.Create(remoteFilePath)
		if err != nil {
			return errors.New("create file " + remoteFilePath + " error:" + err.Error())
		}
		defer f.Close()

		// 大于2M的文件显示进度
		if rfStat.Size() > 2*1024*1024 {
			logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
				" to remote file: "+remoteFilePath+" start", nil)
			// 开启一个goroutine来不断计算进度
			go func(localFile, remoteFile string) {
				// 每10s计算一次
				ticker := time.NewTicker(2 * time.Second)
				for range ticker.C {
					// 获取远端文件信息
					remoteFileInfo, err := client.Stat(remoteFile)
					if err != nil {
						logger.Error(context.Background(), "stat error", map[string]interface{}{
							"err":         err,
							"remote_file": remoteFile,
						})
						continue
					}
					// 如果远端文件大小等于本地文件大小，说明已经结束了
					size := remoteFileInfo.Size()
					if size >= rfStat.Size() {
						break
					}
					// 计算进度并且打印进度
					percent := int(size * 100 / rfStat.Size())
					logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
						" to remote file: "+remoteFile+fmt.Sprintf(" %v%% %v/%v", percent, size, rfStat.Size()), nil)
				}
			}(filepath.Join(localFolder, relPath), remoteFilePath)
		}

		// 将本地文件并发读取到远端文件
		if _, err := f.ReadFromWithConcurrency(rf, 10); err != nil {
			return errors.New("Write file " + remoteFilePath + " error:" + err.Error())
		}
		// 记录成功信息
		logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
			" to remote file: "+remoteFilePath+" finish", nil)
		return nil
	})
}

// 创建部署的folder
func createDeployFolder(c framework.Container) (string, error) {
	appService := c.MustMake(contract.AppKey).(contract.App)
	deployFolder := appService.DeployFolder()

	deployVersion := time.Now().Format("20060102150405")
	versionFolder := filepath.Join(deployFolder, deployVersion)
	if !util.Exists(versionFolder) {
		return versionFolder, os.MkdirAll(versionFolder, os.ModePerm)
	}
	return versionFolder, nil
}

func toLinuxPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
