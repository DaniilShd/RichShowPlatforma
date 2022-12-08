package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetMasterClassByID(id int) (*models.MasterClass, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var masterClass models.MasterClass

	queryMasterClass := `
	select id_master_class, name_master_class, description, id_check_list 
	from master_class
	where id_master_class = $1
	`

	rows := m.DB.QueryRowContext(ctx, queryMasterClass, id)

	var idChekList int
	err := rows.Scan(&masterClass.ID, &masterClass.Name, &masterClass.Description, &idChekList)
	if err != nil {
		return nil, err
	}

	var checkList *models.CheckList
	checkList, err = m.GetCheckListByID(idChekList)
	masterClass.CheckList = *checkList

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &masterClass, nil
}

func (m *postgresDBRepo) GetAllMasterClass() (*[]models.MasterClass, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var masterClasses []models.MasterClass

	queryMasterClass := `
	select id_master_class, name_master_class, description, id_check_list 
	from master_class
	`

	rows, err := m.DB.QueryContext(ctx, queryMasterClass)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var masterClass models.MasterClass
		var idChekList int
		err := rows.Scan(&masterClass.ID, &masterClass.Name, &masterClass.Description, &idChekList)
		if err != nil {
			return nil, err
		}

		var checkList *models.CheckList
		checkList, err = m.GetCheckListByID(idChekList)
		masterClass.CheckList = *checkList

		masterClasses = append(masterClasses, masterClass)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &masterClasses, nil
}

// Update masterclass in database
func (m *postgresDBRepo) UpdateMasterClass(masterClass *models.MasterClass) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update master_class set name_master_class = $1, description = $2, id_check_list = $3
	where id_master_class = $4
	`

	_, err := m.DB.ExecContext(ctx, query,
		masterClass.Name,
		masterClass.Description,
		masterClass.CheckList.ID,
		masterClass.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// delete master class by id
func (m *postgresDBRepo) DeleteMasterClassByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `delete from master_class
	where id_master_class = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil

} //Add new master-class
func (m *postgresDBRepo) InsertMasterClass(masterClass *models.MasterClass) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `insert into master_class (name_master_class, description, id_check_list)
	VALUES ($1, $2, $3)
	`

	_, err := m.DB.ExecContext(ctx, query,
		masterClass.Name,
		masterClass.Description,
		masterClass.CheckList.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
