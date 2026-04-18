package postgres

import (
	"backend/internal/entity"
	"context"
)

type BookingRepo struct {
	db *Postgres
}

func NewBookingRepo(db *Postgres) *BookingRepo { return &BookingRepo{db: db} }

func (r *BookingRepo) Create(ctx context.Context, b *entity.Booking) error {
	err := r.db.Pool.QueryRow(ctx,
		"INSERT INTO bookings(user_id, place_id, slot_id, status) VALUES($1,$2,$3,$4) RETURNING id",
		b.UserID, b.PlaceID, b.SlotID, b.Status,
	).Scan(&b.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BookingRepo) ListByUser(ctx context.Context, userID int) ([]*entity.Booking, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT id, user_id, place_id, slot_id, status, created_at FROM bookings WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*entity.Booking
	for rows.Next() {
		b := &entity.Booking{}
		if err := rows.Scan(&b.ID, &b.UserID, &b.PlaceID, &b.SlotID, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, nil
}
