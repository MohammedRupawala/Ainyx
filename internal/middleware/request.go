package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := generateRequestID()
		c.Set("X-Request-ID", requestID)
		c.Locals("requestId", requestID)
		return c.Next()
	}
}

func Logger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("request_id", requestID(c)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		}

		if err != nil {
			logger.Error("http request completed with error", append(fields, zap.Error(err))...)
			return err
		}

		logger.Info("http request completed", fields...)
		return nil
	}
}

func requestID(c *fiber.Ctx) string {
	value, _ := c.Locals("requestId").(string)
	return value
}

func generateRequestID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(buf[:])
}
