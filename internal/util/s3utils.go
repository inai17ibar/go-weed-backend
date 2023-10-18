package util

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3バケットとファイルの情報
const bucketName = "todo-weed"
const fileKey = "database/todos.db"

func ConnectS3AWS() (string, *s3.S3, string, string) {
	// AWSセッションの設定
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"), // リージョンを指定
		},
		Profile: "myprofile-inai17ibar", // ここに使用したいプロファイル名を指定
	}))
	// S3クライアントの作成
	s3Client := s3.New(sess)

	// ファイルの存在を確認
	localDBPath := "local-database.db"
	if _, err := os.Stat(localDBPath); err == nil {
		// ファイルが存在する場合は削除
		if err := os.Remove(localDBPath); err != nil {
			log.Fatal(err)
		}
	}

	// S3からデータベースファイルをダウンロード
	downloadedFile, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer downloadedFile.Body.Close()

	// ダウンロードしたデータベースファイルをローカルに保存
	fileData, err := ioutil.ReadAll(downloadedFile.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(localDBPath, fileData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return localDBPath, s3Client, bucketName, fileKey
}
