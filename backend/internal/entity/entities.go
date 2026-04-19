package entity

import "time"

type Place struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Address      string `json:"address"`
	Description  string `json:"description,omitempty"`
	Image        string `json:"image,omitempty"`
	PricePerHour int    `json:"price_per_hour"`
}

type Slot struct {
	ID          int       `json:"id"`
	PlaceID     int       `json:"place_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsAvailable bool      `json:"is_available"`
}

type Booking struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id"`
	PlaceID   int       `json:"place_id"`
	SlotID    int       `json:"slot_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type User struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}
