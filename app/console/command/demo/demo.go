package demo

import (
	"log"

	"github.com/y19941115mx/ygo/framework/cobra"
)

// FooCommand 代表Foo命令
var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo的简要说明",
	Long:    "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		log.Println("execute foo command")
		return nil
	},
}
