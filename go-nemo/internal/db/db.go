package db

import (
	"fmt"
	"log"
	"netease-kit/nemo/internal/config"

	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func InitDB() {
	cfg := config.AppConfig.Spring.Datasource
	// Construct DSN. Ideally parse from cfg.Url (JDBC format) but for now hardcode structure with config creds
	// JDBC: jdbc:mysql://127.0.0.1:3306/nemo?useUnicode=true&characterEncoding=utf8&serverTimezone=Asia/Shanghai
	// Go DSN: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/nemo?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	if cfg.Hikari.MaximumPoolSize > 0 {
		sqlDB.SetMaxOpenConns(cfg.Hikari.MaximumPoolSize)
	}

	// Init Redis
	redisCfg := config.AppConfig.Spring.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.Database,
	})

	_, err = RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("failed to connect redis: %v", err)
		// Don't fatal here as Redis might be optional for some parts or local dev without redis
	}
}
