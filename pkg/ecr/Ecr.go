package ecr

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

var ecrUrlRegex = regexp.MustCompile(`^(\d+)\.dkr\.ecr\.(?:[^.]+)\.amazonaws\.com\/(.+)$`)

func IsEcrUrl(url string) bool {
	matches := ecrUrlRegex.FindStringSubmatch(url)
	return len(matches) <= 3
}

func FindAvailableTags(url string, constraint string) ([]string, error) {
	accountId, repoPath, err := parseEcrUrl(url)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := ecr.New(sess)
	input := &ecr.DescribeImagesInput{
		RegistryId:     aws.String(accountId),
		RepositoryName: aws.String(repoPath),
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		return nil, err
	}

	for _, detail := range result.ImageDetails {
		for _, tag := range detail.ImageTags {
			if *tag == constraint {
				return convertPtrSliceToStringSlice(detail.ImageTags), nil
			}
		}
	}

	return nil, fmt.Errorf("tag not found %s:%s", url, constraint)
}

func parseEcrUrl(url string) (accountID, repoPath string, err error) {
	matches := ecrUrlRegex.FindStringSubmatch(url)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("input URL does not match the expected ECR format")
	}

	accountID = matches[1]
	repoPath = matches[2]
	return accountID, repoPath, nil
}

func convertPtrSliceToStringSlice(ptrSlice []*string) []string {
	stringSlice := make([]string, len(ptrSlice))
	for i, ptr := range ptrSlice {
		if ptr != nil {
			stringSlice[i] = *ptr
		}
	}
	return stringSlice
}
