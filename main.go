package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Brix101/network-file-manager/internal/api/users"
	"github.com/Brix101/network-file-manager/internal/middlewares"
	"github.com/Brix101/network-file-manager/internal/utils/db"
	"github.com/Brix101/network-file-manager/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	t := templates.New()
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	// Be careful to use constant time comparison to prevent timing attacks
	// 	if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
	// 		subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))
	e.Static("/static", "templates/static")
	e.Renderer = t
	e.Validator = middlewares.NewValidator()

	d := db.Init()
	us := users.NewUserServices(d)

	sr := users.UserHandler{UserServices: us}

	v1 := e.Group("/api")
	sr.Routes(v1)
	// Routes
	e.GET("/", hello)
	e.POST("/sign-in", func(c echo.Context) error {
		email := c.FormValue("email")

		password := c.FormValue("password")

		fmt.Println("++++++++++++++++++++", email, password)

		time.Sleep(30 * time.Second)
		return c.Render(http.StatusOK, "sign-in.html", "asdfasdf")
	})
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.Render(http.StatusOK, "sign-in.html", "asdfasdf")
}
