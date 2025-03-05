# Приложение для загрузки вакансий из сайта hh.ru

забирает вакасии по апи api.hh.ru, где start это началальная страница, а end это конец
curl --location 'http://localhost:8080/request' \
--header 'Content-Type: application/json' \
--data '{
"start": 14,
"end": 15
}'

вывод сохраненной вакансии из базы данных
curl --location 'http://localhost:8080/vacancy/116766981'

в проекте есть докер, можно запустить docker compose up