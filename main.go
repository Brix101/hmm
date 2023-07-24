package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Brix101/network-file-manager/internal/api/files"
	"github.com/Brix101/network-file-manager/internal/api/users"
	"github.com/Brix101/network-file-manager/internal/middlewares"
	"github.com/Brix101/network-file-manager/internal/utils/db"
	"github.com/Brix101/network-file-manager/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
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
	e.Renderer = templates.New()
	e.Validator = middlewares.NewValidator()

	d := db.Init()
	us := users.NewUserServices(d)
	reader := files.NewReader("", false)

	sr := users.UserHandler{UserServices: us}
	fr := files.FileHandler{Reader: reader}

	v1 := e.Group("/api")
	sr.Routes(v1)
	fr.Routes(v1)

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
	s := http.Server{
		Addr:    "0.0.0.0:5000",
		Handler: e,
		// ReadTimeout: 30 * time.Second, // customize http.Server timeouts
	}
	log.Println("🚀🚀🚀 Server at http://" + s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Handler
func hello(c echo.Context) error {
	return c.Render(http.StatusOK, "sign-in.html", "asdfasdf")
}
