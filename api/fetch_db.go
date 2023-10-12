package api

import (
	"go-weed-backend/db"
	"go-weed-backend/internal/model"
	"log"

	"github.com/jinzhu/gorm"
)

func FetchAndSaveContribution() {
	dbInstance := db.GetDB() // この変数を以下のコードで利用

	// 最後のコントリビューションの日付を取得
	var lastContribution model.ContributionDayDB
	dbInstance.Order("date desc").First(&lastContribution)

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
		existingData := model.ContributionDayDB{}
		if err := dbInstance.Where("date = ?", c.Date).First(&existingData).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := dbInstance.Create(&c).Error; err != nil {
					log.Printf("error creating contribution: %v\n", err)
				}
			} else {
				log.Printf("error querying contribution: %v\n", err)
			}
		} else {
			if err := dbInstance.Model(&existingData).Updates(&c).Error; err != nil {
				log.Printf("error updating contribution: %v\n", err)
			}
		}
	}
}

func FetchAndSaveCommits() {
	dbInstance := db.GetDB() // この変数を以下のコードで利用

	// 最後のコミットの日付を取得
	var lastCommit model.MyCommit
	dbInstance.Order("date desc").First(&lastCommit)

	var commits []model.MyCommit

	commits, err := CallGithubAllCommitAPI()
	if err != nil {
		log.Fatalf("Error fetching commits from GitHub API: %v", err)
	}

	for _, c := range commits {
		if c.Date.After(lastCommit.Date) {
			dbInstance.Save(&c)
		}
	}
}
