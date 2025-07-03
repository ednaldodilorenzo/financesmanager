package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func LogRequests(c *fiber.Ctx) error {
	err := c.Next()
	log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("status", c.Response().StatusCode()).
		Msg("Handled request")
	return err
}
