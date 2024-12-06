package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kemov/todo-app"
	"github.com/kemov/todo-app/pkg/handler"
	"github.com/kemov/todo-app/pkg/repository"
	"github.com/kemov/todo-app/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка инициализации configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("ошибка загрузки переменных окружения: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("Не удалось инициализировать базу данных: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("произошла ошибка при запуске http-сервера: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

/*
	if err := viper.ReadInConfig(); err != nil {
		logrus.Printf("ошибка чтения кофигурации: %s", err)
		return err
	}

	return nil
*/

