package main

import (
	"context"
	"fmt"

	"github.com/Revachol/iu5_web_5sem/internal/app/config"
	"github.com/Revachol/iu5_web_5sem/internal/app/dsn"
	"github.com/Revachol/iu5_web_5sem/internal/app/handler"
	"github.com/Revachol/iu5_web_5sem/internal/app/redis"
	"github.com/Revachol/iu5_web_5sem/internal/app/repository"

	"github.com/Revachol/iu5_web_5sem/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title История и Культура 26. Расчет бюджета создания постройки/произведения с учетом исторического индекса цен
// @version 1.0
// @description Bmstu Open IT Platform

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes http
// @BasePath /

// main инициализирует конфигурацию, репозиторий, хендлеры и запускает приложение
func main() {
	ctx := context.Background()

	router := gin.Default() // создание нового роутера Gin

	// Загрузка конфигурации приложения
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	// Получение строки подключения к PostgreSQL
	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, err := repository.NewRepository(
		postgresString,
		conf.Minio.Endpoint,
		conf.Minio.AccessKey,
		conf.Minio.SecretKey,
		conf.Minio.Bucket,
	)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	redisClient, err := redis.New(ctx, conf.Redis)
	if err != nil {
		logrus.Fatalf("failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	// Создание хендлера с подключённым репозиторием
	hand := handler.NewHandler(rep, conf, redisClient)

	// Инициализация приложения и запуск сервера
	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
