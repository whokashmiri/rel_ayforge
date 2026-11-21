package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/whokashmiri/rel_ayforge/internal/domain"
	"github.com/whokashmiri/rel_ayforge/internal/repository"
)

var (
	ErrInvalidName    = errors.New("invalid tenant name")
	ErrInvalidID      = errors.New("invalid tenant ID")
	ErrTenantNotFound = errors.New("tenant not found")
)

var uuidPattern = regexp.MustCompile(
	`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
)

type TenantStore interface {
	Create(ctx context.Context, name string) (domain.Tenant, error)
	GetByID(ctx context.Context, id string) (domain.Tenant, error)
}

type TenantService struct {
	repo TenantStore
}

func NewTenantService(repo TenantStore) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) Register(ctx context.Context, name string) (domain.Tenant, error) {
	name = strings.TrimSpace(name)

	if name == "" || utf8.RuneCountInString(name) > 150 {
		return domain.Tenant{}, ErrInvalidName
	}
	tenant, err := s.repo.Create(ctx, name)

	if err != nil {
		return domain.Tenant{}, fmt.Errorf("register tenant : %w", err)
	}
	return tenant, nil
}

func (s *TenantService) GetByID(ctx context.Context, id string) (domain.Tenant, error) {
	if !uuidPattern.MatchString(id) {
		return domain.Tenant{}, ErrInvalidID
	}
	tenant, err := s.repo.GetByID(ctx, id)

	if errors.Is(err, repository.ErrNotFound) {
		return domain.Tenant{}, ErrTenantNotFound
	}
	if err != nil {
		return domain.Tenant{}, fmt.Errorf("get tenant: %w", err)
	}
	return tenant, nil
}
