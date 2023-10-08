package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Todo struct {
	gorm.Model
	Title        string `json:"Title"`
	Completed    bool   `json:"Completed"`
	Created_Date string `json:"Created_date"`
}

type MyCommit struct {
	gorm.Model
	SHA       string    `json:"Sha"`
	Message   string    `json:"Message"`
	Date      time.Time `json:"Date"`
	Additions int       `json:"Additions"`
	Deletions int       `json:"Deletions"`
	Total     int       `json:"Total"`
}

type CommitData struct {
	Date      string `json:"Date"`
	Count     int    `json:"Count"`
	Additions int    `json:"Additions"`
	Deletions int    `json:"Deletions"`
	Total     int    `json:"Total"`
}

// GraphQLResponse represents the structure of the GraphQL response
type GraphQLResponse struct {
	Data Data `json:"data" graphql:"data"`
}

// Data represents the data field in the GraphQL response
type Data struct {
	User User `json:"user" graphql:"user"`
}

// User represents the user field in the GraphQL response
type User struct {
	ContributionsCollection ContributionsCollection `json:"contributionsCollection" graphql:"contributionsCollection"`
}

// ContributionsCollection represents the contributionsCollection field in the GraphQL response
type ContributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar" graphql:"contributionCalendar"`
}

// ContributionCalendar represents the contributionCalendar field in the GraphQL response
type ContributionCalendar struct {
	TotalContributions int    `json:"totalContributions" graphql:"totalContributions"`
	Weeks              []Week `json:"weeks" graphql:"weeks"`
}

// Week represents each week's data in the contributionCalendar
type Week struct {
	ContributionDays []ContributionDay `json:"contributionDays" graphql:"contributionDays"`
}

type ContributionDay struct {
	Date              string `json:"date" graphql:"date"`
	ContributionCount int    `json:"contributionCount" graphql:"contributionCount"`
}
