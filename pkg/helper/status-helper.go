package helper

import (
	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/model"
)

func GetStatusByName(name string) (*model.StatusModel, error) {
	selectedStatus, err := db.GetSingleDataByQuery[model.StatusModel]("select * from public.status where lower(name) = lower($1) limit 1", name)
	if err != nil {
		return nil, err
	}

	return selectedStatus, nil
}

func GetStatusByCode(code, category string) (*model.StatusModel, error) {
	selectedStatus, err := db.GetSingleDataByQuery[model.StatusModel]("select * from public.status where (lower(code) = lower($1) or lower(name) = lower($1)) and lower(category) = lower($2)", code, category)
	if err != nil {
		return nil, err
	}

	return selectedStatus, nil
}
