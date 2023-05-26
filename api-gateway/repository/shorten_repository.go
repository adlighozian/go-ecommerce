package repository

import (
	"api-gateway-go/helper/timeout"
	"api-gateway-go/model"
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type ShortenRepo struct {
	db    *sql.DB
	redis *redis.Client
}

func NewShortenRepo(db *sql.DB, redis *redis.Client) ShortenRepoI {
	repo := new(ShortenRepo)
	repo.db = db
	repo.redis = redis
	return repo
}

func (repo *ShortenRepo) Get(hashedURL string) (*model.APIManagement, error) {
	apiManagement := new(model.APIManagement)

	// note: implemetation of cache aside of read cache strategy
	cachedData, errGetCache := repo.getDataFromCache(hashedURL)
	if errGetCache != nil {
		data, errGetDB := repo.getHashURLFromDatabase(hashedURL)
		if errGetDB != nil {
			return nil, errGetDB
		}

		dataByte, errJSON := json.Marshal(data)
		if errJSON != nil {
			return nil, errJSON
		}
		// Store the data in the cache for future reads
		errSetCache := repo.redis.Set(context.Background(), hashedURL, dataByte, 10*time.Minute).Err()
		if errSetCache != nil {
			return nil, errSetCache
		}

		return data, nil
	}

	errJSONUn := json.Unmarshal([]byte(cachedData), &apiManagement)
	if errJSONUn != nil {
		return nil, errJSONUn
	}

	return apiManagement, nil
}

func (repo *ShortenRepo) getDataFromCache(key string) (string, error) {
	cachedData, errGet := repo.redis.Get(context.Background(), key).Result()
	if errGet != nil {
		return "", errGet
	}
	return cachedData, nil
}

func (repo *ShortenRepo) getHashURLFromDatabase(hashedURL string) (*model.APIManagement, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := `
	SELECT id, api_name, service_name, endpoint_url, 
				 hashed_endpoint_url, is_available, created_at, updated_at 
	FROM api_managements 
	WHERE hashed_endpoint_url = $1 AND is_available = TRUE 
	LIMIT 1
	`
	stmt, errStmt := repo.db.PrepareContext(ctx, sqlQuery)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	apiManagement := new(model.APIManagement)
	row := stmt.QueryRowContext(ctx, hashedURL)
	scanErr := row.Scan(
		&apiManagement.ID, &apiManagement.APIName, &apiManagement.ServiceName, &apiManagement.EndpointURL,
		&apiManagement.HashedEndpointURL, &apiManagement.IsAvailable, &apiManagement.CreatedAt, &apiManagement.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return apiManagement, nil
}

func (repo *ShortenRepo) Create(apiManagement *model.APIManagement) (*model.APIManagement, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	sqlQuery := `
	INSERT INTO api_managements (api_name, service_name, endpoint_url, hashed_endpoint_url, is_available)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	stmt, errStmt := repo.db.PrepareContext(ctx, sqlQuery)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt.Close()

	// apiManagement := new(model.APIManagement)
	row := stmt.QueryRowContext(
		ctx,
		&apiManagement.APIName, &apiManagement.ServiceName,
		&apiManagement.EndpointURL, &apiManagement.HashedEndpointURL,
		&apiManagement.IsAvailable,
	)
	scanErr := row.Scan(
		&apiManagement.ID, &apiManagement.CreatedAt, &apiManagement.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	return apiManagement, nil
}
