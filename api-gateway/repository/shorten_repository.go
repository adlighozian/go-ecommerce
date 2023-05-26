package repository

import (
	"api-gateway-go/helper/timeout"
	"api-gateway-go/model"
	"database/sql"
)

type ShortenRepo struct {
	db *sql.DB
}

func NewShortenRepo(db *sql.DB) ShortenRepoI {
	repo := new(ShortenRepo)
	repo.db = db
	return repo
}

func (repo *ShortenRepo) Get(hashedURL string) (*model.APIManagement, error) {
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
