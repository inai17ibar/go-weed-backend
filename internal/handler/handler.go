package handler

import (
	"log"
	"net/http"
	"sort"

	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"go-weed-backend/internal/model"
	"go-weed-backend/internal/util"
)

// dbはデータベースへの参照を保持します。
// 実際のプロジェクトでは、データベースへのアクセス方法をより適切に構造化することが重要です。
var db *gorm.DB

// var svc *s3.S3
// var bucketName string
// var fileKey string

// Initはhandlerパッケージを初期化します。
func Init(database *gorm.DB) {
	db = database
	// svc = s3Service
	// bucketName = bName
	// fileKey = fKey
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []model.Todo
	db.Find(&todos)

	// JSONデータとしてクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := util.ConvertToJapanTime(time.Now()).Truncate(24 * time.Hour)
	todo.Created_Date = currentDate.Format("2006-01-02")

	db.Create(&todo)

	// Debug: Check the content of todo after insertion
	//fmt.Printf("After insertion: %+v\n", todo)

	// 新しいTodoの情報をJSONとしてレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// func uploadTodoToS3(todo model.Todo) error {
// 	// データベースファイルを読み込み
// 	databaseBytes, err := ioutil.ReadFile("local-database.db")
// 	if err != nil {
// 		return err
// 	}

// 	// データベースファイルをS3にアップロード
// 	_, err = svc.PutObject(&s3.PutObjectInput{ //空になる可能性がある。エラーハンドリングできていないかも
// 		Bucket:        aws.String(bucketName),
// 		Key:           aws.String(fileKey),
// 		Body:          bytes.NewReader(databaseBytes),       // バイナリデータを指定
// 		ContentLength: aws.Int64(int64(len(databaseBytes))), // データの長さを指定
// 	})

// 	return err
// }

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("ID")
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	var todo model.Todo
	db.First(&todo, todoID)
	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "TODO item not found")
		return
	}

	db.Delete(&todo)
	fmt.Fprintln(w, "TODO item deleted successfully")
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("ID")
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	var todo model.Todo
	db.First(&todo, todoID)
	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "TODO item not found")
		return
	}

	var updateData struct {
		Title     string `json:"Title"`
		Completed bool   `json:"Completed"`
		Favorite  bool   `json:"Favorite"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}

	todo.Title = updateData.Title
	todo.Completed = updateData.Completed
	todo.Favorite = updateData.Favorite

	db.Save(&todo)
	fmt.Fprintln(w, "TODO item updated successfully")

	// Update後のTodoの情報をJSONとしてレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func GetTodosByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("Created_date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing 'Created_date' parameter in the request")
		return
	}

	var todos []model.Todo
	db.Where("created_date = ?", date).Find(&todos)

	// TODOリストをJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func AggregateCommitDataByDate(w http.ResponseWriter, r *http.Request) {
	var commits []model.MyCommit
	if err := db.Find(&commits).Error; err != nil {
		log.Fatalf("Failed to get commits from database: %v", err)
	}

	commitDataByDate := make(map[string]*model.CommitData)
	for _, commit := range commits {
		date := commit.Date.Format("2006-01-02")

		// 日付ごとのデータがなければ新しく作成
		if commitDataByDate[date] == nil {
			commitDataByDate[date] = &model.CommitData{}
		}

		// 日付ごとのコミット数とコード変更量を累計
		commitDataByDate[date].Count++
		commitDataByDate[date].Additions += commit.Additions
		commitDataByDate[date].Deletions += commit.Deletions
		commitDataByDate[date].Total += commit.Total
	}

	var commitDataList []model.CommitData
	for date, data := range commitDataByDate {
		commitDataList = append(commitDataList, model.CommitData{
			Date:      date,
			Count:     data.Count,
			Additions: data.Additions,
			Deletions: data.Deletions,
			Total:     data.Total,
		})
	}

	//日付で昇順ソートする
	sort.Slice(commitDataList, func(i, j int) bool {
		iDate, err := time.Parse("2006-01-02", commitDataList[i].Date)
		if err != nil {
			// エラーハンドリング
			return false
		}

		jDate, err := time.Parse("2006-01-02", commitDataList[j].Date)
		if err != nil {
			// エラーハンドリング
			return false
		}

		// iDateがjDateより前（昇順）であればtrueを返す
		return iDate.Before(jDate)
	})

	jsonData, err := json.Marshal(commitDataList)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
