package helper

import "github.com/gofiber/fiber/v2"

func Err(ctx *fiber.Ctx, statusCode int, message string, err error) error {
	return ctx.Status(statusCode).JSON(fiber.Map{"message": message, "error": err.Error()})
}

func Success(ctx *fiber.Ctx, statusCode int, data any) error {
	return ctx.Status(statusCode).JSON(fiber.Map{"data": data})
}
