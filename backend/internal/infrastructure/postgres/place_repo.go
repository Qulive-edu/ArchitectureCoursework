package postgres

import (
	"backend/internal/entity"
	"context"
	"fmt"
)

type PlaceRepo struct {
	db *Postgres
}

// ListSlots implements usecase.PlaceRepo.
func (r *PlaceRepo) ListSlots(ctx context.Context, placeID int) ([]*entity.Slot, error) {
	return NewSlotRepo(r.db).ListByPlace(ctx, placeID)
}

func NewPlaceRepo(db *Postgres) *PlaceRepo { return &PlaceRepo{db: db} }

func (r *PlaceRepo) List(ctx context.Context) ([]*entity.Place, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT id, title, address, floor_type, description, image, price_per_hour FROM places")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*entity.Place
	for rows.Next() {
		p := &entity.Place{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Address, &p.FloorType, &p.Description, &p.Image, &p.PricePerHour); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (r *PlaceRepo) GetByID(ctx context.Context, id int) (*entity.Place, error) {
	row := r.db.Pool.QueryRow(ctx, "SELECT id, title, address, floor_type, description, image, price_per_hour FROM places WHERE id=$1", id)
	p := &entity.Place{}
	if err := row.Scan(&p.ID, &p.Title, &p.Address, &p.FloorType, &p.Description, &p.Image, &p.PricePerHour); err != nil {
		return nil, fmt.Errorf("not found")
	}
	return p, nil
}
