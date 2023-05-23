package apply

import (
	"github.com/spf13/cobra"
	"minik8s/cmd/edit"
	"minik8s/pkg/client"
	"minik8s/pkg/exeFile"
	// "fmt"
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
	edit.ApplyLog(key, tp)
}

func init() {
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "Path to yaml file")
	applyCmd.MarkFlagRequired("file")
}

func Apply() *cobra.Command {
	return applyCmd
}
