package usecase

import (
    "backend/internal/entity"
    "context"
    "errors"
    "golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
    Create(ctx context.Context, u *entity.User) error
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    GetByID(ctx context.Context, id int) (*entity.User, error)
}

type UserService interface {
    Register(ctx context.Context, name, email, password string) (*entity.User, error)
    Login(ctx context.Context, email, password string) (*entity.User, error)
    GetByID(ctx context.Context, id int) (*entity.User, error)
}

type userService struct {
    repo UserRepo
}

func NewUserService(r UserRepo) UserService { return &userService{repo: r} }

func (s *userService) Register(ctx context.Context, name, email, password string) (*entity.User, error) {
    h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    u := &entity.User{Name: name, Email: email, PasswordHash: string(h)}
    if err := s.repo.Create(ctx, u); err != nil {
        return nil, err
    }
    return u, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (*entity.User, error) {
    u, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    return u, nil
}

func (s *userService) GetByID(ctx context.Context, id int) (*entity.User, error) {
    return s.repo.GetByID(ctx, id)
}
