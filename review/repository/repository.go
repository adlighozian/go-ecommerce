package  repository

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"auth-go/db"
)

type  repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repositorier {
	return &repository{
		db: db,
	}
}

func (repo *repository) Get(ctx *gin.Context) {
}

func (repo *repository) GetDetail(ctx *gin.Context) {
}

func (repo *repository) Create(ctx *gin.Context) {
}

func (repo *repository) Update(ctx *gin.Context) {
}

func (repo *repository) Delete(ctx *gin.Context) {
}