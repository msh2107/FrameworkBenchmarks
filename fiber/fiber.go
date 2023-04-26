package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
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
	app := fiber.New()

	app.Get("/image", sendImage)
	app.Static("/images", "./images")

	app.Get("/users", getUsers)
	app.Post("/user", newUser)
	app.Get("/user/:id", getUser)
	app.Put("/user/:id", updateUser)
	app.Delete("/user/:id", deleteUser)

	app.Get("/sleep", sleep)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}

func sleep(c *fiber.Ctx) error {
	time.Sleep(5 * time.Second)
	return c.SendString("Вы проспали 5 секунд")
}

func sendImage(c *fiber.Ctx) error {
	return c.Redirect("/images/bug.jpg", fiber.StatusMovedPermanently)
}

func getUsers(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.JSON(users)
}

func newUser(c *fiber.Ctx) error {
	var jsonUser user
	err := c.BodyParser(&jsonUser)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	users = append(users, jsonUser)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.JSON(jsonUser)
}

func getUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	for _, user := range users {
		if user.Id == id {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.JSON(user)
		}
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func updateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var jsonUser user
	err = c.BodyParser(&jsonUser)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	for index, user := range users {
		if user.Id == id {
			users[index] = jsonUser
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.JSON(jsonUser)
		}
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func deleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	for index, user := range users {
		if user.Id == id {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			_ = c.JSON(user)
			users = append(users[:index], users[index+1:]...)
			return nil
		}
	}
	return c.SendStatus(fiber.StatusBadRequest)
}
