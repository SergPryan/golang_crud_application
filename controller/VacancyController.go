package controller

import (
	"crud_application/entity"
	"crud_application/repository"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func HandlerGetVacancy(writer http.ResponseWriter, request *http.Request) {
	idPath := request.URL.Path
	idPath = strings.TrimPrefix(idPath, "/vacancy/")
	var ids []string
	ids = append(ids, idPath)
	res := repository.SelectVacancy(ids)
	resp, ok := res[idPath]
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		js, err := json.Marshal(resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		writer.Write(js)
	}
}

type VacancyPage struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func HandlerPostRequestNewVacancies(w http.ResponseWriter, r *http.Request) {
	//TODO замена на json

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var vacancyPage VacancyPage
	err = json.Unmarshal(body, &vacancyPage)

	log.Println(vacancyPage)

	var countPage = vacancyPage.End - vacancyPage.Start

	if countPage < 0 {
		log.Fatal("Vacancy page is out of range")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chanPageVacancy := make(chan entity.PageVacancy, countPage)

	var wg sync.WaitGroup
	for pageNumber := vacancyPage.Start; pageNumber < vacancyPage.End; pageNumber++ {
		log.Println("Получаем вакансии с сайта, page:" + strconv.Itoa(pageNumber))
		wg.Add(1)
		go getVacancyFromHH(chanPageVacancy, pageNumber, &wg)
	}
	log.Println("ждем данных из канала")
	wg.Wait()

	close(chanPageVacancy)
	log.Println("Закрытие канала")

	var ids []string
	var vacancies []entity.Vacancy
	for value := range chanPageVacancy {
		fmt.Println(value)
		for _, vacancy := range value.Vacancy {
			ids = append(ids, vacancy.Id)
			vacancies = append(vacancies, vacancy)
		}
	}

	storedVacancy := repository.SelectVacancy(ids)
	//fmt.Println("вакансии которые уже хранятся в базе:", storedVacancy)

	var idsResult []int
	for _, v := range vacancies {
		_, ok := storedVacancy[v.Id]
		if !ok {
			log.Println("вакансия которая не храниться в базе:", v)
			id := repository.InsertVacancy(&v)
			idsResult = append(idsResult, id)
		} else {
			log.Println("вакасия уже храниться в базе:", v.Id)
		}
	}
	strRes := fmt.Sprintf("Обработанные id:%d", idsResult)
	log.Println(strRes)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strRes))
}

func getVacancyFromHH(chanVacancies chan entity.PageVacancy, page int, wg *sync.WaitGroup) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	url := os.Getenv("URL_HH_RU")

	go func() {
		defer wg.Done()
		response, err := client.Get(url + strconv.Itoa(page))
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		log.Printf("ответ от hh.ru:%d page:%d", response.StatusCode, page)
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var data entity.PageVacancy
		errJson := json.Unmarshal(body, &data)
		if errJson != nil {
			log.Fatal(errJson)
		}

		fmt.Println("отправляем данные в канал, page: " + strconv.Itoa(page))
		chanVacancies <- data
	}()
}
