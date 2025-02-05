package repository

import (
	"crud_application/entity"
	"database/sql"
	"github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

func createConnection() *sql.DB {
	driver := os.Getenv("DATABASE_DRIVER")
	url := os.Getenv("DATABASE_URL")

	db, err := sql.Open(driver, url)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func InsertVacancy(vacancy *entity.Vacancy) int {

	db := createConnection()
	defer db.Close()

	sqlQuery := `INSERT INTO Vacancy (id, name, url, salary_from, salary_to, published_at) VALUES ($1, $2, $3, $4, $5, $6) returning id`

	var id int
	idVacancy, err := strconv.Atoi(vacancy.Id)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow(sqlQuery, idVacancy, vacancy.Name, vacancy.AlternateUrl, vacancy.Salary.From, vacancy.Salary.To, vacancy.DatePublication).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func SelectVacancy(ids []string) map[string]entity.Vacancy {

	log.Println("выбираем вакансии из базы:", ids)
	db := createConnection()
	defer db.Close()

	var idsInt []int
	for _, v := range ids {
		idInt, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		idsInt = append(idsInt, idInt)
	}

	sqlQuery := `SELECT id, name, salary_from, salary_to FROM vacancy WHERE id = ANY($1)`

	rows, err := db.Query(sqlQuery, pq.Array(idsInt))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var vacancyList = make(map[string]entity.Vacancy)
	for rows.Next() {
		var vacancy entity.Vacancy
		err := rows.Scan(&vacancy.Id, &vacancy.Name, &vacancy.Salary.From, &vacancy.Salary.To)
		if err != nil {
			log.Fatal("ошибка маппинга", err)
		}
		vacancyList[vacancy.Id] = vacancy
	}
	return vacancyList
}
