package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

// Get All check-list
func (m *postgresDBRepo) GetAllCheckList() (*[]models.CheckList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var checklists []models.CheckList

	queryCheckLists := `
	select id_check_list, name_check_list, description, id_type_of_list
	from check_lists
	`

	queryPointsCheckList := `
	select name_point 
	from check_list_points
	where id_check_list = $1
	`

	queryNameOfType := `
	select name_type 
	from check_list_type
	where id_type_of_list = $1
	`
	rows, err := m.DB.QueryContext(ctx, queryCheckLists)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var checklist models.CheckList
		var idType int
		err := rows.Scan(&checklist.ID, &checklist.Name, &checklist.Description, &idType)
		if err != nil {
			return nil, err
		}

		rowPoints, err := m.DB.QueryContext(ctx, queryPointsCheckList, checklist.ID)
		if err != nil {
			return nil, err
		}
		defer rowPoints.Close()

		var namePoints []string
		for rowPoints.Next() {
			var point string

			err := rowPoints.Scan(&point)
			if err != nil {
				return nil, err
			}

			namePoints = append(namePoints, point)
		}

		checklist.NameOfPoints = namePoints

		rowName := m.DB.QueryRowContext(ctx, queryNameOfType, idType)
		err = rowName.Scan(&checklist.TypeOfList)
		if err != nil {
			return nil, err
		}

		checklists = append(checklists, checklist)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &checklists, nil
}

func (m *postgresDBRepo) GetCheckListByID(id int) (*models.CheckList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var checklist models.CheckList

	queryCheckList := `
	select id_check_list, name_check_list, description, id_type_of_list
	from check_lists
	where id_check_list = $1
	`

	queryPointsCheckList := `
	select name_point 
	from check_list_points
	where id_check_list = $1
	`

	queryNameOfType := `
	select name_type 
	from check_list_type
	where id_type_of_list = $1
	`
	row := m.DB.QueryRowContext(ctx, queryCheckList, id)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var idType int
	err := row.Scan(&checklist.ID, &checklist.Name, &checklist.Description, &idType)
	if err != nil {
		return nil, err
	}

	rowPoints, err := m.DB.QueryContext(ctx, queryPointsCheckList, checklist.ID)
	if err != nil {
		return nil, err
	}
	defer rowPoints.Close()
	var namePoints []string
	for rowPoints.Next() {
		var point string

		err := rowPoints.Scan(&point)
		if err != nil {
			return nil, err
		}

		namePoints = append(namePoints, point)
	}

	checklist.NameOfPoints = namePoints

	rowName := m.DB.QueryRowContext(ctx, queryNameOfType, idType)
	err = rowName.Scan(&checklist.TypeOfList)
	if err != nil {
		return nil, err
	}

	return &checklist, nil
}
