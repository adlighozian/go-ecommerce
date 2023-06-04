package repository_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"
	"user-go/mocks"
	"user-go/model"
	"user-go/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepoSuite struct {
	suite.Suite
	mockDB      *sql.DB
	mockSQL     sqlmock.Sqlmock
	mockRedisDB *redis.Client
	mockRedis   redismock.ClientMock
	mockRmq     *mocks.RabbitMQClient
	repo        repository.UserRepositoryI
}

func (s *UserRepoSuite) SetupTest() {
	var err error

	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Require().NoError(err)
	}

	// * gorm.Config handle internally, which can not mock explisitly
	gormConf := new(gorm.Config)
	gormConf.Logger = logger.Default.LogMode(logger.Info)
	gormConf.PrepareStmt = true
	gormConf.SkipDefaultTransaction = true

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, gormConf)
	if err != nil {
		s.Require().NoError(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		s.Require().NoError(err)
	}

	redisDB, mockRedis := redismock.NewClientMock()

	rmq := mocks.NewRabbitMQClient(s.T())

	repo := repository.NewUserRepository(sqlDB, redisDB, rmq)

	s.mockDB = sqlDB
	s.mockSQL = mockSQL
	s.mockRedisDB = redisDB
	s.mockRedis = mockRedis
	s.mockRmq = rmq
	s.repo = repo
}

func (s *UserRepoSuite) TearDownTest() {
	if errMockSQL := s.mockSQL.ExpectationsWereMet(); errMockSQL != nil {
		s.Require().NoError(errMockSQL)
	}

	if errMockRedis := s.mockRedis.ExpectationsWereMet(); errMockRedis != nil {
		s.Require().NoError(errMockRedis)
	}

	defer func() {
		s.mockDB.Close()
		s.mockRedisDB.Close()
	}()
}

func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserRepoSuite))
}

func (s *UserRepoSuite) TestUserRepository_GetByID() {
	notif, darkmode := true, false

	user := model.User{
		ID:        1,
		Username:  "test",
		Email:     "test@test.com",
		Role:      "user",
		Provider:  "email",
		FullName:  "testtist",
		Age:       21, // i swear
		ImageURL:  "some-img",
		CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
		UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),

		UserSetting: model.UserSetting{
			ID:           0,
			UserID:       0,
			Notification: &notif,
			DarkMode:     &darkmode,
			LanguageID:   0,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},

			Language: model.Language{
				ID:        0,
				Name:      "English",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
	}

	type args struct {
		userID uint
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, redismock.ClientMock)
		want       *model.User
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success - redis",
			args: args{
				userID: user.ID,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock) {
				data, _ := json.Marshal(user)
				cm.ExpectGet("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetVal(string(data))
			},
			want: &user,
		},
		{
			name: "success - db",
			args: args{
				userID: user.ID,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock) {
				cm.ExpectGet("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetErr(errors.New("redis error"))

				row := s.NewRows(
					[]string{
						"id", "username", "email", "role", "provider",
						"full_name", "age", "image_url", "created_at", "updated_at",
						"notification", "dark_mode",
						"language",
					}).
					AddRow(
						user.ID, user.Username, user.Email, user.Role, user.Provider,
						user.FullName, user.Age, user.ImageURL, user.CreatedAt, user.UpdatedAt,
						user.UserSetting.Notification, user.UserSetting.DarkMode,
						user.UserSetting.Language.Name,
					)

				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.ID).
					WillReturnRows(row)

				data, _ := json.Marshal(user)
				cm.ExpectSet("user_id:"+strconv.FormatUint(uint64(user.ID), 10), data, 10*time.Minute).SetVal("OK")
			},
			want: &user,
		},
		{
			name: "failed - db - stmt",
			args: args{
				userID: user.ID,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock) {
				cm.ExpectGet("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetErr(errors.New("redis error"))

				s.ExpectPrepare("SELECT .* FROM users .*").
					WillReturnError(errors.New("prepare stmt error"))
			},
			wantErr: true,
		},
		{
			name: "failed - db - scan",
			args: args{
				userID: user.ID,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock) {
				cm.ExpectGet("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetErr(errors.New("redis error"))

				row := s.NewRows(
					[]string{
						"id", "username", "email", "role", "provider",
						"full_name", "age", "image_url", "created_at", "updated_at",
						"notification", "dark_mode",
						"language",
					}).
					AddRow(
						user.ID, user.Username, user.Email, user.Role, user.Provider,
						nil, user.Age, user.ImageURL, user.CreatedAt, user.UpdatedAt,
						user.UserSetting.Notification, user.UserSetting.DarkMode,
						user.UserSetting.Language.Name,
					).
					RowError(1, errors.New("scan error"))

				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.ID).
					WillReturnRows(row)
			},
			wantErr: true,
		},
		{
			name: "success - db",
			args: args{
				userID: user.ID,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock) {
				cm.ExpectGet("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetErr(errors.New("redis error"))

				row := s.NewRows(
					[]string{
						"id", "username", "email", "role", "provider",
						"full_name", "age", "image_url", "created_at", "updated_at",
						"notification", "dark_mode",
						"language",
					}).
					AddRow(
						user.ID, user.Username, user.Email, user.Role, user.Provider,
						user.FullName, user.Age, user.ImageURL, user.CreatedAt, user.UpdatedAt,
						user.UserSetting.Notification, user.UserSetting.DarkMode,
						user.UserSetting.Language.Name,
					)

				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.ID).
					WillReturnRows(row)

				data, _ := json.Marshal(user)
				cm.ExpectSet("user_id:"+strconv.FormatUint(uint64(user.ID), 10), data, 10*time.Minute).
					SetErr(errors.New("redis error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, s.mockRedis)
			}
			got, err := s.repo.GetByID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("UserRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				s.T().Errorf("UserRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *UserRepoSuite) TestUserRepository_UpdateByID() {
	user := model.User{
		ID:        1,
		Username:  "test",
		Email:     "test@test.com",
		Role:      "user",
		Provider:  "email",
		FullName:  "testtist",
		Age:       21, // i swear
		ImageURL:  "some-img",
		CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
		UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
	}
	type args struct {
		profile *model.User
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, redismock.ClientMock, *mocks.RabbitMQClient)
		want       *model.User
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				profile: &user,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				cm.ExpectDel("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetVal(1)

				// publish
				rm.On(
					"Publish",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(nil)

				row := s.NewRows(
					[]string{
						"id", "username", "email", "role", "provider",
						"full_name", "age", "image_url", "created_at", "updated_at",
						"notification", "dark_mode",
						"language",
					}).
					AddRow(
						user.ID, user.Username, user.Email, user.Role, user.Provider,
						user.FullName, user.Age, user.ImageURL, user.CreatedAt, user.UpdatedAt,
						user.UserSetting.Notification, user.UserSetting.DarkMode,
						user.UserSetting.Language.Name,
					)

				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.ID).
					WillReturnRows(row)

				data, _ := json.Marshal(user)
				cm.ExpectSet("user_id:"+strconv.FormatUint(uint64(user.ID), 10), data, 10*time.Minute).SetVal("OK")
			},
			want: &user,
		},
		{
			name: "failed redis",
			args: args{
				profile: &user,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				cm.ExpectDel("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetErr(errors.New("redis error"))

				// publish
				rm.On(
					"Publish",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(nil)
			},
			wantErr: true,
		},
		{
			name: "failed publish",
			args: args{
				profile: &user,
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				cm.ExpectDel("user_id:" + strconv.FormatUint(uint64(user.ID), 10)).SetVal(1)

				// publish
				rm.On(
					"Publish",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(errors.New("pub error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, s.mockRedis, s.mockRmq)
			}

			got, err := s.repo.UpdateByID(tt.args.profile)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("UserRepository.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				s.T().Errorf("UserRepository.UpdateByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
