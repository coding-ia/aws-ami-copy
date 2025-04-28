package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cobra"
	"log"
	awsclient "main/client"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("")

	rootCmd := &cobra.Command{
		Use:   "aac_cleanup",
		Short: "Removes AMI and snapshot for copied AMI",
		Run: func(cmd *cobra.Command, args []string) {
			Cleanup(cmd.Context())
		},
	}

	cobra.CheckErr(rootCmd.Execute())
}

func Cleanup(ctx context.Context) {
	actions := githubactions.New()

	imageId := actions.Getenv("COPIED_AMI_ID")
	snapshotId := actions.Getenv("COPIED_AMI_SNAPSHOT_ID")

	client := awsclient.CreateEC2Client(ctx)

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
