package global

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"main/app/internal/model/config"
)

var (
	Config  config.Config
	Logger  *zap.Logger
	MysqlDB *sql.DB
	Rdb     *redis.Client
)
