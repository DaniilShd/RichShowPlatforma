package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetPartyByID(id int) (*models.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var party models.Party

	queryparty := `
	select id_party_quest, name_party_quest, description, id_check_list 
	from parties_quests
	where id_party_quest = $1
	`

	rows := m.DB.QueryRowContext(ctx, queryparty, id)

	var idChekList int
	err := rows.Scan(&party.ID, &party.Name, &party.Description, &idChekList)
	if err != nil {
		return nil, err
	}

	var checkList *models.CheckList
	checkList, err = m.GetCheckListByID(idChekList)
	party.CheckList = *checkList

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &party, nil
}

func (m *postgresDBRepo) GetAllParty() (*[]models.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var partys []models.Party

	queryparty := `
	select id_party_quest, name_party_quest, description, id_check_list 
	from parties_quests
	`

	rows, err := m.DB.QueryContext(ctx, queryparty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var party models.Party
		var idChekList int
		err := rows.Scan(&party.ID, &party.Name, &party.Description, &idChekList)
		if err != nil {
			return nil, err
		}

		var checkList *models.CheckList
		checkList, err = m.GetCheckListByID(idChekList)
		party.CheckList = *checkList

		partys = append(partys, party)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &partys, nil
}

// Update party in database
func (m *postgresDBRepo) UpdateParty(party *models.Party) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update parties_quests set name_party_quest = $1, description = $2, id_check_list = $3
	where id_party_quest = $4
	`

	_, err := m.DB.ExecContext(ctx, query,
		party.Name,
		party.Description,
		party.CheckList.ID,
		party.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// delete master class by id
func (m *postgresDBRepo) DeletePartyByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `delete from parties_quests
	where id_party_quest = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil

} //Add new master-class
func (m *postgresDBRepo) InsertParty(party *models.Party) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `insert into parties_quests (name_party_quest, description, id_check_list)
	VALUES ($1, $2, $3)
	`

	_, err := m.DB.ExecContext(ctx, query,
		party.Name,
		party.Description,
		party.CheckList.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
