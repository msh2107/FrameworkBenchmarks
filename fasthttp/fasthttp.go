package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type user struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Hobby     string `json:"hobby"`
}

var users []user

func main() {
	r := fasthttprouter.New()

	r.GET("/image", sendImage)
	r.ServeFiles("/images/*filepath", "images")

	r.GET("/users", getUsers)
	r.POST("/user", newUser)
	r.GET("/user/:id", getUser)
	r.PATCH("/user/:id", updateUser)
	r.DELETE("/user/:id", deleteUser)

	r.GET("/sleep", sleep)

	log.Fatal(fasthttp.ListenAndServe(":8000", r.Handler))
}

func sleep(ctx *fasthttp.RequestCtx) {
	time.Sleep(5 * time.Second)
	_, err := ctx.WriteString("Вы проспали 5 секунд")
	if err != nil {
		return
	}
}

func sendImage(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/images/bug.jpg", fasthttp.StatusMovedPermanently)
}

func getUsers(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	err := json.NewEncoder(ctx).Encode(users)
	if err != nil {
		return
	}
}

func newUser(ctx *fasthttp.RequestCtx) {
	var jsonUser user
	err := json.Unmarshal(ctx.Request.Body(), &jsonUser)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	users = append(users, jsonUser)
	ctx.Response.Header.Set("Content-Type", "application/json")
	err = json.NewEncoder(ctx).Encode(jsonUser)
	if err != nil {
		return
	}
}

func getUser(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.Error("Bad id", fasthttp.StatusBadRequest)
		return
	}
	for _, user := range users {
		if user.Id == id {
			ctx.Response.Header.Set("Content-Type", "application/json")
			err := json.NewEncoder(ctx).Encode(user)
			if err != nil {
				return
			}
			return
		}
	}
	ctx.Error("No such user", fasthttp.StatusBadRequest)
}

func updateUser(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.Error("Bad id", fasthttp.StatusBadRequest)
		return
	}
	var jsonUser user
	err = json.Unmarshal(ctx.Request.Body(), &jsonUser)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	for index, user := range users {
		if user.Id == id {
			users[index] = jsonUser
			ctx.Response.Header.Set("Content-Type", "application/json")
			err := json.NewEncoder(ctx).Encode(jsonUser)
			if err != nil {
				return
			}
			return
		}
	}
	ctx.Error("No such user", fasthttp.StatusBadRequest)
}

func deleteUser(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.Error("Bad id", fasthttp.StatusBadRequest)
		return
	}
	for index, user := range users {
		if user.Id == id {
			ctx.Response.Header.Set("Content-Type", "application/json")
			err := json.NewEncoder(ctx).Encode(user)
			if err != nil {
				return
			}
			users = append(users[:index], users[index+1:]...)
			return
		}
	}
	ctx.Error("No such user", fasthttp.StatusBadRequest)
}
