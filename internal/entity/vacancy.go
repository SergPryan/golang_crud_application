package entity

type Vacancy struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	DatePublication string `json:"published_at,omitempty"`
	Salary          Salary `json:"salary,omitempty"`
	AlternateUrl    string `json:"alternate_url,omitempty"`
}

type Salary struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type PageVacancy struct {
	Vacancy []Vacancy `json:"items"`
}

type Area struct {
	Name string `json:"name"`
}
