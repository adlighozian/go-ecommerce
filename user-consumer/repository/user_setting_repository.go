package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"user-consumer-go/helper/timeout"
	"user-consumer-go/model"
)

type UserSettingRepository struct {
	db *sql.DB
}

func NewUserSettingRepository(db *sql.DB) UserSettingRepositoryI {
	repo := new(UserSettingRepository)
	repo.db = db
	return repo
}

func (repo *UserSettingRepository) UpdateByUserID(setting *model.UserSetting) (*model.UserSetting, error) {
	// Prepare the SQL statement
	var columns []string
	var args []interface{}
	var argPos = 1

	if setting.Notification != nil {
		columns = append(columns, fmt.Sprintf("notification = $%d", argPos))
		args = append(args, *setting.Notification)
		argPos++
	}

	if setting.DarkMode != nil {
		columns = append(columns, fmt.Sprintf("dark_mode = $%d", argPos))
		args = append(args, *setting.DarkMode)
		argPos++
	}

	if setting.LanguageID != 0 {
		columns = append(columns, fmt.Sprintf("language_id = $%d", argPos))
		args = append(args, setting.LanguageID)
		argPos++
	}

	if len(columns) == 0 {
		return setting, nil // no update needed
	}

	// Append the user ID at the end
	args = append(args, setting.ID)
	//nolint:gosec // hard to avoid this
	query := fmt.Sprintf(
		`
		UPDATE user_settings SET %s 
		WHERE user_id = $%d 
		RETURNING id, user_id, notification, dark_mode, languange_id
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
		&setting.ID, &setting.UserID, &setting.Notification, &setting.DarkMode, &setting.LanguageID,
	)
	if errScan != nil {
		return nil, errScan
	}

	return setting, nil
}
