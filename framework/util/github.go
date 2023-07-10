package util

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/go-github/v39/github"
)

func ResolveRateLimitErrorByAuth(err error) (*github.Client, error) {
	if _, ok := err.(*github.RateLimitError); ok {
		var client *github.Client
		perPage := 10
		opts := &github.ListOptions{Page: 1, PerPage: perPage}

		fmt.Println("使用github.com帐号登录方式来增加调用次数")
		githubUserName := ""
		prompt := &survey.Input{
			Message: "请输入github帐号用户名：",
		}
		if err := survey.AskOne(prompt, &githubUserName); err != nil {
			fmt.Println("任务终止：" + err.Error())
			return nil, err
		}
		githubPassword := ""
		promptPwd := &survey.Password{
			Message: "请输入github帐号密码：",
		}
		if err := survey.AskOne(promptPwd, &githubPassword); err != nil {
			fmt.Println("任务终止：" + err.Error())
			return nil, err
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					req.SetBasicAuth(githubUserName, githubPassword)
					return nil, nil
				},
			},
		}
		client = github.NewClient(httpClient)
		releases, rsp, err := client.Repositories.ListReleases(context.Background(), "y19941115mx", "ygo", opts)
		if err != nil {
			fmt.Println("错误提示：" + err.Error())
			fmt.Println("用户名密码错误，请重新开始")
			return client, nil
		}
		if len(releases) == 0 {
			fmt.Println("用户名密码错误，请重新开始")
			return client, nil
		}
		fmt.Println(rsp.Rate.String())
	}

	fmt.Println("github.com的连接异常：" + err.Error())
	return nil, err
}
