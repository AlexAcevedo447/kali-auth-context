package middleware

import (
	"crypto/sha256"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"kali-auth-context/internal/ports"
)

// NewIdempotencyMiddleware returns a Fiber handler that enforces idempotency for
// mutating POST endpoints. Clients must send a unique Idempotency-Key header (UUID).
//
// Behaviour:
//   - Missing key        → 400 Bad Request
//   - Key seen before, same fingerprint → replay cached status + body, header Idempotent-Replayed: true
//   - Key seen before, different fingerprint → 422 (key reuse with different request)
//   - First time         → process normally and persist result for 24 h
func NewIdempotencyMiddleware(repo ports.IIdempotencyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("Idempotency-Key")
		if key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Idempotency-Key header is required",
			})
		}

		fingerprint := computeFingerprint(c.Method(), c.Path(), c.Body())

		record, err := repo.Find(c.UserContext(), key)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		if record != nil {
			if record.Fingerprint != fingerprint {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"error": "idempotency key was already used for a different request",
				})
			}
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			c.Set("Idempotent-Replayed", "true")
			return c.Status(record.StatusCode).Send(record.Body)
		}

		if err := c.Next(); err != nil {
			return err
		}

		statusCode := c.Response().StatusCode()
		body := make([]byte, len(c.Response().Body()))
		copy(body, c.Response().Body())

		_ = repo.Save(c.UserContext(), &ports.IdempotencyRecord{
			Key:         key,
			Fingerprint: fingerprint,
			StatusCode:  statusCode,
			Body:        body,
		})

		return nil
	}
}

func computeFingerprint(method, path string, body []byte) string {
	h := sha256.New()
	h.Write([]byte(method))
	h.Write([]byte("|"))
	h.Write([]byte(path))
	h.Write([]byte("|"))
	h.Write(body)
	return fmt.Sprintf("%x", h.Sum(nil))
}
