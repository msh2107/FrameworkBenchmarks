package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type user struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Hobby     string `json:"hobby"`
}

var users []user

func main() {
	r := gin.Default()

	r.Static("/image", "./images")

	r.GET("/image", sendImage)

	r.POST("/user", newUser)
	r.GET("/user/:id", getUser)

	r.GET("/sleep", sleep)
	r.DELETE("/user/:id", deleteUser)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func sleep(ctx *gin.Context) {
	time.Sleep(5 * time.Second)
	ctx.Writer.WriteString("Вы проспали 5 секунд")
}

func sendImage(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/image/bug.jpg")
}

func newUser(ctx *gin.Context) {
	var jsonUser user
	err := ctx.ShouldBindJSON(&jsonUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad json",
		})
		return
	}
	users = append(users, jsonUser)
}

func getUser(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad id",
		})
		return
	}
	for _, user := range users {
		if user.Id == i {
			ctx.JSON(http.StatusOK, user)
			return
		}

	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "no such user",
	})

}

func deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad id",
		})
		return
	}
	for index, user := range users {
		if user.Id == i {
			ctx.JSON(http.StatusOK, user)
			users = append(users[:index], users[index+1:]...)

			return
		}

	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "no such user",
	})
}
