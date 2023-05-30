package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
	"user-go/helper/timeout"
	"user-go/model"
	"user-go/package/rmq"

	"github.com/redis/go-redis/v9"
)

type UserSettingRepository struct {
	db    *sql.DB
	redis *redis.Client
	rmq   rmq.RabbitMQClient
}

func NewUserSettingRepository(db *sql.DB, redis *redis.Client, rmq rmq.RabbitMQClient) UserSettingRepositoryI {
	repo := new(UserSettingRepository)
	repo.db = db
	repo.redis = redis
	repo.rmq = rmq
	return repo
}

func (repo *UserSettingRepository) UpdateByUserID(newSetting *model.UserSetting) (*model.User, error) {
	userID := newSetting.UserID

	errCache := repo.redis.Del(context.Background(), strconv.FormatUint(uint64(userID), 10)).Err()
	if errCache != nil {
		return nil, errCache
	}

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	newSettingBytes, errJSON := json.Marshal(newSetting)
	if errJSON != nil {
		return nil, errJSON
	}

	errPub := repo.rmq.Publish(
		ctx,
		"user_setting.updated",
		"topic",
		"application/json",
		"user_setting.updated",
		newSettingBytes,
	)
	if errPub != nil {
		return nil, errPub
	}

	time.Sleep(1 * time.Second)

	user, errGetByID := repo.getByID(userID)
	if errGetByID != nil {
		return nil, errGetByID
	}

	return user, nil
}

func (repo *UserSettingRepository) getByID(userID uint) (*model.User, error) {
	user := new(model.User)

	cachedData, errGetCache := repo.getDataFromCache(strconv.FormatUint(uint64(userID), 10))
	if errGetCache != nil {
		data, errGetDB := repo.getByUserIDFromDatabase(userID)
		if errGetDB != nil {
			return nil, errGetDB
		}

		dataByte, errJSON := json.Marshal(data)
		if errJSON != nil {
			return nil, errJSON
		}

		// Store the data in the cache for future reads
		errSetCache := repo.redis.Set(
			context.Background(),
			strconv.FormatUint(uint64(userID), 10), dataByte, 10*time.Minute,
		).Err()
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

func (repo *UserSettingRepository) getDataFromCache(key string) (string, error) {
	cachedData, errGet := repo.redis.Get(context.Background(), key).Result()
	if errGet != nil {
		return "", errGet
	}
	return cachedData, nil
}

func (repo *UserSettingRepository) getByUserIDFromDatabase(userID uint) (*model.User, error) {
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
	WHERE users.id = $1
	LIMIT 1
	`
	stmt, errStmt := repo.db.PrepareContext(ctx, sqlQuery)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	user := new(model.User)
	row := stmt.QueryRowContext(ctx, userID)
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
