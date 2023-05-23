package reset

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"minik8s/cmd/edit"
	"minik8s/pkg/client"

	// "fmt"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset initial status of k8s",
	Long:  "this is the main cmd to reset initial status of k8s",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	lines := edit.ResetLog()
	for _, line := range lines {
		var key, tp string
		if _, err := fmt.Sscanf(line, "%s %s", &key, &tp); err != nil {
			log.Fatal(err)
		}
		client.Delete_object(key, tp)
	}
}

func init() {
}

func Reset() *cobra.Command {
	return resetCmd
}
