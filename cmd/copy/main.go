package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cobra"
	"log"
	awsclient "main/client"
	"time"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("")

	rootCmd := &cobra.Command{
		Use:   "aac",
		Short: "Copy AWS AMI",
		Run: func(cmd *cobra.Command, args []string) {
			CopyAMI(cmd.Context())
		},
	}

	cobra.CheckErr(rootCmd.Execute())
}

func CopyAMI(ctx context.Context) {
	actions := githubactions.New()

	log.Println("Starting AMI copy process")

	amiId := actions.GetInput("ami-id")
	ssmParamAMIId := actions.GetInput("ssm-param-ami-id")
	awsRegion := actions.GetInput("aws-region")
	amiDescription := actions.GetInput("description")

	client := awsclient.CreateEC2Client(ctx)

	if ssmParamAMIId != "" && amiId == "" {
		ssmClient := awsclient.CreateSSMClient(ctx)
		getParameterOutput, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
			Name:           aws.String(ssmParamAMIId),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatalf("failed to get parameter: %v", err)
		}

		amiId = *getParameterOutput.Parameter.Value
	}

	if amiId == "" {
		log.Fatalf("ami-id or ssm-param-ami-id is required")
	}

	describeImagesInput := &ec2.DescribeImagesInput{
		ImageIds: []string{amiId},
	}

	describeImagesOutput, err := client.DescribeImages(ctx, describeImagesInput)
	if err != nil {
		log.Fatal(err)
	}

	name := describeImagesOutput.Images[0].Name

	copyImageInput := &ec2.CopyImageInput{
		Name:          name,
		SourceImageId: aws.String(amiId),
		SourceRegion:  aws.String(awsRegion),
		Description:   aws.String(amiDescription),
	}

	copyImageOutput, err := client.CopyImage(ctx, copyImageInput)
	if err != nil {
		log.Fatal(err)
	}

	newImageID := *copyImageOutput.ImageId

	waiter := ec2.NewImageAvailableWaiter(client)
	err = waiter.Wait(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{newImageID},
	}, 10*time.Minute)
	if err != nil {
		log.Fatalf("waiting for image to become available failed: %v", err)
	}

	snapshots, err := getImageSnapshots(ctx, client, newImageID)
	if err != nil {
		log.Fatalf("error getting snapshot info: %v", err)
	}

	actions.SetOutput("copied-ami-id", newImageID)
	actions.SetEnv("COPIED_AMI_ID", newImageID)

	actions.SetOutput("copied-ami-snapshot-id", snapshots[0])
	actions.SetEnv("COPIED_AMI_SNAPSHOT_ID", snapshots[0])

	log.Printf("Copied AMI: %s", newImageID)
	log.Printf("Copied AMI Snapshot: %s", snapshots[0])
}

func getImageSnapshots(ctx context.Context, client *ec2.Client, imageID string) ([]string, error) {
	output, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{imageID},
	})
	if err != nil {
		return nil, err
	}

	if len(output.Images) == 0 {
		return nil, fmt.Errorf("image %s not found", imageID)
	}

	image := output.Images[0]

	var snapshots []string
	for _, mapping := range image.BlockDeviceMappings {
		if mapping.Ebs != nil && mapping.Ebs.SnapshotId != nil {
			snapshots = append(snapshots, *mapping.Ebs.SnapshotId)
		}
	}

	return snapshots, nil
}
