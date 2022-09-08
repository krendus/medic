package main

import (
	"log"
	"os"
	"time"

	// _ "medic/docs"
	"medic/entity"

	"github.com/gin-contrib/cors"
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

// @host                       127.0.0.1:8000
// @BasePath                   /api/v1
// @schemes                    http
// @query.collection.format    multi
// @securityDefinitions.basic  BasicAuth
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		//  return origin == "https://mseamless.herokuapp.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	api := r.Group("/api/v1")
	{
		api.POST("/user/signup/:role", entity.Signup)
		api.POST("/user/signin", entity.Signin)

		api.POST("/book/appointment", entity.BookAppoitment)
		api.GET("/appointments", entity.GetAppointments)
		api.PUT("/appointment/:id", entity.UpdateAppointment)

	}

	r.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + port)
}
