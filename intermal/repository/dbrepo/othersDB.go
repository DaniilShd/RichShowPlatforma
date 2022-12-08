package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetOtherByID(id int) (*models.Other, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var other models.Other

	queryother := `
	select id_other, name_other, description, id_check_list 
	from others
	where id_other = $1
	`

	rows := m.DB.QueryRowContext(ctx, queryother, id)

	var idChekList int
	err := rows.Scan(&other.ID, &other.Name, &other.Description, &idChekList)
	if err != nil {
		return nil, err
	}

	var checkList *models.CheckList
	checkList, err = m.GetCheckListByID(idChekList)
	other.CheckList = *checkList

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &other, nil
}

func (m *postgresDBRepo) GetAllOther() (*[]models.Other, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var others []models.Other

	queryother := `
	select id_others, name_other, description, id_check_list 
	from others
	`

	rows, err := m.DB.QueryContext(ctx, queryother)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var other models.Other
		var idChekList int
		err := rows.Scan(&other.ID, &other.Name, &other.Description, &idChekList)
		if err != nil {
			return nil, err
		}

		var checkList *models.CheckList
		checkList, err = m.GetCheckListByID(idChekList)
		other.CheckList = *checkList

		others = append(others, other)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &others, nil
}

// Update other in database
func (m *postgresDBRepo) UpdateOther(other *models.Other) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update others set name_other = $1, description = $2, id_check_list = $3
	where id_other = $4
	`

	_, err := m.DB.ExecContext(ctx, query,
		other.Name,
		other.Description,
		other.CheckList.ID,
		other.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// delete master class by id
func (m *postgresDBRepo) DeleteOtherByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `delete from others
	where id_other = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil

} //Add new master-class
func (m *postgresDBRepo) InsertOther(other *models.Other) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `insert into others (name_other, description, id_check_list)
	VALUES ($1, $2, $3)
	`

	_, err := m.DB.ExecContext(ctx, query,
		other.Name,
		other.Description,
		other.CheckList.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
