package repository_test

import (
	"api-gateway-go/model"
	"api-gateway-go/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ShortenRepoSuite struct {
	suite.Suite
	mockDB      *sql.DB
	mockSQL     sqlmock.Sqlmock
	mockRedisDB *redis.Client
	mockRedis   redismock.ClientMock
	repo        repository.ShortenRepoI
}

func (s *ShortenRepoSuite) SetupTest() {
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

	repo := repository.NewShortenRepo(sqlDB, redisDB)

	s.mockDB = sqlDB
	s.mockSQL = mockSQL
	s.mockRedisDB = redisDB
	s.mockRedis = mockRedis
	s.repo = repo
}

func (s *ShortenRepoSuite) TearDownTest() {
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

func TestShortenRepoSuite(t *testing.T) {
	suite.Run(t, new(ShortenRepoSuite))
}

func (s *ShortenRepoSuite) TestShortenRepo_Get() {
	type args struct {
		hashedURL string
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, redismock.ClientMock, string)
		want       *model.APIManagement
		wantErr    bool
	}{
		{
			name: "success - redis",
			args: args{
				hashedURL: "hashedURL1",
			},
			beforeTest: func(s sqlmock.Sqlmock, r redismock.ClientMock, key string) {
				apiManagement := model.APIManagement{
					ID:                1,
					APIName:           "api1",
					ServiceName:       "service1",
					EndpointURL:       "http://localhost:5013/api1",
					HashedEndpointURL: key,
					IsAvailable:       true,
					NeedBypass:        false,
					CreatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				data, _ := json.Marshal(apiManagement)
				r.ExpectGet(key).SetVal(string(data))
			},
			want: &model.APIManagement{
				ID:                1,
				APIName:           "api1",
				ServiceName:       "service1",
				EndpointURL:       "http://localhost:5013/api1",
				HashedEndpointURL: "hashedURL1",
				IsAvailable:       true,
				NeedBypass:        false,
				CreatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
			},
		},
		{
			name: "success - db",
			args: args{
				hashedURL: "hashedURL1",
			},
			beforeTest: func(s sqlmock.Sqlmock, r redismock.ClientMock, key string) {
				apiManagement := model.APIManagement{
					ID:                1,
					APIName:           "api1",
					ServiceName:       "service1",
					EndpointURL:       "http://localhost:5013/api1",
					HashedEndpointURL: key,
					IsAvailable:       true,
					NeedBypass:        false,
					CreatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				r.ExpectGet(key).SetErr(errors.New("redis error"))

				row := s.NewRows(
					[]string{
						"id", "api_name", "service_name", "endpoint_url",
						"hashed_endpoint_url", "is_available", "need_bypass",
						"created_at", "updated_at",
					}).
					AddRow(
						apiManagement.ID, apiManagement.APIName, apiManagement.ServiceName, apiManagement.EndpointURL,
						apiManagement.HashedEndpointURL, apiManagement.IsAvailable, apiManagement.NeedBypass,
						apiManagement.CreatedAt, apiManagement.UpdatedAt,
					)

				s.ExpectPrepare("SELECT .* FROM api_managements .*").
					ExpectQuery().
					WithArgs(key).
					WillReturnRows(row)

				data, _ := json.Marshal(apiManagement)
				r.ExpectSet(key, data, 10*time.Minute).SetVal("OK")
			},
			want: &model.APIManagement{
				ID:                1,
				APIName:           "api1",
				ServiceName:       "service1",
				EndpointURL:       "http://localhost:5013/api1",
				HashedEndpointURL: "hashedURL1",
				IsAvailable:       true,
				NeedBypass:        false,
				CreatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
			},
		},
		{
			name: "failed - db - stmt",
			args: args{
				hashedURL: "hashedURL1",
			},
			beforeTest: func(s sqlmock.Sqlmock, r redismock.ClientMock, key string) {
				r.ExpectGet(key).SetErr(errors.New("redis error"))

				s.ExpectPrepare("SELECT .* FROM api_managements .*").
					WillReturnError(errors.New("prepare stmt error"))
			},
			wantErr: true,
		},
		{
			name: "failed - db - scan",
			args: args{
				hashedURL: "hashedURL1",
			},
			beforeTest: func(s sqlmock.Sqlmock, r redismock.ClientMock, key string) {
				apiManagement := model.APIManagement{
					ID:                1,
					APIName:           "api1",
					ServiceName:       "service1",
					EndpointURL:       "http://localhost:5013/api1",
					HashedEndpointURL: key,
					IsAvailable:       true,
					NeedBypass:        false,
					UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				r.ExpectGet(key).SetErr(errors.New("redis error"))

				row := s.NewRows(
					[]string{
						"id", "api_name", "service_name", "endpoint_url",
						"hashed_endpoint_url", "is_available", "need_bypass",
						"created_at", "updated_at",
					}).
					AddRow(
						apiManagement.ID, apiManagement.APIName, apiManagement.ServiceName, apiManagement.EndpointURL,
						apiManagement.HashedEndpointURL, apiManagement.IsAvailable, apiManagement.NeedBypass,
						nil, apiManagement.UpdatedAt,
					).
					RowError(1, errors.New("scan error"))

				s.ExpectPrepare("SELECT .* FROM api_managements .*").
					ExpectQuery().
					WithArgs(key).
					WillReturnRows(row)
			},
			wantErr: true,
		},
		{
			name: "failed - redis - set",
			args: args{
				hashedURL: "hashedURL1",
			},
			beforeTest: func(s sqlmock.Sqlmock, r redismock.ClientMock, key string) {
				apiManagement := model.APIManagement{
					ID:                1,
					APIName:           "api1",
					ServiceName:       "service1",
					EndpointURL:       "http://localhost:5013/api1",
					HashedEndpointURL: key,
					IsAvailable:       true,
					NeedBypass:        false,
					CreatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
					UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				}

				r.ExpectGet(key).SetErr(errors.New("redis error"))

				row := s.NewRows(
					[]string{
						"id", "api_name", "service_name", "endpoint_url",
						"hashed_endpoint_url", "is_available", "need_bypass",
						"created_at", "updated_at",
					}).
					AddRow(
						apiManagement.ID, apiManagement.APIName, apiManagement.ServiceName, apiManagement.EndpointURL,
						apiManagement.HashedEndpointURL, apiManagement.IsAvailable, apiManagement.NeedBypass,
						apiManagement.CreatedAt, apiManagement.UpdatedAt,
					)

				s.ExpectPrepare("SELECT .* FROM api_managements .*").
					ExpectQuery().
					WithArgs(key).
					WillReturnRows(row)

				data, _ := json.Marshal(apiManagement)
				r.ExpectSet(key, data, 10*time.Minute).SetErr(errors.New("redis error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, s.mockRedis, tt.args.hashedURL)
			}
			got, err := s.repo.Get(tt.args.hashedURL)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("ShortenRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				s.T().Errorf("ShortenRepo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *ShortenRepoSuite) TestShortenRepo_Create() {
	type args struct {
		apiManagement *model.APIManagement
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, string)
		want       *model.APIManagement
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				apiManagement: &model.APIManagement{
					APIName:           "wislists",
					ServiceName:       "wislists",
					EndpointURL:       "http://localhost:5013/wislists",
					HashedEndpointURL: "CUbICN8VgNk",
					IsAvailable:       true,
					NeedBypass:        false,
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(uint(1), time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local), time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local))

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs("wislists", "wislists", "http://localhost:5013/wislists", "CUbICN8VgNk", true, false).
					WillReturnRows(rows)
			},
			want: &model.APIManagement{
				ID:                1,
				APIName:           "wislists",
				ServiceName:       "wislists",
				EndpointURL:       "http://localhost:5013/wislists",
				HashedEndpointURL: "CUbICN8VgNk",
				IsAvailable:       true,
				NeedBypass:        false,
				CreatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
			},
		},
		{
			name: "failed to prepare stmt",
			args: args{
				apiManagement: &model.APIManagement{
					APIName:           "wislists",
					ServiceName:       "wislists",
					EndpointURL:       "http://localhost:5013/wislists",
					HashedEndpointURL: "CUbICN8VgNk",
					IsAvailable:       true,
					NeedBypass:        false,
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectPrepare(regexp.QuoteMeta(query)).
					WillReturnError(errors.New("prepare stmt error"))
			},
			wantErr: true,
		},
		{
			name: "failed to scan",
			args: args{
				apiManagement: &model.APIManagement{
					APIName:           "wislists",
					ServiceName:       "wislists",
					EndpointURL:       "http://localhost:5013/wislists",
					HashedEndpointURL: "CUbICN8VgNk",
					IsAvailable:       true,
					NeedBypass:        false,
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(uint(1), nil, time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local)).
					RowError(1, errors.New("scan error"))

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs("wislists", "wislists", "http://localhost:5013/wislists", "CUbICN8VgNk", true, false).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := `
			INSERT INTO api_managements (api_name, service_name, endpoint_url, hashed_endpoint_url, is_available, need_bypass)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at
			`
			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}
			got, err := s.repo.Create(tt.args.apiManagement)
			if (err != nil) != tt.wantErr {
				s.T().Errorf("ShortenRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				s.T().Errorf("ShortenRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
