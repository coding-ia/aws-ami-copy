package cmd

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cobra"
	"log"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Delete copied AMI",
	Run: func(cmd *cobra.Command, args []string) {
		Cleanup(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}

func Cleanup(ctx context.Context) {
	actions := githubactions.New()

	imageId := actions.Getenv("COPIED_AMI_ID")
	snapshotId := actions.Getenv("COPIED_AMI_SNAPSHOT_ID")

	client := createEC2Client(ctx)

	if imageId != "" {
		_, err := client.DeregisterImage(ctx, &ec2.DeregisterImageInput{
			ImageId: aws.String(imageId),
		})
		if err != nil {
			log.Fatalf("failed to deregister image: %v", err)
		}
	}

	if snapshotId != "" {
		_, err := client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
			SnapshotId: aws.String(snapshotId),
		})
		if err != nil {
			log.Fatalf("failed to delete snapshot: %v", err)
		}
	}
}
