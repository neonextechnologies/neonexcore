package user

import (
	"neonexcore/internal/core"

	"github.com/gofiber/fiber/v2"
)

func (m *UserModule) Routes(app *fiber.App, c *core.Container) {
	group := app.Group("/user")

	group.Get("/", func(ctx *fiber.Ctx) error {
		controller := core.Resolve[*UserController](c)
		return controller.GetUsers(ctx)
	})

	group.Get("/search", func(ctx *fiber.Ctx) error {
		controller := core.Resolve[*UserController](c)
		return controller.SearchUsers(ctx)
	})

	group.Get("/:id", func(ctx *fiber.Ctx) error {
		controller := core.Resolve[*UserController](c)
		return controller.GetUserByID(ctx)
	})

	group.Post("/", func(ctx *fiber.Ctx) error {
		controller := core.Resolve[*UserController](c)
		return controller.CreateUser(ctx)
	})

	group.Put("/:id", func(ctx *fiber.Ctx) error {
		controller := core.Resolve[*UserController](c)
		return controller.UpdateUser(ctx)
	})

	group.Delete("/:id", func(ctx *fiber.Ctx) error {
		controller := core.Resolve[*UserController](c)
		return controller.DeleteUser(ctx)
	})
}
