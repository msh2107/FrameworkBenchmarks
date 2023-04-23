package main

import (
	"errors"
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

	r.GET("/users", getUsers)
	r.POST("/user", newUser)
	r.GET("/user/:id", getUser)
	r.PATCH("/user/:id", updateUser)
	r.DELETE("/user/:id", deleteUser)

	r.GET("/sleep", sleep)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func sleep(ctx *gin.Context) {
	time.Sleep(5 * time.Second)
	_, err := ctx.Writer.WriteString("Вы проспали 5 секунд")
	if err != nil {
		return
	}
}

func sendImage(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/image/bug.jpg")
}

func getUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, users)
}

func newUser(ctx *gin.Context) {
	var jsonUser user
	err := ctx.ShouldBindJSON(&jsonUser)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	users = append(users, jsonUser)
	ctx.JSON(http.StatusOK, jsonUser)
}

func getUser(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	for _, user := range users {
		if user.Id == i {
			ctx.JSON(http.StatusOK, user)
			return
		}
	}
	_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no such user"))

}

func updateUser(ctx *gin.Context) {
	var jsonUser user
	err := ctx.ShouldBindJSON(&jsonUser)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	for index, user := range users {
		if user.Id == i {
			users[index] = jsonUser
			ctx.JSON(http.StatusOK, jsonUser)
			return
		}
	}
	_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no such user"))
}

func deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	for index, user := range users {
		if user.Id == i {
			ctx.JSON(http.StatusOK, user)
			users = append(users[:index], users[index+1:]...)
			return
		}

	}
	_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no such user"))
}
