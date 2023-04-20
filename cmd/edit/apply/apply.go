package apply

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/client"
	"minik8s/pkg/exeFile"
)

var file string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "add or change the items of k8s",
	Long:  "this is the main cmd to add or change the items int k8s",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	value, key, tp := exeFile.ReadYaml(file)
	client.Put_object(key, value, tp)
	//fmt.Println(value)
	//fmt.Println(key)
	//fmt.Println(tp)
}

func init() {
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "")
}

func Apply() *cobra.Command {
	return applyCmd
}
