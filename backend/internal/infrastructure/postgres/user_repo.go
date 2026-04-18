package postgres

import (
    "backend/internal/entity"
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
)

type UserRepo struct {
    db *Postgres
}

func NewUserRepo(db *Postgres) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u *entity.User) error {
    err := r.db.Pool.QueryRow(ctx, "INSERT INTO users(name, email, password_hash) VALUES($1,$2,$3) RETURNING id", u.Name, u.Email, u.PasswordHash).Scan(&u.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    row := r.db.Pool.QueryRow(ctx, "SELECT id, name, email, password_hash FROM users WHERE email=$1", email)
    u := &entity.User{}
    if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash); err != nil {
        if err == pgx.ErrNoRows { return nil, fmt.Errorf("not found") }
        return nil, err
    }
    return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
    row := r.db.Pool.QueryRow(ctx, "SELECT id, name, email, password_hash FROM users WHERE id=$1", id)
    u := &entity.User{}
    if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash); err != nil {
        return nil, fmt.Errorf("not found")
    }
    return u, nil
}
