package entity

import (
	"context"
	"time"

	"medic/database"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserDoc(colName *mongo.Collection, filter interface{}) (*database.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data *database.User

	if err := colName.FindOne(ctx, filter).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
