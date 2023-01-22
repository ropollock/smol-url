package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"
	"net/http"
	"os"
	"server/config"
	"server/controller"
	"server/dao"
	"server/data"
	"server/model"
	"server/service"
)

func main() {
	fmt.Println("Loading config.")

	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables config.", err)
	}

	config.AppConfig = &conf

	databaseProvider := data.PostgresDBProvider()
	databaseProvider.Connect(conf.DBUri)
	databaseProvider.GetDB().AutoMigrate(&model.User{})

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
	}))

	logger := lecho.New(
		os.Stdout,
		lecho.WithLevel(log.DEBUG),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)
	e.Logger = logger

	e.Use(middleware.RequestID())
	e.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.GET("/api/healthcheck", healthcheck)

	userDao := dao.UserDao(databaseProvider)
	userService := service.UserService(userDao)
	authService := service.AuthService(userService)

	usersController := controller.UsersController(userService, authService)
	usersController.RegisterUserRoutes(e)

	authController := controller.AuthController(userService, authService)
	authController.RegisterLoginRoutes(e)

	_, err = userService.FindUserByUsername("admin")
	if err != nil {
		adminPass, _ := userService.HashPassword("admin")
		adminUser := model.User{IsAdmin: true, Username: "admin", Password: adminPass}
		userService.CreateUser(&adminUser)
	}

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &model.Claims{},
		SigningKey:              []byte(authService.GetJWTSecret()),
		TokenLookup:             "cookie:access-token,header:Authorization",
		ErrorHandlerWithContext: authController.JWTErrorChecker,
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/login" {
				return true
			}
			return false
		},
	}))

	e.Use(authController.TokenRefresherMiddleware)
	e.Logger.Fatal(e.Start(":" + conf.ServerPort))
}

func healthcheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}
