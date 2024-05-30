package handlers

import (
    "github.com/gofiber/fiber/v2"
)

// HealthCheckHandler é o manipulador para verificar o status de saúde do serviço
func HealthCheckHandler(c *fiber.Ctx) error {
    // Retorna um status "OK" como resposta
    return c.JSON(fiber.Map{"status": "OK"})
}
