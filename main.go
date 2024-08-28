package main

import (
	"rental-mobil/controllers"
	"rental-mobil/initializers"
	"rental-mobil/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	// Example test router
	r.POST("/posts", controllers.PostsCreate)
	r.POST("/posts/:id", controllers.PostsCreateById)
	r.PUT("/posts/:id", controllers.PostsUpdate)
	r.GET("/posts", controllers.PostsIndex)
	r.GET("/posts/:id", controllers.PostsShow)
	r.DELETE("/posts/:id", controllers.PostsDelete)

	// users
	r.POST("/register", controllers.RegisterUsers)
	r.POST("/login", controllers.LoginUser)
	// admin
	r.POST("/register-admin", controllers.RegisterAdmin)
	r.POST("/login-admin", controllers.LoginAdmin)
	// get role
	r.GET("/protected-route", middleware.AuthRequired(),
		controllers.ProtectedRouteHandler)

	// orders
	r.POST("/create-orders", middleware.AuthRequired(), controllers.PostOrder)
	r.GET("/orders", middleware.AuthRequired(), controllers.GetListOrders)
	r.PUT("/order/:id", middleware.AuthRequired(), middleware.UserOwnsOrderOrAdmin(), controllers.UpdateDataOrder)
	r.DELETE("/order/:id", middleware.AuthRequired(), controllers.DeleteOrder)

	// cars
	r.POST("/create-cars", middleware.AuthRequired(), middleware.AdminRequired(), controllers.PostCar)
	r.GET("/cars", controllers.GetListCar)
	r.GET("/car/:id", controllers.GetListCarByID)
	r.PUT("/car/:id", middleware.AuthRequired(), middleware.AdminRequired(), controllers.UpdateDataCar)
	r.DELETE("/car/:id", middleware.AuthRequired(), middleware.AdminRequired(), controllers.DeleteCar)

	r.Run(":4000")
}
