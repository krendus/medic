package entity

import (
	"errors"
	"fmt"
	"medic/database"
	"medic/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func CreateDoctor(c *gin.Context) {
// 	var d database.Doctor

// 	if err := c.BindJSON(&d); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	if err := database.Validate.Struct(&d); err != nil {
// 		c.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	filter := bson.D{{Key: "email", Value: d.DoctorName}}
// 	_, emailErr := database.GetMongoDoc(database.UserCollection, filter)
// 	if emailErr != nil {
// 		d.Created_At = time.Now()
// 		d.ID = primitive.NewObjectID()

// 		d.Role = "doctor"

// 		_, insertErr := database.CreateMongoDoc(database.UserCollection, &d)
// 		if insertErr != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": insertErr.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"user": userInfo,
// 		})
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": errors.New("doctor name taken").Error(),
// 		})
// 		return
// 	}
// }

// create new user 	godoc
// @Summary      create new user
// @Description  this endpoint is used create a user with role as either patient or doctor by passing the role of the user to the URL
// @Tags         user
// @Accept       json
// @Produce      json
// @param        user  body  database.UserModel  true  "user"
// @Success      200
// @Router       /api/v1/user/signup/:role [post]
func Signup(c *gin.Context) {
	var user database.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := database.Validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	param := c.Param("role")
	if param == "" {
		c.JSON(http.StatusBadRequest, errors.New("a user must be passed with role as either doctor or patient").Error())
	}

	emailFilter := bson.D{{Key: "email", Value: user.Email}}
	_, emailErr := database.GetMongoDoc(database.UserCollection, emailFilter)
	if emailErr != nil {
		user.Created_At = time.Now()
		user.ID = primitive.NewObjectID()

		user.Role = param
		hashedPass, _ := helper.Hash(user.Password)
		user.Password = string(hashedPass)

		_, insertErr := database.CreateMongoDoc(database.UserCollection, &user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": insertErr.Error(),
			})
			return
		}

		userInfo := database.PublicUser(&user)

		c.JSON(http.StatusOK, gin.H{
			"user": userInfo,
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("email taken").Error(),
		})
		return
	}
}

// signin user 	godoc
// @Summary      signin user
// @Description  this endpoint is used signin a user
// @Tags         user
// @Accept       json
// @param        appointment  body  database.Appointment  true  "appointment"
// @Produce      json
// @param        user  body  database.SignIn  true  "user"
// @Success      200
// @Router       /api/v1/book/appointment [post]
func Signin(c *gin.Context) {
	var user *database.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := database.Validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	usernameilter := bson.D{{Key: "username", Value: user.Username}}
	foundUser, usernameErr := GetUserDoc(database.UserCollection, usernameilter)
	if usernameErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("username not found").Error(),
		})
		return
	} else {
		if err := helper.VerifyPassword(foundUser.Password, user.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors.New("incorrect password").Error(),
			})
			return
		}
	}

	userInfo := database.PublicUser(foundUser)

	//generate a token for the user on signup
	token, refreshToken, _ := database.GenerateAllTokens(user.ID.Hex())
	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          userInfo,
	})
}

// book appointment 	godoc
// @Summary      book appointment
// @Description  this endpoint is used to create new appointment with a doctor
// @Tags         appointment
// @Accept       json
// @param        appointment  body  database.AppointmentModel  true  "appointment"
// @Produce      json
// @Success      200
// @Router       /api/v1/book/appointment [post]
func BookAppoitment(c *gin.Context) {
	var app *database.Appointment
	if err := c.BindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := database.Validate.Struct(&app); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	app.ID = primitive.NewObjectID()
	app.Created_At = time.Now()

	insertId, insertErr := database.CreateMongoDoc(database.AppCollection, app)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": insertErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"insertId": insertId,
		"message":  "appointment booked",
	})
}

// get all appointments 	godoc
// @Summary      get all appointments
// @Description  this endpoint is used to get all the appointments
// @Tags         appointment
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/v1/appointments [get]
func GetAppointments(c *gin.Context) {
	filter := bson.M{}
	res, err := database.GetMongoDocs(database.AppCollection, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.New("no appointments at the moment").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// get all doctors 	godoc
// @Summary      get all doctors
// @Description  this endpoint is used to get all the doctors
// @Tags         appointment
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/v1/doctors [get]
func GetDoctors(c *gin.Context) {
	filter := bson.M{"role": "doctor"}
	res, err := database.GetMongoDocs(database.UserCollection, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.New("no doctors at the moment").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// update appointments 	godoc
// @Summary      update appointments
// @Description  this endpoint is used to update an appointments
// @Tags         appointment
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/v1/appointment/:id [put]
func UpdateAppointment(c *gin.Context) {
	var app *database.Appointment
	if err := c.BindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := database.Validate.Struct(&app); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.New("id parameter is empty").Error())
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
	}
	filter := bson.D{{Key: "_id", Value: _id}}
	data := bson.D{
		{Key: "firstname", Value: app.FirstName},
		{Key: "lastname", Value: app.LastName},
		{Key: "email", Value: app.Email},
		{Key: "phonenumber", Value: app.PhoneNumber},
		{Key: "time", Value: app.Time},
		{Key: "date", Value: app.Date},
		{Key: "specialist", Value: app.Specialist},
		{Key: "message", Value: app.Message},
	}
	updateRes, updateErr := database.UpdateMongoDoc(database.AppCollection, filter, data)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"update_res": updateRes,
		"message":    "info updated",
	})
}
