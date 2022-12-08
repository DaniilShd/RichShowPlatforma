package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetProgramShowByID(id int) (*models.ProgramShow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var programShow models.ProgramShow

	queryprogramShow := `
	select id_show, name_show, description, id_check_list 
	from program_show
	where id_show = $1
	`

	rows := m.DB.QueryRowContext(ctx, queryprogramShow, id)

	var idChekList int
	err := rows.Scan(&programShow.ID, &programShow.Name, &programShow.Description, &idChekList)
	if err != nil {
		return nil, err
	}

	var checkList *models.CheckList
	checkList, err = m.GetCheckListByID(idChekList)
	programShow.CheckList = *checkList

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &programShow, nil
}

func (m *postgresDBRepo) GetAllProgramShow() (*[]models.ProgramShow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var programShows []models.ProgramShow

	queryprogramShow := `
	select id_show, name_show, description, id_check_list 
	from program_show
	`

	rows, err := m.DB.QueryContext(ctx, queryprogramShow)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var programShow models.ProgramShow
		var idChekList int
		err := rows.Scan(&programShow.ID, &programShow.Name, &programShow.Description, &idChekList)
		if err != nil {
			return nil, err
		}

		var checkList *models.CheckList
		checkList, err = m.GetCheckListByID(idChekList)
		programShow.CheckList = *checkList

		programShows = append(programShows, programShow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &programShows, nil
}

// Update programShow in database
func (m *postgresDBRepo) UpdateProgramShow(programShow *models.ProgramShow) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update program_show set name_show = $1, description = $2, id_check_list = $3
	where id_show = $4
	`

	_, err := m.DB.ExecContext(ctx, query,
		programShow.Name,
		programShow.Description,
		programShow.CheckList.ID,
		programShow.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// delete master class by id
func (m *postgresDBRepo) DeleteProgramShowByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `delete from program_show
	where id_show = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil

} //Add new master-class
func (m *postgresDBRepo) InsertProgramShow(programShow *models.ProgramShow) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `insert into program_show (name_show, description, id_check_list)
	VALUES ($1, $2, $3)
	`

	_, err := m.DB.ExecContext(ctx, query,
		programShow.Name,
		programShow.Description,
		programShow.CheckList.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
