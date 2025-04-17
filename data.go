package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SarkiMudboy/shuttle/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("mysql", "shuttle:2580@/shuttle?parseTime=true")
	if err != nil {
		fmt.Printf("An error occured initializing the database: %s", err.Error())
		os.Exit(1)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// migrate any changes
	runMigrations()
}

func runMigrations() {

	m, err := migrate.New("file://database/migrations", "mysql://shuttle:2580@/shuttle?")

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("An error occured: %s", err.Error())
		os.Exit(1)
	}

	log.Println("Migration Sucessful")
}

func getRequestHistoryMethod(method string) database.RequestHistoryMethod {

	switch method {
	case "GET":
		return database.RequestHistoryMethodGET
	case "POST":
		return database.RequestHistoryMethodPOST
	case "PATCH":
		return database.RequestHistoryMethodPATCH
	case "PUT":
		return database.RequestHistoryMethodPUT
	case "DELETE":
		return database.RequestHistoryMethodDELETE
	case "TRACE":
		return database.RequestHistoryMethodTRACE
	}

	return ""
}

func SaveRequestToHistory(r *request) error {

	ctx := context.Background()
	queries := database.New(db)

	//add to database
	_, err := queries.CreateRequest(ctx, database.CreateRequestParams{
		Endpoint: r.location,
		Headers:  sql.NullString{String: r.headers.rawHeaders},
		Method: database.NullRequestHistoryMethod{
			RequestHistoryMethod: getRequestHistoryMethod(r.method),
			Valid:                true,
		},
		Body: sql.NullString{String: r.body},
	})

	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("request to %s saved\n", r.location)
	return nil
}

func getRequestHistory() error {

	ctx := context.Background()

	queries := database.New(db)
	requests, err := queries.GetRequestHistory(ctx)

	if err != nil {
		return err
	}

	for r := range requests {
		log.Printf("%v\n", r)
	}

	return nil
}

func parseRequestFromDatabase(r *request, request database.RequestHistory) {

	r.location = request.Endpoint

	if request.Method.Valid {
		r.method = string(request.Method.RequestHistoryMethod)
	}

	if request.Headers.Valid {
		r.headers.rawHeaders = request.Headers.String
	}

	if request.Body.Valid {
		r.body = request.Body.String
	}

}

func getLastRequest(r *request) error {

	ctx := context.Background()
	queries := database.New(db)

	request, err := queries.GetlastRequest(ctx)

	if err != nil {
		return err
	}

	parseRequestFromDatabase(r, request)

	return nil
}
