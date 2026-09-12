package db

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetSingleDataByQuery[T any](query string, param ...interface{}) (*T, error) {
	rows, err := Conn.Query(context.Background(), query, param...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &result, nil
}

func GetMultipleDataByQuery[T any](query string, param ...interface{}) (*[]T, error) {
	rows, err := Conn.Query(context.Background(), query, param...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &result, nil
}

func ExecuteQuery(query string, param ...interface{}) error {
	_, err := Conn.Exec(context.Background(), query, param...)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func InsertReturnUUID(query string, param ...interface{}) (*uuid.UUID, error) {
	var id uuid.UUID
	err := Conn.QueryRow(context.Background(), query, param...).Scan(&id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &id, nil
}
