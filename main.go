package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"waroong-be/apps/routers"
	"waroong-be/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader: "USER-SERVICE",
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
	})

	// to enable Cross-Origin Resource Sharing with various options.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PATCH,DELETE",
	}))

	// To recover from a panic thrown by any handler in the stack
	app.Use(recover.New())
	// for Fiber to let's caches be more efficient and save bandwidth,
	// as a web server does not need to resend a full response if the content has not changed.
	app.Use(etag.New())

	app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format: "[${time}][${ip}]:${port} ${status} - ${method} ${path} -- query: ${queryParams} -- body: ${body}\n",
	}))
	// no need for now
	// metricPage := app.Group("/v1")
	// metricPage.Get("/metrics", monitor.New(monitor.Config{Title: "Waroong - Backend", Refresh: 5 * time.Second}))

	// set up DB connection here
	DBConnection := config.InitDatabase()

	// setup initial routes
	apiV1 := app.Group("/api/v1")
	// set dispatch to the main router
	routers.Dispatch(DBConnection, apiV1)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	var serverShutdown sync.WaitGroup

	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(60 * time.Second)
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = config.GoDotEnvVariable("APP_PORT")
	}
	fmt.Println("⚡️ [" + config.GoDotEnvVariable("APPLICATION_ENV") + "] - " + config.GoDotEnvVariable("APP_NAME") + " IS RUNNING ON PORT - " + port)
	if err := app.Listen(":" + port); err != nil {
		log.Panic(err)
		panic(err)
	}

	serverShutdown.Wait()

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	fmt.Println("Fiber was successful shutdown.")
}
