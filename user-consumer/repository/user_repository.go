package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"user-consumer-go/helper/timeout"
	"user-consumer-go/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryI {
	repo := new(UserRepository)
	repo.db = db
	return repo
}

func (repo *UserRepository) Create(user *model.User) (*model.User, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	tx, errTx := repo.db.BeginTx(ctx, nil)
	if errTx != nil {
		return nil, errTx
	}

	sqlQuery1 := `
	INSERT INTO users (username, email, password, role, full_name, age, image_url) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id
	`
	stmt1, errStmt := tx.PrepareContext(ctx, sqlQuery1)
	if errStmt != nil {
		return nil, errStmt
	}
	defer stmt1.Close()

	errScan := stmt1.QueryRowContext(ctx,
		user.Username, user.Email, user.Password, user.Role,
		user.FullName, user.Age, user.ImageURL,
	).
		Scan(&user.ID)
	if errScan != nil {
		return nil, errScan
	}

	sqlQuery2 := `
	INSERT INTO user_settings (user_id) VALUES ($1) 
	`
	stmt2, errStmt2 := tx.PrepareContext(ctx, sqlQuery2)
	if errStmt2 != nil {
		return nil, errStmt2
	}
	defer stmt2.Close()

	_, errExec := stmt2.ExecContext(ctx, user.ID)
	if errExec != nil {
		return nil, errExec
	}

	_ = tx.Commit()

	return user, nil
}

func (repo *UserRepository) UpdateByID(user *model.User) (*model.User, error) {
	// Prepare the SQL statement
	var columns []string
	var args []interface{}
	var argPos = 1

	if user.Username != "" {
		columns = append(columns, fmt.Sprintf("username = $%d", argPos))
		args = append(args, user.Username)
		argPos++
	}

	if user.FullName != "" {
		columns = append(columns, fmt.Sprintf("full_name = $%d", argPos))
		args = append(args, user.FullName)
		argPos++
	}

	if user.Age != 0 {
		columns = append(columns, fmt.Sprintf("age = $%d", argPos))
		args = append(args, user.Age)
		argPos++
	}

	if user.ImageURL != "" {
		columns = append(columns, fmt.Sprintf("image_url = $%d", argPos))
		args = append(args, user.ImageURL)
		argPos++
	}

	if len(columns) == 0 {
		return user, nil // no update needed
	}

	// Append the user ID at the end
	args = append(args, user.ID)
	//nolint:gosec // hard to avoid this
	query := fmt.Sprintf(
		`
	UPDATE users SET %s 
	WHERE id = $%d 
	RETURNING id, username, full_name, age, image_url 
	`, strings.Join(columns, ", "), argPos,
	)

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	stmt, errStmt2 := repo.db.PrepareContext(ctx, query)
	if errStmt2 != nil {
		return nil, errStmt2
	}
	defer stmt.Close()

	errScan := stmt.QueryRowContext(ctx, args...).Scan(
		&user.ID, &user.Username, &user.FullName, &user.Age, &user.ImageURL,
	)
	if errScan != nil {
		return nil, errScan
	}

	return user, nil
}
