package main

import (
	"fmt"

	"github.com/Revachol/iu5_web_5sem/internal/app/config"
	"github.com/Revachol/iu5_web_5sem/internal/app/dsn"
	"github.com/Revachol/iu5_web_5sem/internal/app/handler"
	"github.com/Revachol/iu5_web_5sem/internal/app/repository"
	"github.com/Revachol/iu5_web_5sem/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// main инициализирует конфигурацию, репозиторий, хендлеры и запускает приложение
func main() {
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

	// Создание хендлера с подключённым репозиторием
	hand := handler.NewHandler(rep)

	// Инициализация приложения и запуск сервера
	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
