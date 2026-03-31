package commands

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type UpdateRoleDto struct {
	Id          identity.RoleId
	TenantId    identity.TenantId
	Name        string
	Description string
}

type UpdateRoleCommand struct {
	repo ports.IUpdateRoleCommandRepository
}

func NewUpdateRoleCommand(repo ports.IUpdateRoleCommandRepository) *UpdateRoleCommand {
	return &UpdateRoleCommand{repo: repo}
}

func (c *UpdateRoleCommand) Execute(ctx context.Context, dto *UpdateRoleDto) error {
	role, err := identity.NewRole(dto.Id, dto.TenantId, dto.Name, dto.Description)
	if err != nil {
		return err
	}

	return c.repo.Update(ctx, role)
}
