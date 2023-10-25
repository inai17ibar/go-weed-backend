package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var instance *mongo.Client
var dbName = "testDB"

func InitDB(uri string) {

	if instance != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	instance, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	// defer instance.Disconnect(ctx)

	err = instance.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping: %s", err)
	}
	//fmt.Println("Successfully connected to MongoDB!")
}

func GetDB() *mongo.Database {
	return instance.Database(dbName)
}

func CloseDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := instance.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}

func EnsureCollectionExists(client *mongo.Client, dbName, collectionName string) error {
	// コレクションの一覧を取得
	collections, err := client.Database(dbName).ListCollectionNames(context.TODO(), bson.M{"name": collectionName})
	if err != nil {
		return err
	}

	// コレクションが存在しない場合、新しいコレクションを作成
	if len(collections) == 0 {
		err := client.Database(dbName).CreateCollection(context.TODO(), collectionName, options.CreateCollection().SetCapped(false))
		if err != nil {
			return err
		}
	}

	return nil
}
