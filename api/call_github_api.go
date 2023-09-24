package api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func CallGithubAPI() {
	apiKey, exists := os.LookupEnv("GITHUB_API_KEY")
	if !exists {
		log.Fatal("Error: GITHUB_API_KEY not set")
	}

	fmt.Println("API Key:", apiKey)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Get all repositories
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			log.Fatalf("Error fetching repositories: %v", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Get commit history for each repository
	for _, repo := range allRepos {
		commits, _, err := client.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, nil)
		if err != nil {
			log.Printf("Error fetching commits for repository %s: %v", *repo.Name, err)
			continue
		}
		fmt.Printf("Commits in %s:\n", *repo.Name)
		for _, commit := range commits {
			fmt.Printf("  %s - %s\n", *commit.SHA, *commit.Commit.Message)
		}
	}
}
