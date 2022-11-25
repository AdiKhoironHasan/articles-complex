package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AdiKhoironHasan/articles-complex/pkg/database"

	integ "github.com/AdiKhoironHasan/articles-complex/internal/integration"
	SqlRepo "github.com/AdiKhoironHasan/articles-complex/internal/repository/postgresql"
	NoSqlRepo "github.com/AdiKhoironHasan/articles-complex/internal/repository/redis"
	"github.com/AdiKhoironHasan/articles-complex/internal/services"
	handlers "github.com/AdiKhoironHasan/articles-complex/internal/transport/http"
	"github.com/AdiKhoironHasan/articles-complex/internal/transport/http/middleware"

	"github.com/apex/log"
	"github.com/labstack/echo"

	"github.com/spf13/viper"
)

func main() {

	// create error channerl
	errChan := make(chan error)

	// create echo and middleware
	e := echo.New()
	m := middleware.NewMidleware()

	// set config env
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config-dev")

	// viper to try read config file
	err := viper.ReadInConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// fill variable postgresql from config
	dbhost, dbUser, dbPassword, dbName, dbPort :=
		viper.GetString("db.postgre.host"),
		viper.GetString("db.postgre.user"),
		viper.GetString("db.postgre.password"),
		viper.GetString("db.postgre.dbname"),
		viper.GetString("db.postgre.port")

	// try connect db postgresql
	PostgreSqlDB, err := database.PostgreSqllInitialize(dbhost, dbUser, dbPassword, dbName, dbPort)

	if err != nil {
		log.Fatal("Failed to Connect Postgre SQL Database: " + err.Error())
	}

	// fill variable redis from config
	dbhost, dbUser, dbPassword, dbPort =
		viper.GetString("db.redis.host"),
		viper.GetString("db.redis.user"),
		viper.GetString("db.redis.password"),
		viper.GetString("db.redis.port")

	// try connect db redis
	RedisDB, err := database.RedislInitialize(dbhost, dbUser, dbPassword, dbPort)
	if err != nil {
		log.Fatal("Failed to Connect Redis Database: " + err.Error())
	}

	// execute function in end line of program
	defer func() {
		// disconnect db postgresql
		err := PostgreSqlDB.Conn.Close()
		if err != nil {
			log.Fatal(err.Error())
		}

		// disconnect db redis
		err = RedisDB.Conn.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	// echo use cors middleware
	e.Use(m.CORS)

	// create repo, handler and service
	sqlrepo := SqlRepo.NewRepo(PostgreSqlDB.Conn)
	noSqlRepo := NoSqlRepo.NewRepo(RedisDB.Conn)
	integSrv := integ.NewService()
	srv := services.NewService(sqlrepo, noSqlRepo, integSrv)
	handlers.NewHttpHandler(e, srv)

	// run func with concurrency
	go func() {
		// send signall error to channel
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		// fill channel from start echo
		errChan <- e.Start(":" + viper.GetString("server.port"))
	}()

	// service succesfully run
	e.Logger.Print("Starting ", viper.GetString("appName"))

	// fill error from channel
	err = <-errChan
	log.Error(err.Error())

}
