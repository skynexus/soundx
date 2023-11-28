package soundx

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/skynexus/soundx/api"
	"github.com/skynexus/soundx/db/sql"
	"github.com/skynexus/soundx/pkg/dbx"
)

func Start() {
	ctx := context.Background()
	log.Printf("booting soundx")
	defer log.Printf("shutdown")

	var dsn *url.URL
	if v, ok := os.LookupEnv("DSN"); !ok {
		log.Fatalf("env variable DSN not set")
	} else if u, err := url.Parse(v); err != nil {
		log.Fatalf("could not parse DSN env variable: %s", err)
	} else {
		dsn = u
	}

	db, dbErr := dbx.Open(ctx, dsn)
	if dbErr != nil {
		log.Fatalf("unable to open database connection: %s", dbErr)
	}
	defer db.Close()

	mErr := dbx.Migrate(db, sql.Archive)
	if errors.Is(mErr, dbx.ErrNoChange) {
		log.Printf("migration: already up to date")
	} else if mErr != nil {
		log.Fatalf("unable to run database migrations: %s", mErr)
	} else {
		log.Printf("migration: changes applied successfully")
	}

	if version, dirty, vErr := dbx.Version(db); vErr != nil {
		log.Fatalf("unable to determine migration version: %s", vErr)
	} else {
		log.Printf("migration: current version=%d, dirty=%t", version, dirty)
	}

	e := echo.New()

	if spec, specErr := api.GetSwagger(); specErr != nil {
		log.Printf("unable to resolve openapi specification: %s", specErr)
	} else {
		// Allow requests regardless of Hostname header.
		spec.Servers = nil

		e.Use(echomiddleware.OapiRequestValidator(spec))
	}

	s := newServer(db)
	api.RegisterHandlers(e, s)
	e.Logger.Fatal(e.Start(":8080"))
}
