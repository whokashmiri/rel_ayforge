package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whokashmiri/rel_ayforge/internal/domain"
)

var ErrNotFound = errors.New("record not found")

type TenantRepository struct {
	db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{db: db}
}

func (repo *TenantRepository) Create(ctx context.Context, name string) (domain.Tenant, error) {
	const query = `
	INSERT INTO tenants (id , name)
	VALUES (gen_random_uuid(), $1)
	RETURNING id::text , name, created_at, updated_at
	`
	var tenant domain.Tenant
	err := repo.db.QueryRow(ctx, query, name).Scan(
		&tenant.ID, &tenant.Name, &tenant.CreateAt, &tenant.UpdatedAt,
	)
	if err != nil {
		return domain.Tenant{}, fmt.Errorf("insert tenant : %w", err)
	}
	return tenant, nil
}
func (repo *TenantRepository) GetByID(ctx context.Context, id string) (domain.Tenant, error) {
	const query = `
		SELECT id::text , name , created_at , updated_at
		FROM tenants
		WHERE id = $1::uuid
		`
	var tenant domain.Tenant
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&tenant.ID, &tenant.Name, &tenant.CreateAt, &tenant.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Tenant{}, ErrNotFound
	}

	if err != nil {
		return domain.Tenant{}, fmt.Errorf("select tenant : %w", err)
	}
	return tenant, nil
}
