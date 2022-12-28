package doc_test

import (
	"bytes"
	"fmt"

	"gitee.com/y19941115mx/ygo/framework/cobra"
	"gitee.com/y19941115mx/ygo/framework/cobra/doc"
)

func ExampleGenManTree() {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "my test program",
	}
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	cobra.CheckErr(doc.GenManTree(cmd, header, "/tmp"))
}

func ExampleGenMan() {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "my test program",
	}
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	out := new(bytes.Buffer)
	cobra.CheckErr(doc.GenMan(cmd, header, out))
	fmt.Print(out.String())
}
