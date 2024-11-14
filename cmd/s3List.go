package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func init() {
	s3listCommand.Flags().Int32P("limit", "l", 50, "Maximum number of buckets to fetch")
	s3listCommand.Flags().StringP("region", "r", "", "Maximum number of buckets to fetch")
}

var s3listCommand = &cobra.Command{
	Use: "list",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		limit, err := cmd.Flags().GetInt32("limit");
		checkErr(err)

		var region string;
		regionFlag, err := cmd.Flags().GetString("region");
		checkErr(err)
		if regionFlag != "" {
			region = regionFlag
		}

		// Load config
		cfg, err :=  config.LoadDefaultConfig(context.TODO())
		checkErr(err)

		client := s3.NewFromConfig(cfg)

		params := &s3.ListBucketsInput{
			MaxBuckets: &limit,
			BucketRegion: &region,
		}
		output, err := client.ListBuckets(context.TODO(), params)
		checkErr(err)

		w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
		fields := []string{"#", "CreationDate", "Name"}
		heading := strings.Join(fields, "\t")
		fmt.Fprintln(w, heading)

		for index, object := range output.Buckets {
			fields := []string{
				fmt.Sprintf("%d", index + 1),
				object.CreationDate.Format("2006-01-02"),
				*object.Name,
			}
			line := strings.Join(fields, "\t")
			fmt.Fprintln(w, line)
		}
		err = w.Flush()
		checkErr(err)
	},
}
