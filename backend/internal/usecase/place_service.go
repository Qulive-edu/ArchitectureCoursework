package usecase

import (
    "backend/internal/entity"
    "context"
    "strconv"
)

type PlaceRepo interface {
    List(ctx context.Context) ([]*entity.Place, error)
    GetByID(ctx context.Context, id int) (*entity.Place, error)
    ListSlots(ctx context.Context, placeID int) ([]*entity.Slot, error)
}

type PlaceService interface {
    ListPlaces(ctx context.Context) ([]*entity.Place, error)
    GetPlaceByID(ctx context.Context, id string) (*entity.Place, error)
    ListSlots(ctx context.Context, placeID string) ([]*entity.Slot, error)
}

type placeService struct {
    repo PlaceRepo
}

func NewPlaceService(r PlaceRepo, slotRepo interface{}) PlaceService {
    return &placeService{repo: r}
}

func (s *placeService) ListPlaces(ctx context.Context) ([]*entity.Place, error) {
    return s.repo.List(ctx)
}

func (s *placeService) GetPlaceByID(ctx context.Context, id string) (*entity.Place, error) {
    i, _ := strconv.Atoi(id)
    return s.repo.GetByID(ctx, i)
}

func (s *placeService) ListSlots(ctx context.Context, placeID string) ([]*entity.Slot, error) {
    i, _ := strconv.Atoi(placeID)
    return s.repo.ListSlots(ctx, i)
}
