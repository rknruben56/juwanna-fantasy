package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rknruben56/juwanna-fantasy/api/internal/domain"
)

type MappingRepo struct {
	pool *pgxpool.Pool
}

func NewMappingRepo(pool *pgxpool.Pool) *MappingRepo {
	return &MappingRepo{pool: pool}
}

func (r *MappingRepo) List(ctx context.Context) ([]domain.OwnerSleeperMappingWithName, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT m.id, m.owner_id, o.name, m.sleeper_user_id
		FROM owner_sleeper_mapping m
		JOIN owners o ON m.owner_id = o.id
		ORDER BY o.name
	`)
	if err != nil {
		return nil, fmt.Errorf("list mappings: %w", err)
	}
	defer rows.Close()

	var mappings []domain.OwnerSleeperMappingWithName
	for rows.Next() {
		var m domain.OwnerSleeperMappingWithName
		if err := rows.Scan(&m.ID, &m.OwnerID, &m.OwnerName, &m.SleeperUserID); err != nil {
			return nil, fmt.Errorf("scan mapping: %w", err)
		}
		mappings = append(mappings, m)
	}
	return mappings, nil
}

func (r *MappingRepo) GetByOwnerID(ctx context.Context, ownerID int) (*domain.OwnerSleeperMapping, error) {
	var m domain.OwnerSleeperMapping
	err := r.pool.QueryRow(ctx, `
		SELECT id, owner_id, sleeper_user_id, created_at, updated_at
		FROM owner_sleeper_mapping
		WHERE owner_id = $1
	`, ownerID).Scan(&m.ID, &m.OwnerID, &m.SleeperUserID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get mapping by owner: %w", err)
	}
	return &m, nil
}

func (r *MappingRepo) GetBySleeperUserID(ctx context.Context, sleeperUserID string) (*domain.OwnerSleeperMapping, error) {
	var m domain.OwnerSleeperMapping
	err := r.pool.QueryRow(ctx, `
		SELECT id, owner_id, sleeper_user_id, created_at, updated_at
		FROM owner_sleeper_mapping
		WHERE sleeper_user_id = $1
	`, sleeperUserID).Scan(&m.ID, &m.OwnerID, &m.SleeperUserID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get mapping by sleeper user: %w", err)
	}
	return &m, nil
}

func (r *MappingRepo) Create(ctx context.Context, ownerID int, sleeperUserID string) (*domain.OwnerSleeperMapping, error) {
	var m domain.OwnerSleeperMapping
	err := r.pool.QueryRow(ctx, `
		INSERT INTO owner_sleeper_mapping (owner_id, sleeper_user_id)
		VALUES ($1, $2)
		RETURNING id, owner_id, sleeper_user_id, created_at, updated_at
	`, ownerID, sleeperUserID).Scan(&m.ID, &m.OwnerID, &m.SleeperUserID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create mapping: %w", err)
	}
	return &m, nil
}

func (r *MappingRepo) Update(ctx context.Context, id int, sleeperUserID string) (*domain.OwnerSleeperMapping, error) {
	var m domain.OwnerSleeperMapping
	err := r.pool.QueryRow(ctx, `
		UPDATE owner_sleeper_mapping
		SET sleeper_user_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, owner_id, sleeper_user_id, created_at, updated_at
	`, sleeperUserID, id).Scan(&m.ID, &m.OwnerID, &m.SleeperUserID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update mapping: %w", err)
	}
	return &m, nil
}

func (r *MappingRepo) Delete(ctx context.Context, id int) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM owner_sleeper_mapping WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete mapping: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("mapping not found")
	}
	return nil
}
