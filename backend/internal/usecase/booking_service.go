package usecase

import (
    "backend/internal/entity"
    "context"
    "errors"
)

type BookingRepo interface {
    Create(ctx context.Context, b *entity.Booking) error
    ListByUser(ctx context.Context, userID int) ([]*entity.Booking, error)
}

type SlotRepo interface {
    GetForUpdate(ctx context.Context, id int) (*entity.Slot, error)
    MarkUnavailable(ctx context.Context, id int) error
}

type BookingService interface {
    CreateBooking(ctx context.Context, userID, placeID, slotID int) error
    ListMyBookings(ctx context.Context, userID int) ([]*entity.Booking, error)
}

type bookingService struct {
    bookingRepo BookingRepo
    slotRepo SlotRepo
}

func NewBookingService(br BookingRepo, sr SlotRepo) BookingService {
    return &bookingService{bookingRepo: br, slotRepo: sr}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID, placeID, slotID int) error {
    // 1. check slot
    slot, err := s.slotRepo.GetForUpdate(ctx, slotID)
    if err != nil {
        return err
    }
    if !slot.IsAvailable {
        return errors.New("slot not available")
    }
    // mark unavailable
    if err := s.slotRepo.MarkUnavailable(ctx, slotID); err != nil {
        return err
    }
    b := &entity.Booking{UserID: userID, PlaceID: placeID, SlotID: slotID, Status: "confirmed"}
    if err := s.bookingRepo.Create(ctx, b); err != nil {
        return err
    }
    return nil
}

func (s *bookingService) ListMyBookings(ctx context.Context, userID int) ([]*entity.Booking, error) {
    return s.bookingRepo.ListByUser(ctx, userID)
}
