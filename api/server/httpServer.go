package server

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
	"os"
)

var DagFolder = os.Getenv("AIRFLOW__CORE__DAGS_FOLDER")
var Username = os.Getenv("_AIRFLOW_WWW_USER_USERNAME")
var Password = os.Getenv("_AIRFLOW_WWW_USER_PASSWORD")
var Address = os.Getenv("dagagent_server_address")

func init() {
	if len(DagFolder) == 0 {
		DagFolder = "./dags"
	}
	if len(Username) == 0 {
		Username = "airflow"
	}
	if len(Password) == 0 {
		Password = "airflow"
	}
	if len(Address) == 0 {
		Address = ":1323"
	}
}

func Start() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(Username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(Password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", HealthCheck)
	e.POST("/upload", Upload)

	e.Logger.Fatal(e.Start(Address))
}
