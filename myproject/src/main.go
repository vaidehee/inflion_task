package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	initDB()
	router := gin.Default()
	router.GET("/person/:person_id/info", getPersonInfo)
	router.POST("/person/create", createPerson)
	router.Run(":8080")
}
