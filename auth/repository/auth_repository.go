package repository

import (
	"auth-go/helper/timeout"
	"auth-go/model"
	"auth-go/package/rmq"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	db    *sql.DB
	redis *redis.Client
	rmq   rmq.RabbitMQClient
}

func NewAuthRepository(db *sql.DB, redis *redis.Client, rmq rmq.RabbitMQClient) AuthRepositoryI {
	repo := new(AuthRepository)
	repo.db = db
	repo.redis = redis
	repo.rmq = rmq
	return repo
}

func (repo *AuthRepository) isEmailUnique(ctx context.Context, user *model.User) error {
	sqlQuery := `
	SELECT COUNT(*) FROM users WHERE email = $1
	`
	stmt, errStmt := repo.db.PrepareContext(ctx, sqlQuery)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()

	count := 0
	errScan := stmt.QueryRowContext(ctx, user.Email).Scan(&count)
	if errScan != nil {
		return errScan
	}

	if count > 0 {
		return errors.New("duplicated key not allowed")
	}
	return nil
}

func (repo *AuthRepository) Create(user *model.User) (*model.User, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	// check if an email is unique before inserting
	errEmailUnique := repo.isEmailUnique(ctx, user)
	if errEmailUnique != nil {
		return nil, errEmailUnique
	}

	userBytes, errJSON := json.Marshal(user)
	if errJSON != nil {
		return nil, errJSON
	}

	errPub := repo.rmq.Publish(
		ctx,
		"user.created",
		"topic",
		"application/json",
		"user.created",
		userBytes,
	)
	if errPub != nil {
		return nil, errPub
	}

	time.Sleep(3 * time.Second)

	return repo.LoginByEmail(user.Email)
}

func (repo *AuthRepository) FirstOrCreate(user *model.User) (*model.User, error) {
	foundUser, errGetByEmail := repo.LoginByEmail(user.Email)
	if errGetByEmail != nil {
		if errors.Is(errGetByEmail, sql.ErrNoRows) {
			return repo.Create(user)
		}

		return nil, errGetByEmail
	}

	return foundUser, nil
}

func (repo *AuthRepository) getDataFromCache(key string) (string, error) {
	return repo.redis.Get(context.Background(), key).Result()
}

func (repo *AuthRepository) LoginByEmail(email string) (*model.User, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := `
	SELECT id, username, email, password, role, provider,
	    	  full_name, age, image_url, created_at, updated_at
	FROM users 
	WHERE email = $1
	LIMIT 1
	`
	stmt, errStmt := repo.db.PrepareContext(ctx, sqlQuery)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	user := new(model.User)
	row := stmt.QueryRowContext(ctx, email)
	scanErr := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.Provider,
		&user.FullName, &user.Age, &user.ImageURL,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return user, nil
}

func (repo *AuthRepository) SetRefreshToken(token string, dataByte []byte, refreshTokenDur time.Duration) error {
	return repo.redis.Set(context.Background(), "refresh_token:"+token, dataByte, refreshTokenDur).Err()
}

func (repo *AuthRepository) GetByRefreshToken(token string) (*model.RefreshToken, error) {
	refreshToken := new(model.RefreshToken)

	cachedData, errGetCache := repo.getDataFromCache("refresh_token:" + token)
	if errGetCache != nil {
		return nil, errGetCache
	}

	errJSONUn := json.Unmarshal([]byte(cachedData), &refreshToken)
	if errJSONUn != nil {
		return nil, errJSONUn
	}

	return refreshToken, nil
}
