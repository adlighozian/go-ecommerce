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

	time.Sleep(1 * time.Second)

	user, errGetByEmail := repo.GetByEmail(user.Email)
	if errGetByEmail != nil {
		return nil, errGetByEmail
	}

	return user, nil
}

func (repo *AuthRepository) FirstOrCreate(user *model.User) (*model.User, error) {
	foundUser, errGetByEmail := repo.GetByEmail(user.Email)
	if errGetByEmail != nil {
		if errors.Is(errGetByEmail, sql.ErrNoRows) {
			return repo.Create(user)
		}

		return nil, errGetByEmail
	}

	return foundUser, nil
}

func (repo *AuthRepository) GetByEmail(email string) (*model.User, error) {
	user := new(model.User)

	cachedData, errGetCache := repo.getDataFromCache(email)
	if errGetCache != nil {
		data, errGetDB := repo.getByEmailFromDatabase(email)
		if errGetDB != nil {
			return nil, errGetDB
		}

		dataByte, errJSON := json.Marshal(data)
		if errJSON != nil {
			return nil, errJSON
		}

		// Store the data in the cache for future reads
		errSetCache := repo.redis.Set(context.Background(), email, dataByte, 10*time.Minute).Err()
		if errSetCache != nil {
			return nil, errSetCache
		}

		return data, nil
	}

	errJSONUn := json.Unmarshal([]byte(cachedData), &user)
	if errJSONUn != nil {
		return nil, errJSONUn
	}

	return user, nil
}

func (repo *AuthRepository) getDataFromCache(key string) (string, error) {
	cachedData, errGet := repo.redis.Get(context.Background(), key).Result()
	if errGet != nil {
		return "", errGet
	}
	return cachedData, nil
}

func (repo *AuthRepository) getByEmailFromDatabase(email string) (*model.User, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := `
	SELECT users.id, users.username, users.email, users.role, users.full_name, 
	    	 users.age, users.image_url, users.created_at, users.updated_at, 
				 user_settings.notification, user_settings.dark_mode, 
				 languages.name AS language
	FROM users 
	INNER JOIN user_settings on users.id = user_settings.user_id
	INNER JOIN languages on user_settings.language_id= languages.id
	WHERE users.email = $1
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
		&user.ID, &user.Username, &user.Email, &user.Role, &user.FullName,
		&user.Age, &user.ImageURL, &user.CreatedAt, &user.UpdatedAt,
		&user.UserSetting.Notification, &user.UserSetting.DarkMode,
		&user.UserSetting.Language.Name,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return user, nil
}

func (repo *AuthRepository) SetRefreshToken(refreshToken string, dataByte []byte, refreshTokenDur time.Duration) error {
	errSetCache := repo.redis.Set(context.Background(), refreshToken, dataByte, refreshTokenDur).Err()
	if errSetCache != nil {
		return errSetCache
	}
	return nil
}
func (repo *AuthRepository) GetByRefreshToken(token string) (*model.RefreshToken, error) {
	refreshToken := new(model.RefreshToken)

	cachedData, errGetCache := repo.getDataFromCache(token)
	if errGetCache != nil {
		return nil, errGetCache
	}

	errJSONUn := json.Unmarshal([]byte(cachedData), &refreshToken)
	if errJSONUn != nil {
		return nil, errJSONUn
	}

	return refreshToken, nil
}
