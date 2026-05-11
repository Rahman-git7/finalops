package repository

import (
	"context"
	"time"

	"github.com/finalops/auth-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO auth.users (id, email, username, password, created_at, updated_at)
         VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		u.ID, u.Email, u.Username, u.Password,
	)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, password, created_at, updated_at
         FROM auth.users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, password, created_at, updated_at
         FROM auth.users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO auth.refresh_tokens (user_id, token_hash, expires_at)
         VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	return err
}

func (r *UserRepository) FindRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.db.QueryRow(ctx,
		`SELECT user_id FROM auth.refresh_tokens
         WHERE token_hash = $1 AND expires_at > NOW()`, tokenHash,
	).Scan(&userID)
	return userID, err
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM auth.refresh_tokens WHERE token_hash = $1`, tokenHash,
	)
	return err
}

func (r *UserRepository) GetProfile(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	p := &model.Profile{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, weight_kg, age, sex, created_at, updated_at
         FROM auth.profiles WHERE user_id = $1`, userID,
	).Scan(&p.ID, &p.UserID, &p.WeightKg, &p.Age, &p.Sex, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *UserRepository) UpsertProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.Profile, error) {
	p := &model.Profile{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO auth.profiles (user_id, weight_kg, age, sex, created_at, updated_at)
         VALUES ($1, $2, $3, $4, NOW(), NOW())
         ON CONFLICT (user_id) DO UPDATE
         SET weight_kg = EXCLUDED.weight_kg,
             age = EXCLUDED.age,
             sex = EXCLUDED.sex,
             updated_at = NOW()
         RETURNING id, user_id, weight_kg, age, sex, created_at, updated_at`,
		userID, req.WeightKg, req.Age, req.Sex,
	).Scan(&p.ID, &p.UserID, &p.WeightKg, &p.Age, &p.Sex, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}
