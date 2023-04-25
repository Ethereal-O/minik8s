package printTable

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func Print(columns ...string) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, '\t', 0)
	formatStr := ""
	for _, column := range columns {
		formatStr += column + "\t"
	}
	_, err := fmt.Fprintln(w, formatStr)
	if err != nil {
		return
	}
	err = w.Flush()
	if err != nil {
		return
	}
}
