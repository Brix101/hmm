package main

import (
	"log"
	"net/http"

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://192.168.254.180:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Static("/assets", "web/dist/assets")
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
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "web/dist",
		Index:  "Index.html",
		Browse: false,
		HTML5:  true,
	}))

	e.File("/*", "assets/index.html")

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
