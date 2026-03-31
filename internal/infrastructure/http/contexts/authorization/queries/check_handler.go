package queries

import (
	authorizationqueries "kali-auth-context/internal/application/authorization/queries"
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/http/shared"

	"github.com/gofiber/fiber/v2"
)

type CheckHandler struct {
	query *authorizationqueries.AuthorizeQuery
}

func NewCheckHandler(query *authorizationqueries.AuthorizeQuery) *CheckHandler {
	return &CheckHandler{query: query}
}

type checkRequest struct {
	TenantId string `json:"tenant_id"`
	UserId   string `json:"user_id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

func (h *CheckHandler) Handle(c *fiber.Ctx) error {
	var req checkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	decision, err := h.query.Handle(c.UserContext(), &authorizationqueries.AuthorizeDto{
		TenantId: identity.TenantId(req.TenantId),
		UserId:   identity.UserId(req.UserId),
		Resource: req.Resource,
		Action:   req.Action,
	})
	if err != nil {
		return shared.WriteError(c, err)
	}
	if !decision.Allowed {
		return c.Status(fiber.StatusForbidden).JSON(decision)
	}

	return c.Status(fiber.StatusOK).JSON(decision)
}
