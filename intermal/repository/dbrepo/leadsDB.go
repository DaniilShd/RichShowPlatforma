package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) InsertLead(lead *models.Lead) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	idClient, err := m.insertClient(ctx, &lead.Client)
	if err != nil {
		return err
	}
	fmt.Println(idClient)
	idChild, err := m.insertChild(ctx, &lead.Child, idClient)
	if err != nil {
		return err
	}
	fmt.Println(idChild)
	queryInsertLead := `insert into leads
	(id_client, amount_of_children, average_age_of_children, address, date, time, id_client_child)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	returning id_lead
	`

	var idLead int
	if err := m.DB.QueryRowContext(ctx, queryInsertLead,
		idClient,
		lead.AmountOfChildren,
		lead.AverageAgeOfChildren,
		lead.Address,
		lead.Date,
		lead.Time,
		idChild).Scan(&idLead); err != nil {

		return err
	}
	if lead.MasterClasses != nil {
		err = m.insertMasterCLass(ctx, &lead.MasterClasses, idLead)
		if err != nil {

			return err

		}
	}
	if lead.Shows != nil {
		err = m.insertShows(ctx, &lead.Shows, idLead)
		if err != nil {

			return err
		}
	}
	if lead.PartyAndQuests != nil {
		err = m.insertPartyAndQuest(ctx, &lead.PartyAndQuests, idLead)
		if err != nil {

			return err
		}
	}
	if lead.Animations != nil {
		err = m.insertAnimation(ctx, &lead.Animations, idLead)
		if err != nil {

			return err
		}
	}
	err = m.insertOthers(ctx, &lead.Others, idLead)
	if err != nil {
		return err
	}
	err = m.insertHeroes(ctx, &lead.Heroes, idLead)
	if err != nil {
		return err
	}
	err = m.insertAssistants(ctx, &lead.Assistants, idLead)
	if err != nil {
		return err
	}
	return nil
}

// Пакет функций для добавления рограмм в заказ (Лид) Start-------------------------------------------------------------------------
func (m *postgresDBRepo) insertMasterCLass(ctx context.Context, programs *[]models.MasterClass, idLead int) error {
	var ID []int
	for _, item := range *programs {
		ID = append(ID, item.ID)
	}
	return m.insertPrograms(ctx, ID, idLead)
}

func (m *postgresDBRepo) insertShows(ctx context.Context, programs *[]models.Show, idLead int) error {
	var ID []int
	for _, item := range *programs {
		ID = append(ID, item.ID)
	}
	return m.insertPrograms(ctx, ID, idLead)
}

func (m *postgresDBRepo) insertPartyAndQuest(ctx context.Context, programs *[]models.PartyAndQuest, idLead int) error {
	var ID []int
	for _, item := range *programs {
		ID = append(ID, item.ID)
	}
	return m.insertPrograms(ctx, ID, idLead)
}

func (m *postgresDBRepo) insertAnimation(ctx context.Context, programs *[]models.Animation, idLead int) error {
	var ID []int
	for _, item := range *programs {
		ID = append(ID, item.ID)
	}
	return m.insertPrograms(ctx, ID, idLead)
}

func (m *postgresDBRepo) insertOthers(ctx context.Context, programs *[]models.Other, idLead int) error {
	var ID []int
	for _, item := range *programs {
		ID = append(ID, item.ID)
	}
	return m.insertPrograms(ctx, ID, idLead)
}

func (m *postgresDBRepo) insertPrograms(ctx context.Context, programs []int, idLead int) error {
	queryInsertProgram := `insert into lead_programs
	(id_check_list, id_lead)
	VALUES ($1, $2)
	`

	for _, program := range programs {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, program, idLead)
		if err != nil {
			return err
		}
	}
	return nil
}

// Пакет функций для добавления рограмм в заказ (Лид) End-------------------------------------------------------------------------

func (m *postgresDBRepo) insertHeroes(ctx context.Context, heroes *[]models.LeadHero, idLead int) error {
	queryInsertHeroes := `insert into lead_heroes
	(id_hero, id_lead, id_artist)
	VALUES ($1, $2, $3)
	`

	for _, hero := range *heroes {
		_, err := m.DB.ExecContext(ctx, queryInsertHeroes, hero.HeroID, idLead, hero.ArtistID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *postgresDBRepo) insertAssistants(ctx context.Context, assistants *[]models.Assistant, idLead int) error {
	queryInsertAssistants := `insert into lead_assistants
	(id_assistant, id_lead)
	VALUES ($1, $2)
	`

	for _, assistant := range *assistants {
		_, err := m.DB.ExecContext(ctx, queryInsertAssistants, assistant.ID, idLead)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *postgresDBRepo) insertClient(ctx context.Context, client *models.Client) (int, error) {

	querySelectClient := `SELECT id_client from clients
	where phone_number = $1
	`
	var idClient int
	row := m.DB.QueryRowContext(ctx, querySelectClient, client.PhoneNumber)
	if err := row.Scan(&idClient); err != nil {
		idClient = 0
	}

	//Если такой клиент уже есть в базе данных пропускаем запись клиента
	if idClient == 0 {
		queryInsertClient := `insert into clients
		(first_name, last_name, phone_number, telegram_client)
		VALUES ($1, $2, $3, $4)
		`
		if err := m.DB.QueryRowContext(ctx, queryInsertClient,
			client.FirstName,
			client.LastName,
			client.PhoneNumber,
			client.Telegram).Scan(&idClient); err != nil {
			return 0, err
		}
	}

	return idClient, nil
}

func (m *postgresDBRepo) insertChild(ctx context.Context, child *models.Child, idClient int) (int, error) {

	querySelectChild := `SELECT id_child from clients
	where name_child = $1 AND id_client = $2 AND date_of_birthday_child = $3
	`
	var idChild int
	row := m.DB.QueryRowContext(ctx, querySelectChild, child.Name, idClient, child.DateOfBirthDay)

	if err := row.Scan(&idChild); err == nil {
		return idChild, nil
	}
	//если такой ребенок уже есть в базе, записывать его снова не нужно, возвращаем id этого ребенка

	queryInsertChild := `insert into client_child
		(name_child, date_of_birthday_child, id_client, age, id_gender_child)
		VALUES ($1, $2, $3, $4, $5)
		Returning id_client_child
		`
	if err := m.DB.QueryRowContext(ctx, queryInsertChild,
		child.Name,
		child.DateOfBirthDay,
		idClient,
		child.Age,
		child.Gender).Scan(&idChild); err != nil {
		return 0, err
	}
	return idChild, nil
}

//Get all leads -----------------------------------------------------------------------------------------------------------------

func (m *postgresDBRepo) GetAllLeads() ([]models.Lead, error) {

}
