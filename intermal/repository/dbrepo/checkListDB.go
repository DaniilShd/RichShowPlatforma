package dbrepo

import (
	"context"
	"errors"
	"fmt"
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
	rows, err := m.DB.QueryContext(ctx, queryCheckLists)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var checklist models.CheckList
		err := rows.Scan(&checklist.ID, &checklist.Name, &checklist.Description, &checklist.TypeOfList)
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
	row := m.DB.QueryRowContext(ctx, queryCheckList, id)

	if err := row.Err(); err != nil {
		return nil, err
	}

	err := row.Scan(&checklist.ID, &checklist.Name, &checklist.Description, &checklist.TypeOfList)
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

	return &checklist, nil
}

func (m *postgresDBRepo) GetAllCheckListsOfType(id_type_of_list int) (*[]models.CheckList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var checklists []models.CheckList

	queryCheckLists := `
	select id_check_list, name_check_list, description, id_type_of_list
	from check_lists
	where id_type_of_list = $1
	`

	queryPointsCheckList := `
	select name_point 
	from check_list_points
	where id_check_list = $1
	`

	rows, err := m.DB.QueryContext(ctx, queryCheckLists, id_type_of_list)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var checklist models.CheckList

		err := rows.Scan(&checklist.ID, &checklist.Name, &checklist.Description, &checklist.TypeOfList)
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

		checklists = append(checklists, checklist)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &checklists, nil
}

// Delete check list by ID
func (m *postgresDBRepo) DeleteCheckListByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryCheckList := `
	delete 
	from check_lists
	where id_check_list = $1
	`

	queryPointsCheckList := `
	delete  
	from check_list_points
	where id_check_list = $1
	`
	tx, err := m.DB.Begin()
	if err != nil {
		fmt.Println("Error while starting a new transaction for bank account transaction: " + err.Error())
		return errors.New("Unexpected database error")
	}

	_, err = tx.ExecContext(ctx, queryCheckList, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryPointsCheckList, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (m *postgresDBRepo) InsertCheckList(checkList *models.CheckList) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `insert into check_lists (name_check_list, description, id_type_of_list)
	VALUES ($1, $2, $3) 
	RETURNING id_check_list
	`

	var id int
	if err := m.DB.QueryRowContext(ctx, query,
		checkList.Name,
		checkList.Description,
		checkList.TypeOfList).Scan(&id); err != nil {
		return err
	}

	fmt.Println(id)

	queryNames := `insert into check_list_points (id_check_list, name_point)
	VALUES ($1, $2)
	`

	for _, value := range checkList.NameOfPoints {
		_, err := m.DB.ExecContext(ctx, queryNames,
			id,
			value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *postgresDBRepo) UpdateCheckList(checkList *models.CheckList) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update check_lists set name_check_list = $1, description = $2
	where id_check_list = $3
	`

	_, err := m.DB.ExecContext(ctx, query,
		checkList.Name,
		checkList.Description,
		checkList.ID)
	if err != nil {
		return err
	}

	queryNamesDelete := `delete
	from check_list_points 
	where id_check_list = $1
	`

	_, err = m.DB.ExecContext(ctx, queryNamesDelete, checkList.ID)
	if err != nil {
		return err
	}

	queryNamesInsert := `insert into check_list_points (id_check_list, name_point)
	VALUES ($1, $2)
	`

	for _, value := range checkList.NameOfPoints {
		_, err := m.DB.ExecContext(ctx, queryNamesInsert,
			checkList.ID,
			value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
