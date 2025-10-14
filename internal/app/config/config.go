package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// userSingleton содержит фиксированные ID пользователей.
type userSingleton struct {
	CreatorID   int
	ModeratorID int
}

// user — единственный экземпляр userSingleton.
var user = &userSingleton{
	CreatorID:   1,
	ModeratorID: 2,
}

// GetUserSingleton возвращает единственный экземпляр userSingleton.
func GetUserSingleton() *userSingleton {
	return user
}

// MinioConfig хранит настройки для MinIO
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

type JWTConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
}

type Config struct {
	ServiceHost string
	ServicePort int
	Minio       MinioConfig
	JWT         JWTConfig
	Redis       RedisConfig
}

type RedisConfig struct {
	Host        string
	User        string
	Password    string
	Port        int
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "config"
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	log.Info("Config successfully parsed")

	cfg.Minio = MinioConfig{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Bucket:    os.Getenv("MINIO_BUCKET"),
	}

	cfg.JWT = JWTConfig{
		SecretKey:      os.Getenv("JWT_ACCESS_SECRET"),
		AccessTokenTTL: 15 * time.Minute,
		//RefreshTokenTTL: 7 * 24 * time.Hour,
	}

	cfg.Redis = RedisConfig{
		Host:        viper.GetString("redis.host"),
		User:        viper.GetString("redis.user"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		Port:        viper.GetInt("redis.port"),
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	}

	return cfg, nil
}
