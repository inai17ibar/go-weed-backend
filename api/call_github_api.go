package api

import (
	"context"
	"fmt"
	"go-weed-backend/internal/model"
	"log"
	"os"
	"sort"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func CallGithubAllCommitAPI() ([]model.MyCommit, error) {
	apiKey, exists := os.LookupEnv("API_KEY_GITHUB")
	if !exists {
		log.Fatal("Error: API_KEY_GITHUB not set")
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

	var allCommits []model.MyCommit
	// Get commit history for each repository
	for _, repo := range allRepos {
		if *repo.Fork {
			continue
		}
		commits, _, err := client.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, nil)
		if err != nil {
			log.Printf("Error fetching commits for repository %s: %v", *repo.Name, err)
			continue
		}

		targetUser := "inai17ibar" // Replace with your GitHub username

		fmt.Printf("Commits in %s:\n", *repo.Name)
		for _, commit := range commits {
			authorLogin := commit.Author.GetLogin()
			committerLogin := commit.Committer.GetLogin()

			// Only include commits authored or committed by the target user
			if authorLogin == targetUser || committerLogin == targetUser {
				// コミットの統計情報を取得
				commitDetail, _, err := client.Repositories.GetCommit(ctx, *repo.Owner.Login, *repo.Name, *commit.SHA, nil)
				if err != nil {
					log.Printf("Error fetching commit detail for commit %s in repository %s: %v", *commit.SHA, *repo.Name, err)
					continue
				}
				stats := commitDetail.GetStats()
				//これをしたらおそくなる

				myCommit := model.MyCommit{
					SHA:       *commit.SHA,
					Message:   *commit.Commit.Message,
					Date:      commit.Commit.Author.GetDate(),
					Additions: stats.GetAdditions(), // 追加された行数
					Deletions: stats.GetDeletions(), // 削除された行数
					Total:     stats.GetTotal(),     // 合計変更行数
				}
				allCommits = append(allCommits, myCommit)
				fmt.Printf("  %s - %s\n", *commit.SHA, *commit.Commit.Message)
			}
		}
	}

	// allCommitsを日時でソート
	sort.Slice(allCommits, func(i, j int) bool {
		return allCommits[i].Date.After(allCommits[j].Date)
	})

	return allCommits, nil
}

func CallGithubContributionAPI() (*model.GraphQLResponse, error) {
	apiKey, exists := os.LookupEnv("API_KEY_GITHUB")
	if !exists {
		log.Fatal("Error: API_KEY_GITHUB not set")
	}

	// GitHub APIクライアントの初期化
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var response model.GraphQLResponse

	// 実際のGraphQLクエリ
	var query struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar model.ContributionCalendar
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"` //@arguments(login: { type: \"String!\", value: $login })"
	}

	currentDate := time.Now()                 //.Format("2006-01-02T15:04:05Z")
	startDate := time.Now().AddDate(-1, 0, 0) //.Format("2006-01-02T15:04:05Z")
	fmt.Printf("startDate: %s\n", startDate)
	fmt.Printf("currentDate: %s\n", currentDate)

	// クエリの変数をセット
	variables := map[string]interface{}{
		"login": githubv4.String("inai17ibar"),
		"from":  githubv4.DateTime{Time: startDate},
		"to":    githubv4.DateTime{Time: currentDate},
	}

	// クエリ実行
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Printf("Error fetching GitHub data: %v\n", err)
		return nil, err
	}

	// データの処理
	fmt.Printf("Total Contributions: %d\n", query.User.ContributionsCollection.ContributionCalendar.TotalContributions)

	response.Data.User.ContributionsCollection.ContributionCalendar = query.User.ContributionsCollection.ContributionCalendar

	//fmt.Printf("GraphQL Response: %+v\n", response)
	return &response, nil
}
