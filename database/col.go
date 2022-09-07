package database

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	UserCollection *mongo.Collection = OpenCollection(Client, "users")
	AppCollection  *mongo.Collection = OpenCollection(Client, "app")
	Validate                         = validator.New()
)
