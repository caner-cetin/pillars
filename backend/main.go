package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"pillars-backend/src/api"
	"pillars-backend/src/constants"
	"pillars-backend/src/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/", os.Getenv("MONGO_INITDB_ROOT_USERNAME"), os.Getenv("MONGO_INITDB_ROOT_PASSWORD"), "mongodb")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	constants.MONGODB = client
	db.CreateIndexes()
	var e *echo.Echo
	e = echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogLevel: log.ERROR,
	}))
	// VERY BAD PRACTICE, and I know it.
	// This is just for the sake of the assignment.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// error might be nil but status code may be 500
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	rest := e.Group("/api")
	rest.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})
	rest.GET("/earthquakes/count", api.GetEarthquakeCount)
	rest.GET("/earthquakes/get", api.GetEarthquakesPaged)
	rest.PUT("/earthquakes/insert", api.InsertSingleEarthquake)
	rest.PUT("/earthquakes/bulk-insert", api.InsertBulkEarthquakes)
	rest.DELETE("/earthquakes/delete", api.WipeAllEarthquakes)
	//
	//
	//
	//
	ws := e.Group("/ws")
	ws.GET("/earthquakes/stream", api.GetEarthquakeStream)
	e.Logger.Fatal(e.Start(":1323"))
}
