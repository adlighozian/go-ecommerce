package repository_test

import (
	"auth-go/mocks"
	"auth-go/model"
	"auth-go/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AuthRepoSuite struct {
	suite.Suite
	mockDB      *sql.DB
	mockSQL     sqlmock.Sqlmock
	mockRedisDB *redis.Client
	mockRedis   redismock.ClientMock
	mockRmq     *mocks.RabbitMQClient
	repo        repository.AuthRepositoryI
}

func (s *AuthRepoSuite) SetupTest() {
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

	repo := repository.NewAuthRepository(sqlDB, redisDB, rmq)

	s.mockDB = sqlDB
	s.mockSQL = mockSQL
	s.mockRedisDB = redisDB
	s.mockRedis = mockRedis
	s.mockRmq = rmq
	s.repo = repo
}

func (s *AuthRepoSuite) TearDownTest() {
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

func TestAuthRepoSuite(t *testing.T) {
	suite.Run(t, new(AuthRepoSuite))
}

func (s *AuthRepoSuite) TestAuthRepository_Create() {
	type args struct {
		user *model.User
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
				user: &model.User{
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				user := &model.User{
					ID:        1,
					Username:  "test",
					Email:     "test@test.com",
					Password:  "hashed-password",
					Role:      "user",
					Provider:  "email",
					FullName:  "testtist",
					Age:       21, // i swear
					ImageURL:  "some-img",
					CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				// isEmailUnique
				row1 := s.NewRows([]string{"count"}).AddRow(0)
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row1)

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

				// get user data back
				row2 := s.NewRows(
					[]string{
						"id", "username", "email", "password", "role", "provider",
						"full_name", "age", "image_url", "created_at", "updated_at",
					},
				).
					AddRow(
						user.ID, user.Username, user.Email, user.Password, user.Role, user.Provider,
						user.FullName, user.Age, user.ImageURL, user.CreatedAt, user.UpdatedAt,
					)
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row2)
			},
			want: &model.User{
				ID:        1,
				Username:  "test",
				Email:     "test@test.com",
				Password:  "hashed-password",
				Role:      "user",
				Provider:  "email",
				FullName:  "testtist",
				Age:       21, // i swear
				ImageURL:  "some-img",
				CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
			},
		},
		{
			name: "failed - email - not unique",
			args: args{
				user: &model.User{
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				user := &model.User{
					// ID:        1,
					// Username:  "test",
					Email: "test@test.com",
					// Password:  "hashed-password",
					// Role:      "user",
					// Provider:  "email",
					// FullName:  "testtist",
					// Age:       21, // i swear
					// ImageURL:  "some-img",
					// CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					// UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				// isEmailUnique
				row1 := s.NewRows([]string{"count"}).AddRow(1)
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row1)
			},
			wantErr: true,
		},
		{
			name: "failed - email - stmt err",
			args: args{
				user: &model.User{
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				// isEmailUnique
				s.ExpectPrepare("SELECT .* FROM users .*").
					WillReturnError(errors.New("prepare stmt error"))
			},
			wantErr: true,
		},
		{
			name: "failed - email - scan err",
			args: args{
				user: &model.User{
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				user := &model.User{
					// ID:        1,
					// Username:  "test",
					Email: "test@test.com",
					// Password:  "hashed-password",
					// Role:      "user",
					// Provider:  "email",
					// FullName:  "testtist",
					// Age:       21, // i swear
					// ImageURL:  "some-img",
					// CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					// UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				// isEmailUnique
				row1 := s.NewRows([]string{"count"}).AddRow(nil).RowError(1, errors.New("scan error"))
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row1)
			},
			wantErr: true,
		},
		{
			name: "failed - pub err",
			args: args{
				user: &model.User{
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				user := &model.User{
					// ID:        1,
					// Username:  "test",
					Email: "test@test.com",
					// Password:  "hashed-password",
					// Role:      "user",
					// Provider:  "email",
					// FullName:  "testtist",
					// Age:       21, // i swear
					// ImageURL:  "some-img",
					// CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					// UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				// isEmailUnique
				row1 := s.NewRows([]string{"count"}).AddRow(0)
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row1)

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
		{
			name: "failed - get data back - stmt err",
			args: args{
				user: &model.User{
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				user := &model.User{
					// ID:        1,
					// Username:  "test",
					Email: "test@test.com",
					// Password:  "hashed-password",
					// Role:      "user",
					// Provider:  "email",
					// FullName:  "testtist",
					// Age:       21, // i swear
					// ImageURL:  "some-img",
					// CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					// UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				// isEmailUnique
				row1 := s.NewRows([]string{"count"}).AddRow(0)
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row1)

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

				// get user data back
				s.ExpectPrepare("SELECT .* FROM users .*").
					WillReturnError(errors.New("prepare stmt error"))
			},
			wantErr: true,
		},
		{
			name: "failed - get data back - scan err",
			args: args{
				user: &model.User{
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, cm redismock.ClientMock, rm *mocks.RabbitMQClient) {
				user := &model.User{
					ID:       1,
					Username: "test",
					Email:    "test@test.com",
					Password: "hashed-password",
					Role:     "user",
					Provider: "email",
					// FullName:  "testtist",
					Age:       21, // i swear
					ImageURL:  "some-img",
					CreatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					UpdatedAt: time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				// isEmailUnique
				row1 := s.NewRows([]string{"count"}).AddRow(0)
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row1)

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

				// get user data back
				row2 := s.NewRows(
					[]string{
						"id", "username", "email", "password", "role", "provider",
						"full_name", "age", "image_url", "created_at", "updated_at",
					},
				).
					AddRow(
						user.ID, user.Username, user.Email, user.Password, user.Role, user.Provider,
						nil, user.Age, user.ImageURL, user.CreatedAt, user.UpdatedAt,
					).
					RowError(1, errors.New("scan error"))
				s.ExpectPrepare("SELECT .* FROM users .*").
					ExpectQuery().
					WithArgs(user.Email).
					WillReturnRows(row2)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, s.mockRedis, s.mockRmq)
			}
			got, err := s.repo.Create(tt.args.user)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("AuthRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				s.T().Errorf("AuthRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *AuthRepoSuite) TestAuthRepository_SetRefreshToken() {
	type args struct {
		token           string
		dataByte        []byte
		refreshTokenDur time.Duration
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(redismock.ClientMock)
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				token:           "this-is-refresh-token-key",
				dataByte:        []byte("yes"),
				refreshTokenDur: 10 * time.Minute,
			},
			beforeTest: func(cm redismock.ClientMock) {
				token := "this-is-refresh-token-key"
				dataByte := []byte("yes")
				refreshTokenDur := 10 * time.Minute

				cm.ExpectSet("refresh_token:"+token, dataByte, refreshTokenDur).SetVal("OK")
			},
		},
		{
			name: "error",
			args: args{
				token:           "this-is-refresh-token-key",
				dataByte:        []byte("yes"),
				refreshTokenDur: 10 * time.Minute,
			},
			beforeTest: func(cm redismock.ClientMock) {
				token := "this-is-refresh-token-key"
				dataByte := []byte("yes")
				refreshTokenDur := 10 * time.Minute

				cm.ExpectSet("refresh_token:"+token, dataByte, refreshTokenDur).SetErr(errors.New("redis error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.beforeTest != nil {
				tt.beforeTest(s.mockRedis)
			}
			err := s.repo.SetRefreshToken(tt.args.token, tt.args.dataByte, tt.args.refreshTokenDur)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("AuthRepository.SetRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *AuthRepoSuite) TestAuthRepository_GetByRefreshToken() {
	type args struct {
		token string
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(redismock.ClientMock)
		want       *model.RefreshToken
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				token: "this-is-refresh-token-key",
			},
			beforeTest: func(cm redismock.ClientMock) {
				token := "this-is-refresh-token-key"
				content := model.RefreshToken{
					UserID: 1, UserRole: "user",
				}

				data, _ := json.Marshal(content)
				cm.ExpectGet("refresh_token:" + token).SetVal(string(data))
			},
			want: &model.RefreshToken{
				UserID: 1, UserRole: "user",
			},
		},
		{
			name: "failed",
			args: args{
				token: "this-is-refresh-token-key",
			},
			beforeTest: func(cm redismock.ClientMock) {
				token := "this-is-refresh-token-key"

				cm.ExpectGet("refresh_token:" + token).SetErr(errors.New("redis error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.beforeTest != nil {
				tt.beforeTest(s.mockRedis)
			}
			got, err := s.repo.GetByRefreshToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("AuthRepository.GetByRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				s.T().Errorf("AuthRepository.GetByRefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
