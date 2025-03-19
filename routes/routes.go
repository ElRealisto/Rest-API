package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents) // GET, POST, PUT, PATCH, DELETE
	server.GET("/events/:id", getSingleEvent)

	authenticated := server.Group("/")          // creating a group of protected routes
	authenticated.Use(middlewares.Authenticate) // adding particular middleware for a group of routes
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)

	server.DELETE("/events/:id", middlewares.Authenticate, deleteEvent) //adding an authentication middleware into particular route
	server.POST("/signup", signup)
	server.POST("/login", login)
	server.GET("/users/", getUsers)

}
