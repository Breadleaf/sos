package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "sos",
		Short: "Simple Object Storage CLI",
	}

	var endpoint string
	root.PersistentFlags().StringVar(
		&endpoint,
		"endpoint",
		"http://localhost:8080",
		"SOS server endpoint",
	)

	put := &cobra.Command{
		Use:   "put <bucket> <key> <file>",
		Short: "Upload a file",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%#v, %#v, %#v", cmd, args, endpoint)
		},
	}
	root.AddCommand(put)

	cobra.CheckErr(root.Execute())
}
