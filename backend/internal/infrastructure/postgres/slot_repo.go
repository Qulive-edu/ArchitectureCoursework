package postgres

import (
	"backend/internal/entity"
	"context"
	"fmt"
)

type SlotRepo struct {
	db *Postgres
}

func NewSlotRepo(db *Postgres) *SlotRepo { return &SlotRepo{db: db} }

func (r *SlotRepo) ListByPlace(ctx context.Context, placeID int) ([]*entity.Slot, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT id, place_id, start_time, end_time, is_available FROM time_slots WHERE place_id=$1 AND start_time > NOW()", placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*entity.Slot
	for rows.Next() {
		s := &entity.Slot{}
		if err := rows.Scan(&s.ID, &s.PlaceID, &s.StartTime, &s.EndTime, &s.IsAvailable); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (r *SlotRepo) GetForUpdate(ctx context.Context, id int) (*entity.Slot, error) {
	row := r.db.Pool.QueryRow(ctx, "SELECT id, place_id, start_time, end_time, is_available FROM time_slots WHERE id=$1", id)
	s := &entity.Slot{}
	if err := row.Scan(&s.ID, &s.PlaceID, &s.StartTime, &s.EndTime, &s.IsAvailable); err != nil {
		return nil, fmt.Errorf("slot not found")
	}
	return s, nil
}

func (r *SlotRepo) MarkUnavailable(ctx context.Context, id int) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE time_slots SET is_available=false WHERE id=$1", id)
	return err
}
