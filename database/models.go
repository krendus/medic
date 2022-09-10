package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	FullName   string             `json:"full_name" validate:"required"`
	Username   string             `json:"username" validate:"required"`
	Email      string             `json:"email" validate:"required"`
	Role       string             `json:"role"`
	Password   string             `json:"password" validate:"required"`
	Created_At time.Time          `json:"created_at"`
}

type UserModel struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Appointment struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName   string             `json:"first_name" validate:"required"`
	LastName    string             `json:"last_name" validate:"required"`
	PhoneNumber string             `json:"phone_number" validate:"required"`
	Email       string             `json:"email" validate:"required"`
	Time        string             `json:"time" validate:"required"`
	Date        string             `json:"date" validate:"required"`
	Specialist  string             `jon:"specialist" validate:"required"`
	Message     string             `json:"message" validate:"required"`
	Created_At  time.Time          `json:"created_at"`
}

type AppointmentModel struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Time        string `json:"time"`
	Date        string `json:"date"`
	Specialist  string `jon:"specialist"`
	Message     string `json:"message"`
}

type PubUser struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func PublicUser(user *User) *PubUser {
	return &PubUser{
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
	}
}
