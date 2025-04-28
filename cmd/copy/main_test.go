package main

import (
	"os"
	"testing"
)

func TestCopyAMI(t *testing.T) {
	_ = os.Setenv("INPUT_SSM-PARAM-AMI-ID", "/aws/service/eks/optimized-ami/1.32/amazon-linux-2023/x86_64/standard/recommended/image_id")
	_ = os.Setenv("INPUT_AWS-REGION", "us-east-2")

	tempEnv, err := createTempGitHubEnviornment(".env")
	if err != nil {
		t.Fail()
	}
	tempOutput, err := createTempGitHubEnviornment(".output")
	if err != nil {
		t.Fail()
	}
	_ = os.Setenv("GITHUB_ENV", tempEnv)
	_ = os.Setenv("GITHUB_OUTPUT", tempOutput)

	CopyAMI(t.Context())
}

func createTempGitHubEnviornment(fileName string) (string, error) {
	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", err
	}
	_ = tempFile.Close()
	return tempFile.Name(), nil
}
