package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbInstance() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	MongoDb := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully connected to mongodb")
	return client
}

var Client *mongo.Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	databaseName := os.Getenv("DATABASE_NAME")
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}

func GetMongoDoc(colName *mongo.Collection, filter interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data User

	if err := colName.FindOne(ctx, filter).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func GetMongoDocs(colName *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data []bson.M

	filterCusor, err := colName.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	if err := filterCusor.All(ctx, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func CreateMongoDoc(colName *mongo.Collection, data interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertNum, insertErr := colName.InsertOne(ctx, data)
	if insertErr != nil {
		return nil, insertErr
	}

	return insertNum, nil
}

func UpdateMongoDoc(colName *mongo.Collection, filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// _id, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, err
	// }
	// filter := bson.D{{Key: "_id", Value: _id}}

	updateData := bson.D{{Key: "$set", Value: data}}
	res, err := colName.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return nil, err
	}

	return res, nil

}
