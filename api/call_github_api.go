package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func CallGithubAllCommitAPI() {
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

func CallGithubCommitAPI(commitCount int) ([]*github.RepositoryCommit, error) {
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

	// ユーザーのすべてのリポジトリを取得
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("Error fetching repositories: %v", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// 各リポジトリからコミットを取得し、一つのスライスに結合
	var allCommits []*github.RepositoryCommit
	for _, repo := range allRepos {
		commits, resp, err := client.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, &github.CommitsListOptions{
			ListOptions: github.ListOptions{PerPage: commitCount},
		})
		if err != nil {
			if resp != nil && resp.StatusCode == 409 {
				//log.Printf("Repository %s is empty. Skipping...\n", *repo.Name)
				continue
			}
			return nil, fmt.Errorf("Error fetching commits for repository %s: %v", *repo.Name, err)
		}
		allCommits = append(allCommits, commits...)
	}

	// コミットを日時でソート
	sort.Slice(allCommits, func(i, j int) bool {
		return allCommits[i].GetCommit().GetAuthor().GetDate().After(allCommits[j].GetCommit().GetAuthor().GetDate())
	})

	// 最新のN件のコミットを返す
	if len(allCommits) > commitCount {
		return allCommits[:commitCount], nil
	}
	return allCommits, nil
}
