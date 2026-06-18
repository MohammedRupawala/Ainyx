package routes

import (
	"github.com/Ainyx-backend/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, userHandler *handler.UserHandler) {
	users := app.Group("/users")
	users.Post("", userHandler.Create)
	users.Get("", userHandler.List)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
}
