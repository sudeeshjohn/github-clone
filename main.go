package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"net/url"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Failed to start the Bot:%s", err)
	}
}

func run() error {
	authToken := os.Getenv("GITHUB_OAUTH_TOKEN")
	if len(authToken) == 0 {
		return fmt.Errorf("the environment GITHUB_OAUTH_TOKEN must be set")
	}
	sourceOrg := os.Getenv("GITHUB_SOURCE_ORG")
	if len(sourceOrg) == 0 {
		return fmt.Errorf("the environment variable GITHUB_SOURCE_ORG must be set")
	}
	sourceRepoSit := os.Getenv("GITHUB_SOURCE_REPO")
	if len(sourceRepoSit) == 0 {
		return fmt.Errorf("the environment variable GITHUB_SOURCE_REPO must be set")
	}

	targetOrg := os.Getenv("GITHUB_TARGET_ORG")
	if len(targetOrg) == 0 {
		return fmt.Errorf("the environment variable GITHUB_SOURCE_ORG must be set")
	}
	targetRepoSit := os.Getenv("GITHUB_TARGET_REPO")
	if len(targetRepoSit) == 0 {
		return fmt.Errorf("the environment variable GITHUB_TARGET_REPO must be set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_OAUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	entURL := os.Getenv("GITHUB_ENTERPRISE_URL")
	if len(entURL) == 0 {
		fmt.Printf("the environment variable GITHUB_ENTERPRISE_URL not set, assuming that the repos are in github.com")
	} else {
		enterpriseURL, err := url.Parse(os.Getenv("GITHUB_ENTERPRISE_URL"))
		if err != nil {
			fmt.Println(err)
		}
		client.BaseURL = enterpriseURL
	}
	/*sourceOrg := "isdls"
	sourceRepoSit := "sang"
	targetOrg := "isdls"
	targetRepoSit := "sing"*/

	_, _, _ = createMilestones(ctx, client, sourceOrg, sourceRepoSit, targetOrg, targetRepoSit)
	_, _, _ = createLabels(ctx, client, sourceOrg, sourceRepoSit, targetOrg, targetRepoSit)
	_, _, _ = createIssues(ctx, client, sourceOrg, sourceRepoSit, targetOrg, targetRepoSit)
	return nil

}
