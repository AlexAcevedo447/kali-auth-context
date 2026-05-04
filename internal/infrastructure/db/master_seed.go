package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/domain/policies"
	"kali-auth-context/internal/infrastructure/config"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type MasterSeeder struct {
	db     *bun.DB
	cfg    *config.Config
	hasher ports.IPasswordHasher
}

func NewMasterSeeder(db *bun.DB, cfg *config.Config, hasher ports.IPasswordHasher) *MasterSeeder {
	return &MasterSeeder{db: db, cfg: cfg, hasher: hasher}
}

func (s *MasterSeeder) Seed(ctx context.Context) error {
	if !s.cfg.SeedMasterEnabled {
		return nil
	}

	password, err := s.cfg.ResolveSeedMasterPassword()
	if err != nil {
		return err
	}

	if err := policies.ValidatePasswordStrength(password); err != nil {
		return fmt.Errorf("invalid SEED_MASTER_PASSWORD: %w", err)
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := s.ensureTenant(ctx, tx); err != nil {
			return err
		}

		masterUserID, err := s.ensureMasterUser(ctx, tx, password)
		if err != nil {
			return err
		}

		adminRoleID, err := s.ensureAdminRole(ctx, tx)
		if err != nil {
			return err
		}

		globalPermissionID, err := s.ensureGlobalAdminPermission(ctx, tx)
		if err != nil {
			return err
		}

		if err := s.ensureRolePermission(ctx, tx, adminRoleID, globalPermissionID); err != nil {
			return err
		}

		if err := s.ensureUserRole(ctx, tx, masterUserID, adminRoleID); err != nil {
			return err
		}

		return nil
	})
}

func (s *MasterSeeder) ensureTenant(ctx context.Context, tx bun.Tx) error {
	var tenant models.TenantModel
	err := tx.NewSelect().
		Model(&tenant).
		Where("id = ?", s.cfg.SeedMasterTenantID).
		Limit(1).
		Scan(ctx)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = tx.NewInsert().
		Model(&models.TenantModel{
			Id:     s.cfg.SeedMasterTenantID,
			Name:   s.cfg.SeedMasterTenantName,
			Status: "ACTIVE",
		}).
		Exec(ctx)

	return err
}

func (s *MasterSeeder) ensureMasterUser(ctx context.Context, tx bun.Tx, password string) (identity.UserId, error) {
	var user models.UserModel
	err := tx.NewSelect().
		Model(&user).
		Where("tenant_id = ? AND email = ?", s.cfg.SeedMasterTenantID, s.cfg.SeedMasterEmail).
		Limit(1).
		Scan(ctx)
	if err == nil {
		return user.Id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}

	masterUserID := identity.UserId(s.cfg.SeedMasterUserID)

	_, err = tx.NewInsert().
		Model(&models.UserModel{
			Id:                   masterUserID,
			TenantId:             identity.TenantId(s.cfg.SeedMasterTenantID),
			IdentificationNumber: s.cfg.SeedMasterIdentificationNumber,
			Username:             s.cfg.SeedMasterUsername,
			Email:                s.cfg.SeedMasterEmail,
			Password:             hashedPassword,
			CreatedAt:            time.Now().UTC(),
		}).
		Exec(ctx)
	if err != nil {
		return "", err
	}

	return masterUserID, nil
}

func (s *MasterSeeder) ensureAdminRole(ctx context.Context, tx bun.Tx) (string, error) {
	const roleName = "admin"
	roleID := fmt.Sprintf("%s-admin-role", s.cfg.SeedMasterTenantID)

	var role models.RoleModel
	err := tx.NewSelect().
		Model(&role).
		Where("tenant_id = ? AND name = ?", s.cfg.SeedMasterTenantID, roleName).
		Limit(1).
		Scan(ctx)
	if err == nil {
		return role.Id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	_, err = tx.NewInsert().
		Model(&models.RoleModel{
			Id:          roleID,
			TenantId:    s.cfg.SeedMasterTenantID,
			Name:        roleName,
			Description: "System administrator with full access",
		}).
		Exec(ctx)
	if err != nil {
		return "", err
	}

	return roleID, nil
}

func (s *MasterSeeder) ensureGlobalAdminPermission(ctx context.Context, tx bun.Tx) (string, error) {
	const resource = "*"
	const action = "*"
	permissionID := fmt.Sprintf("%s-admin-all", s.cfg.SeedMasterTenantID)

	var permission models.PermissionModel
	err := tx.NewSelect().
		Model(&permission).
		Where("tenant_id = ? AND resource = ? AND action = ?", s.cfg.SeedMasterTenantID, resource, action).
		Limit(1).
		Scan(ctx)
	if err == nil {
		return permission.Id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	_, err = tx.NewInsert().
		Model(&models.PermissionModel{
			Id:       permissionID,
			TenantId: s.cfg.SeedMasterTenantID,
			Resource: resource,
			Action:   action,
		}).
		Exec(ctx)
	if err != nil {
		return "", err
	}

	return permissionID, nil
}

func (s *MasterSeeder) ensureRolePermission(ctx context.Context, tx bun.Tx, roleID, permissionID string) error {
	var rolePermission models.RolePermissionModel
	err := tx.NewSelect().
		Model(&rolePermission).
		Where("tenant_id = ? AND role_id = ? AND permission_id = ?", s.cfg.SeedMasterTenantID, roleID, permissionID).
		Limit(1).
		Scan(ctx)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = tx.NewInsert().
		Model(&models.RolePermissionModel{
			TenantId:     s.cfg.SeedMasterTenantID,
			RoleId:       roleID,
			PermissionId: permissionID,
		}).
		Exec(ctx)

	return err
}

func (s *MasterSeeder) ensureUserRole(ctx context.Context, tx bun.Tx, userID identity.UserId, roleID string) error {
	var userRole models.UserRoleModel
	err := tx.NewSelect().
		Model(&userRole).
		Where("tenant_id = ? AND user_id = ? AND role_id = ?", s.cfg.SeedMasterTenantID, userID, roleID).
		Limit(1).
		Scan(ctx)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = tx.NewInsert().
		Model(&models.UserRoleModel{
			TenantId: s.cfg.SeedMasterTenantID,
			UserId:   string(userID),
			RoleId:   roleID,
		}).
		Exec(ctx)

	return err
}
