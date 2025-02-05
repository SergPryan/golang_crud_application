-- +goose Up
-- +goose StatementBegin
create table if not exists vacancy (
                                       id BIGINT PRIMARY KEY,
                                       name VARCHAR(1000) NOT NULL,
                                       url varchar(5000),
                                       salary_from int ,
                                       salary_to int,
                                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                       published_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists vacancy
-- +goose StatementEnd
