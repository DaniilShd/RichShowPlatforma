package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) InsertLead(lead *models.Lead) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	idClient, idChild, err := m.insertClient(ctx, lead)
	if err != nil {
		return err
	}
	queryInsertLead := `insert into leads
	(id_client, amount_of_children, average_age_of_children, address, date, id_client_child)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	var idLead int
	if err := m.DB.QueryRowContext(ctx, queryInsertLead,
		idClient,
		lead.AmountOfChildren,
		lead.AverageAgeOfChildren,
		lead.Address,
		lead.Date,
		idChild).Scan(&idLead); err != nil {
		return err
	}
	err = m.insertPrograms(ctx, &lead.Programs, idLead)
	if err != nil {
		return err
	}
	err = m.insertHeroes(ctx, &lead.Heroes, idLead)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) insertPrograms(ctx context.Context, programs *[]models.Program, idLead int) error {
	queryInsertProgram := `insert into lead_program
	(id_check_list, id_lead)
	VALUES ($1, $2)
	`

	for _, program := range *programs {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, idLead)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *postgresDBRepo) insertHeroes(ctx context.Context, heroes *[]models.Hero, idLead int) error {
	queryInsertProgram := `insert into lead_heroes
	(id_hero, id_lead)
	VALUES ($1, $2)
	`

	for _, hero := range *heroes {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, hero.ID, idLead)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *postgresDBRepo) insertClient(ctx context.Context, lead *models.Lead) (int, int, error) {

	querySelectClient := `SELECT * from clients
	where phone_number = $1
	`
	var idClient int
	if err := m.DB.QueryRowContext(ctx, querySelectClient, lead.PhoneNumber).Scan(&idClient); err != nil {
		return 0, 0, err
	}

	//Если такой клиент уже есть в базе данных пропускаем запись клиента
	if idClient == 0 {
		queryInsertClient := `insert into clients
		(first_name, last_name, phone_number, telegram_client)
		VALUES ($1, $2, $3, $4)
		`
		if err := m.DB.QueryRowContext(ctx, queryInsertClient,
			lead.FirstNameClient,
			lead.LastNameClient,
			lead.PhoneNumber,
			lead.Telegram).Scan(&idClient); err != nil {
			return 0, 0, err
		}
	}
	lead.ID = idClient
	idChild, err := m.insertChilds(ctx, lead)
	if err != nil {
		return 0, 0, err
	}
	return idClient, idChild, nil
}

func (m *postgresDBRepo) insertChilds(ctx context.Context, lead *models.Lead) (int, error) {

	querySelectChild := `SELECT * from clients
	where name_child = $1 AND id_client = $2 AND date_of_birthday_child = $3
	`
	var idChild int
	if err := m.DB.QueryRowContext(ctx, querySelectChild, lead.Child.Name, lead.ID, lead.Child.DateOfBirthDay).Scan(&idChild); err != nil {
		return 0, err
	}

	//если такой ребенок уже есть в базе, записывать его снова не нужно, возвращаем id этого ребенка
	if idChild != 0 {
		return idChild, nil
	}
	queryInsertChild := `insert into client_child
		(name_child, date_of_birthday_child, id_client, id_gender_child)
		VALUES ($1, $2, $3, $4)
		`
	if err := m.DB.QueryRowContext(ctx, queryInsertChild,
		lead.Child.Name,
		lead.Child.DateOfBirthDay,
		lead.ID,
		lead.Child.Gender).Scan(&idChild); err != nil {
		return 0, err
	}
	return idChild, nil
}
