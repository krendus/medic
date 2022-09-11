package main

import (
	"log"
	"os"

	_ "medic/docs"
	"medic/entity"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           medic API
// @version         1.0
// @description     This is the API serving the medic frontend
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host                       medic0.herokuapp.com
// @BasePath                   /api/v1
// @schemes                    https
// @query.collection.format    multi
// @securityDefinitions.basic  BasicAuth
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}

	r := gin.Default()
	config := CORSMiddleware()
	r.Use(config)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	api := r.Group("/api/v1")
	{
		api.POST("/user/signup/:role", entity.Signup)
		api.POST("/user/signin", entity.Signin)

		api.POST("/book/appointment", entity.BookAppoitment)
		api.GET("/appointments/all", entity.GetAllAppointments)
		api.GET("/appointments/:id", entity.GetUserAppointments)
		api.PUT("/appointment/:id", entity.UpdateAppointment)
		api.GET("/doctors", entity.GetDoctors)
		// .Use(database.Authentication)

	}

	r.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, token, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
