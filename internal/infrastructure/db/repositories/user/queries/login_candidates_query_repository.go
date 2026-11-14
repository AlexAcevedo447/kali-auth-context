package queries

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type FindLoginCandidatesQueryRepository struct {
	db *bun.DB
}

var _ ports.IFindLoginCandidatesQueryRepository = (*FindLoginCandidatesQueryRepository)(nil)

func NewFindLoginCandidatesQueryRepository(db *bun.DB) *FindLoginCandidatesQueryRepository {
	return &FindLoginCandidatesQueryRepository{db: db}
}

type loginCandidateRow struct {
	Id                   string    `bun:"id"`
	TenantId             string    `bun:"tenant_id"`
	IdentificationNumber string    `bun:"identification_number"`
	Username             string    `bun:"username"`
	Email                string    `bun:"email"`
	Password             string    `bun:"password"`
	CreatedAt            time.Time `bun:"created_at"`
	TenantName           string    `bun:"tenant_name"`
	TenantCode           string    `bun:"tenant_code"`
	TenantStatus         string    `bun:"tenant_status"`
}

func (r *FindLoginCandidatesQueryRepository) FindByEmail(ctx context.Context, email, tenantCode string) ([]ports.LoginCandidate, error) {
	var rows []loginCandidateRow
	// JOIN con tenants: un usuario cuyo tenant no existe nunca puede autenticarse.
	q := r.db.NewSelect().
		TableExpr("user_models AS u").
		Join("JOIN tenants AS t ON t.id = u.tenant_id").
		ColumnExpr("u.id, u.tenant_id, u.identification_number, u.username, u.email, u.password, u.created_at").
		ColumnExpr("t.name AS tenant_name, COALESCE(t.code, '') AS tenant_code, t.status AS tenant_status").
		Where("u.email = ?", email).
		OrderExpr("u.tenant_id")
	if tenantCode != "" {
		q = q.Where("t.code = ?", tenantCode)
	}
	if err := q.Scan(ctx, &rows); err != nil {
		return nil, err
	}

	out := make([]ports.LoginCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, ports.LoginCandidate{
			User: &identity.User{
				Id:                   identity.UserId(row.Id),
				TenantId:             identity.TenantId(row.TenantId),
				IdentificationNumber: row.IdentificationNumber,
				Username:             row.Username,
				Email:                row.Email,
				Password:             row.Password,
			},
			TenantName:   row.TenantName,
			TenantCode:   row.TenantCode,
			TenantStatus: identity.TenantStatus(row.TenantStatus),
		})
	}
	return out, nil
}
