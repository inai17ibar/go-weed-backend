package handler

import (
	"log"
	"net/http"
	"sort"

	"encoding/json"
	"fmt"
	"time"

	"context"

	"go-weed-backend/db"
	"go-weed-backend/internal/model"
	"go-weed-backend/internal/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	collection := db.GetDB().Collection("todos")

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	var todos []model.Todo
	if err = cur.All(context.TODO(), &todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	decoder := json.NewDecoder(r.Body)

	//fmt.Println("AddTodo: Decoding request body")
	if err := decoder.Decode(&todo); err != nil {
		fmt.Println("Error decoding todo:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println("AddTodo: Successfully decoded request body")

	currentDate := util.ConvertToJapanTime(time.Now()).Truncate(24 * time.Hour)
	todo.Created_Date = currentDate.Format("2006-01-02")

	//fmt.Println("AddTodo: Connecting to DB")
	collection := db.GetDB().Collection("todos")

	//fmt.Println("AddTodo: Inserting data into DB")
	_, err := collection.InsertOne(context.TODO(), todo) //add todo to DB
	if err != nil {
		fmt.Println("Error inserting todo:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 新しいTodoの情報をJSONとしてレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("ID")
	fmt.Printf("Received ID: %s\n", todoIDStr)
	if todoIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	// 文字列からObjectIDへの変換
	todoID, err := primitive.ObjectIDFromHex(todoIDStr)
	if err != nil {
		fmt.Printf("Error in ObjectIDFromHex: %v\n", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	//fmt.Printf("Received ID in DeleteTodo: %s", todoIDStr)
	//fmt.Printf("Converted ObjectID: %v", todoID)

	collection := db.GetDB().Collection("todos")
	filter := bson.M{"_id": todoID}
	_, err = collection.DeleteOne(context.TODO(), filter) //delete todo from DB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "TODO item deleted successfully")
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("ID")
	if todoIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing ID parameter")
		return
	}

	// 文字列からObjectIDへの変換
	todoID, err := primitive.ObjectIDFromHex(todoIDStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Title     string `json:"Title"`
		Completed bool   `json:"Completed"`
		Favorite  bool   `json:"Favorite"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request body")
		return
	}

	collection := db.GetDB().Collection("todos")
	filter := bson.M{"_id": todoID}
	update := bson.M{
		"$set": bson.M{
			"Title":     updateData.Title,
			"Completed": updateData.Completed,
			"Favorite":  updateData.Favorite,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "TODO item updated successfully")
	// Update後のTodoの情報をJSONとしてレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateData)
}

func GetTodosByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("Created_Date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing 'Created_Date' parameter in the request")
		return
	}

	collection := db.GetDB().Collection("todos")

	filter := bson.M{"Created_Date": date}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error finding TODOs by date: %v", err) // Add detailed logging
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Cannot find TODOs by date")
		return
	}
	defer cur.Close(context.TODO()) // Move this right after checking the error

	todos := make([]model.Todo, 0)
	for cur.Next(context.TODO()) {
		var todo model.Todo
		if err := cur.Decode(&todo); err != nil {
			log.Printf("Error decoding todo from database: %v", err) // Add detailed logging
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error decoding todo from database")
			return
		}
		todos = append(todos, todo)
	}

	if err := cur.Err(); err != nil {
		log.Printf("Cursor iteration error: %v", err) // Add detailed logging
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Cursor iteration error")
		return
	}

	// If no todos found with the given date
	if len(todos) == 0 {
		log.Println("No data found with the given filter")
	}

	// Send TODO list in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func AggregateCommitDataByDate(w http.ResponseWriter, r *http.Request) {
	collection := db.GetDB().Collection("commits")
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching commits from database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	//MongoDBのGoドライバのcollection.Find()メソッドは、
	//マッチするドキュメントのカーソルを返します。このカーソル(cur変数)を使用して、結果セットをイテレートし、取得したドキュメントをGoの構造体にデコードすることができます。
	var commits []model.MyCommit
	for cur.Next(context.TODO()) {
		var commit model.MyCommit
		err := cur.Decode(&commit)
		if err != nil {
			log.Println("Error decoding commit:", err)
			continue
		}
		commits = append(commits, commit)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, "Cursor error: "+err.Error(), http.StatusInternalServerError)
		return
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
