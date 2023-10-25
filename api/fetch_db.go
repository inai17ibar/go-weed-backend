package api

import (
	"context"
	"go-weed-backend/db"
	"go-weed-backend/internal/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchAndSaveContribution() {
	client := db.GetDB()                             // MongoDB client
	collection := client.Collection("contributions") // 対応するコレクションを取得

	// 最後のコントリビューションの日付を取得
	var lastContribution model.ContributionDayDB
	opts := options.FindOne().SetSort(bson.D{{Key: "date", Value: -1}})
	err := collection.FindOne(context.TODO(), bson.D{}, opts).Decode(&lastContribution)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("error querying last contribution date: %v\n", err)
		return
	}

	respData, err := CallGithubContributionAPI()
	if err != nil {
		log.Printf("error fetching GitHubAPI data: %v\n", err)
		return
	}

	var contributions []model.ContributionDay
	for _, week := range respData.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		contributions = append(contributions, week.ContributionDays...)
	}

	contributionDaysDB := model.ConvertToDBModels(contributions)

	for _, c := range contributionDaysDB {
		filter := bson.M{"date": c.Date}
		update := bson.M{"$set": &c}
		opts := options.Update().SetUpsert(true) // Upsertオプションを使って、存在しない場合は新しく挿入する

		_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			log.Printf("error updating/creating contribution: %v\n", err)
		}
	}
}

func FetchAndSaveCommits() {
	client := db.GetDB()                       // MongoDB client
	collection := client.Collection("commits") // 対応するコレクションを取得

	// 最後のコミットの日付を取得
	var lastCommit model.MyCommit
	opts := options.FindOne().SetSort(bson.D{{Key: "date", Value: -1}})
	err := collection.FindOne(context.TODO(), bson.D{}, opts).Decode(&lastCommit)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("error querying last commit date: %v\n", err)
		return
	}

	commits, err := CallGithubAllCommitAPI()
	if err != nil {
		log.Fatalf("Error fetching commits from GitHub API: %v", err)
	}

	for _, c := range commits {
		if c.Date.After(lastCommit.Date) {
			_, err := collection.InsertOne(context.TODO(), c)
			if err != nil {
				log.Printf("error inserting commit: %v\n", err)
			}
		}
	}
}
