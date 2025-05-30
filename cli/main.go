package main

import (
	"fmt"
	"io"
	"os"

	"github.com/breadleaf/sos/pkg/client"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "sos",
		Short: "Simple Object Storage CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(endpoint)
			f, err := os.Open(args[2])
			if err != nil {
				return err
			}
			defer f.Close()
			return c.PutObject(args[0], args[1], f)
		},
	}

	get := &cobra.Command{
		Use:   "get <bucket> <key> <file>",
		Short: "Download an object",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(endpoint)
			reader, err := c.GetObject(args[0], args[1])
			if err != nil {
				return err
			}
			defer reader.Close()
			out, err := os.Create(args[2])
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, reader)
			return err
		},
	}

	listBuckets := &cobra.Command{
		Use:   "lsBuckets",
		Short: "List all buckets",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(endpoint)
			buckets, err := c.ListBuckets()
			if err != nil {
				return err
			}
			for _, b := range buckets {
				fmt.Println(b)
			}
			return nil
		},
	}

	listObjects := &cobra.Command{
		Use:   "lsObjects <bucket>",
		Short: "List all objects",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(endpoint)
			objs, err := c.ListObjects(args[0])
			if err != nil {
				return err
			}
			for _, o := range objs {
				fmt.Println(o)
			}
			return nil
		},
	}

	deleteBucket := &cobra.Command{
		Use:   "rmBucket <bucket>",
		Short: "Delete a bucket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(endpoint)
			return c.DeleteBucket(args[0])
		},
	}

	deleteObject := &cobra.Command{
		Use:   "rmObject <bucket> <key>",
		Short: "Delete an object",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient(endpoint)
			return c.DeleteObject(args[0], args[1])
		},
	}

	root.AddCommand(put, get, listBuckets, listObjects, deleteBucket, deleteObject)

	cobra.CheckErr(root.Execute())
}
