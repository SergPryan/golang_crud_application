package main

import (
	"crud_application/entity"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {

	var db *sql.DB
	// setup database

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	response, err := client.Get("https://api.hh.ru/vacancies?area=1202&specialization=1")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data entity.PageVacancy

	errJson := json.Unmarshal(body, &data)
	if errJson != nil {
		log.Fatal(errJson)
	}
	for _, v := range data.Vacancy {
		fmt.Println(v.Area, v.Name, v.DatePublication, v.Salary)
	}
}
