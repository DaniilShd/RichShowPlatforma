package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	modelsTelegram "github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"

	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) InsertLead(lead *models.Lead) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	idClient, err := m.insertClient(ctx, &lead.Client)
	if err != nil {
		return err
	}

	idChild, err := m.insertChild(ctx, &lead.Child, idClient)
	if err != nil {
		return err
	}

	queryInsertLead := `insert into leads
	(id_client, amount_of_children, average_age_of_children, address, date, time, id_client_child, description, duration)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
		idChild,
		lead.Description,
		lead.Duration).Scan(&idLead); err != nil {

		return err
	}

	lead.ID = idLead

	err = m.insertPrograms(ctx, lead, idLead)
	if err != nil {
		return err
	}
	err = m.insertHeroes(ctx, &lead.Heroes, lead)
	if err != nil {
		return err
	}
	err = m.insertAssistants(ctx, &lead.Assistants, lead)
	if err != nil {
		return err
	}

	checkArtist, err := m.checkArtist(ctx, idLead)
	if err != nil {
		return err
	}
	queryCheckArtist := `update leads
	set check_artists=$1
	where id_lead=$2
	`
	_, err = m.DB.ExecContext(ctx, queryCheckArtist, checkArtist, idLead)
	if err != nil {
		return err
	}

	checkAssistant, err := m.checkAssistant(ctx, idLead)
	if err != nil {
		return err
	}
	queryCheckAssistant := `update leads
	set check_assistants=$1
	where id_lead=$2
	`
	_, err = m.DB.ExecContext(ctx, queryCheckAssistant, checkAssistant, idLead)
	if err != nil {
		return err
	}

	return nil
}

// Пакет функций для добавления рограмм в заказ (Лид) Start-------------------------------------------------------------------------

func (m *postgresDBRepo) insertPrograms(ctx context.Context, programs *models.Lead, idLead int) error {
	queryInsertProgram := `insert into lead_programs
	(id_check_list, description, id_lead)
	VALUES ($1, $2, $3)
	`

	queryInsertOrderStore := `
	insert into order_store
	(id_check_list, id_lead)
	Values ($1, $2)
	`
	for _, program := range programs.MasterClasses {
		if program.ID != 0 {
			_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
			if err != nil {
				return err
			}
			_, err = m.DB.ExecContext(ctx, queryInsertOrderStore, program.ID, idLead)
			if err != nil {
				return err
			}
		}
	}
	for _, program := range programs.Shows {
		if program.ID != 0 {
			_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
			if err != nil {
				return err
			}
			_, err = m.DB.ExecContext(ctx, queryInsertOrderStore, program.ID, idLead)
			if err != nil {
				return err
			}
		}
	}
	for _, program := range programs.PartyAndQuests {
		if program.ID != 0 {
			_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
			if err != nil {
				return err
			}
			_, err = m.DB.ExecContext(ctx, queryInsertOrderStore, program.ID, idLead)
			if err != nil {
				return err
			}
		}
	}
	for _, program := range programs.Animations {
		if program.ID != 0 {
			_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
			if err != nil {
				return err
			}
			_, err = m.DB.ExecContext(ctx, queryInsertOrderStore, program.ID, idLead)
			if err != nil {
				return err
			}
		}
	}
	for _, program := range programs.Others {
		if program.ID != 0 {
			_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
			if err != nil {
				return err
			}
			_, err = m.DB.ExecContext(ctx, queryInsertOrderStore, program.ID, idLead)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Пакет функций для добавления рограмм в заказ (Лид) End-------------------------------------------------------------------------

func (m *postgresDBRepo) insertHeroes(ctx context.Context, heroes *[]models.LeadHero, lead *models.Lead) error {

	queryInsertHeroes := `insert into lead_heroes
	(id_hero, id_lead)
	VALUES ($1, $2)
	`

	queryInsertArtist := `update lead_heroes
	set id_artist=$1
	where id_lead=$2
	`

	for _, hero := range *heroes {

		_, err := m.DB.ExecContext(ctx, queryInsertHeroes, hero.HeroID, lead.ID)
		if err != nil {
			return err
		}
		if hero.ArtistID != 0 {
			_, err := m.DB.ExecContext(ctx, queryInsertArtist, hero.ArtistID, lead.ID)
			if err != nil {
				return err
			}
		}

		//Отправка сообщения артистам
		chatID, err := m.getArtistChatID(ctx, hero.ArtistID)
		if err != nil {
			return err
		}
		if chatID != 0 && hero.NeedSendMessage {
			text := "Вы назначены на заказ № " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Дата: " + "<strong>" + string(lead.Date.Format("02-01-2006")) + "</strong>" + "\n" + "Время: " + "<strong>" + string(lead.Time.Format("15:04")) + "</strong>"

			m.App.MailChan <- modelsTelegram.MailData{
				ChatID: chatID,
				Text:   text,
			}
		}
	}
	return nil
}

func (m *postgresDBRepo) getArtistChatID(ctx context.Context, artistID int) (int64, error) {
	querySelect := `select id_telegram_chat
	from artists
	where id_artist=$1
	`

	var chatID sql.NullInt64
	var result int64
	err := m.DB.QueryRowContext(ctx, querySelect, artistID).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	if chatID.Valid {
		result = chatID.Int64
	} else {
		return 0, nil
	}
	return result, nil
}

func (m *postgresDBRepo) getStoreChatID(ctx context.Context) (int64, error) {
	querySelect := `select id_telegram_chat
	from id_accounts
	where access_level=3
	`

	var chatID sql.NullInt64
	var result int64
	err := m.DB.QueryRowContext(ctx, querySelect).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	if chatID.Valid {
		result = chatID.Int64
	} else {
		return 0, nil
	}
	return result, nil
}

func (m *postgresDBRepo) insertAssistants(ctx context.Context, assistants *[]models.Assistant, lead *models.Lead) error {
	queryInsertAssistants := `insert into lead_assistants
	(id_assistant, id_lead)
	VALUES ($1, $2)
	`

	for _, assistant := range *assistants {
		_, err := m.DB.ExecContext(ctx, queryInsertAssistants, assistant.ID, lead.ID)
		if err != nil {
			return err
		}
		//Отправка сообщения ассистентам
		chatID, err := m.getAssistantChatID(ctx, assistant.ID)
		if err != nil {
			return err
		}
		if chatID != 0 && assistant.NeedSendMessage {
			text := "Вы назначены на заказ № " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Дата: " + "<strong>" + string(lead.Date.Format("02-01-2006")) + "</strong>" + "\n" + "Время: " + "<strong>" + string(lead.Time.Format("15:04")) + "</strong>"

			m.App.MailChan <- modelsTelegram.MailData{
				ChatID: chatID,
				Text:   text,
			}
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
		(name_child, date_of_birthday_child, id_client, id_gender_child)
		VALUES ($1, $2, $3, $4)
		Returning id_client_child
		`
	if err := m.DB.QueryRowContext(ctx, queryInsertChild,
		child.Name,
		child.DateOfBirthDay,
		idClient,
		child.Gender).Scan(&idChild); err != nil {
		return 0, err
	}
	return idChild, nil
}

func (m *postgresDBRepo) checkArtist(ctx context.Context, idLead int) (bool, error) {
	querySelect := `
	select id_lead_hero, id_artist
	from lead_heroes
	where id_lead = $1
	`

	rows, err := m.DB.QueryContext(ctx, querySelect, idLead)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return false, err
	}
	var idArtists []int
	fmt.Println("check!!!", idLead)
	fmt.Println("check!!!", rows)

	for rows.Next() {
		var idArtist sql.NullInt64
		var idLeadHero int
		err := rows.Scan(&idLeadHero, &idArtist)
		if err != nil {
			return false, err
		}
		if !idArtist.Valid {
			return false, nil
		} else {
			idArtists = append(idArtists, int(idArtist.Int64))
		}
	}
	if len(idArtists) != 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgresDBRepo) checkAssistant(ctx context.Context, idLead int) (bool, error) {
	querySelect := `
	select id_assistant
	from lead_assistants
	where id_lead = $1
	`

	rows, err := m.DB.QueryContext(ctx, querySelect, idLead)
	if err != nil {
		return true, nil
	}

	var idAssistants []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return false, err
		}
		idAssistants = append(idAssistants, id)
	}
	if len(idAssistants) != 0 {
		return true, nil
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	return false, nil
}

//Get all leads --------------------------------------------------------------------------------------------------------------------------

func (m *postgresDBRepo) GetAllRawLeads() ([]models.Lead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryRawLeads := `
	select id_lead, id_client, amount_of_children, average_age_of_children, 
	address, date, check_artists, confirmed, check_assistants, id_client_child, time, description, duration
	from leads
	where (date + time)>(current_time+current_date)  AND (check_artists<>true OR confirmed<>true OR check_assistants<>true)
	`

	rows, err := m.DB.QueryContext(ctx, queryRawLeads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	var leads []models.Lead
	var timeDate string
	for rows.Next() {
		var lead models.Lead
		err := rows.Scan(&lead.ID,
			&lead.Client.ID,
			&lead.AmountOfChildren,
			&lead.AverageAgeOfChildren,
			&lead.Address,
			&lead.Date,
			&lead.CheckArtists,
			&lead.Confirmed,
			&lead.CheckAssistants,
			&lead.Child.ID,
			&timeDate,
			&lead.Description,
			&lead.Duration,
		)
		if err != nil {
			return nil, err
		}

		lead.Time, err = time.Parse("15:04", timeDate[:5])
		if err != nil {
			return nil, err
		}

		client, err := m.getClientByID(ctx, lead.Client.ID)
		if err != nil {
			return nil, err
		}
		lead.Client.FirstName = client.FirstName
		lead.Client.LastName = client.LastName

		lead.Client.PhoneNumber = helpers.ConvertNumberPhone(client.PhoneNumber)

		leads = append(leads, lead)
	}
	return leads, nil
}

func (m *postgresDBRepo) GetAllConfirmedLeads() ([]models.Lead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryLeads := `
	select id_lead, id_client, amount_of_children, average_age_of_children, 
	address, date, check_artists, confirmed, check_assistants, id_client_child, time, description, duration
	from leads
	where (date + time)>(current_time+ current_date) AND (check_artists=true AND confirmed=true AND check_assistants=true)
	`

	rows, err := m.DB.QueryContext(ctx, queryLeads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	var leads []models.Lead
	var timeDate string
	for rows.Next() {
		var lead models.Lead
		err := rows.Scan(&lead.ID,
			&lead.Client.ID,
			&lead.AmountOfChildren,
			&lead.AverageAgeOfChildren,
			&lead.Address,
			&lead.Date,
			&lead.CheckArtists,
			&lead.Confirmed,
			&lead.CheckAssistants,
			&lead.Child.ID,
			&timeDate,
			&lead.Description,
			&lead.Duration,
		)
		if err != nil {
			return nil, err
		}

		lead.Time, err = time.Parse("15:04", timeDate[:5])
		if err != nil {
			return nil, err
		}

		client, err := m.getClientByID(ctx, lead.Client.ID)
		if err != nil {
			return nil, err
		}
		lead.Client.FirstName = client.FirstName
		lead.Client.LastName = client.LastName

		lead.Client.PhoneNumber = helpers.ConvertNumberPhone(client.PhoneNumber)

		leads = append(leads, lead)
	}
	return leads, nil
}

func (m *postgresDBRepo) GetAllArchiveLeads() ([]models.Lead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryLeads := `
	select id_lead, id_client, amount_of_children, average_age_of_children, 
	address, date, check_artists, confirmed, check_assistants, id_client_child, time, description, duration
	from leads
	where (date + time)<(current_time+ current_date)
	`

	rows, err := m.DB.QueryContext(ctx, queryLeads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	var leads []models.Lead
	var timeDate string
	for rows.Next() {
		var lead models.Lead
		err := rows.Scan(&lead.ID,
			&lead.Client.ID,
			&lead.AmountOfChildren,
			&lead.AverageAgeOfChildren,
			&lead.Address,
			&lead.Date,
			&lead.CheckArtists,
			&lead.Confirmed,
			&lead.CheckAssistants,
			&lead.Child.ID,
			&timeDate,
			&lead.Description,
			&lead.Duration,
		)
		if err != nil {
			return nil, err
		}

		lead.Time, err = time.Parse("15:04", timeDate[:5])
		if err != nil {
			return nil, err
		}

		client, err := m.getClientByID(ctx, lead.Client.ID)
		if err != nil {
			return nil, err
		}
		lead.Client.FirstName = client.FirstName
		lead.Client.LastName = client.LastName

		lead.Client.PhoneNumber = helpers.ConvertNumberPhone(client.PhoneNumber)

		leads = append(leads, lead)
	}
	return leads, nil
}

func (m *postgresDBRepo) GetLeadByID(idLead int) (*models.Lead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryLeadID := `
	select id_lead, id_client, amount_of_children, average_age_of_children, 
	address, date, check_artists, confirmed, check_assistants, id_client_child, time, description, duration
	from leads
	where id_lead=$1
	`
	var lead models.Lead
	var timeDate string
	err := m.DB.QueryRowContext(ctx, queryLeadID, idLead).Scan(&lead.ID,
		&lead.Client.ID,
		&lead.AmountOfChildren,
		&lead.AverageAgeOfChildren,
		&lead.Address,
		&lead.Date,

		&lead.CheckArtists,
		&lead.Confirmed,
		&lead.CheckAssistants,
		&lead.Child.ID,
		&timeDate,
		&lead.Description,
		&lead.Duration,
	)
	if err != nil {
		return nil, err
	}

	lead.Time, err = time.Parse("15:04", timeDate[:5])
	if err != nil {
		return nil, err
	}

	client, err := m.getClientByID(ctx, lead.Client.ID)
	if err != nil {
		return nil, err
	}
	lead.Client.FirstName = client.FirstName
	lead.Client.LastName = client.LastName
	lead.Client.Telegram = client.Telegram
	lead.Client.PhoneNumber = helpers.ConvertNumberPhone(client.PhoneNumber)

	child, err := m.getChildByID(ctx, lead.Child.ID)
	if err != nil {
		return nil, err
	}
	lead.Child.Name = child.Name
	lead.Child.DateOfBirthDay = child.DateOfBirthDay
	lead.Child.Gender = child.Gender

	programs, err := m.getProgramsByLeadID(ctx, lead.ID)
	if err != nil {
		return nil, err
	}
	for _, program := range programs[CHECK_LISTS_TYPE_OF_MASTER_CLASS] {
		lead.MasterClasses = append(lead.MasterClasses, program.(models.MasterClass))
	}
	for _, program := range programs[CHECK_LISTS_TYPE_OF_ANIMATION] {
		lead.Animations = append(lead.Animations, program.(models.Animation))
	}
	for _, program := range programs[CHECK_LISTS_TYPE_OF_OTHER] {
		lead.Others = append(lead.Others, program.(models.Other))
	}
	for _, program := range programs[CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS] {
		lead.PartyAndQuests = append(lead.PartyAndQuests, program.(models.PartyAndQuest))
	}
	for _, program := range programs[CHECK_LISTS_TYPE_OF_SHOW] {
		lead.Shows = append(lead.Shows, program.(models.Show))
	}

	lead.Heroes, err = m.getHeroesByID(ctx, lead.ID)
	if err != nil {
		return nil, err
	}

	lead.Assistants, err = m.getLeadAssistants(ctx, lead.ID)
	if err != nil {
		return nil, err
	}

	return &lead, nil
}

func (m *postgresDBRepo) getClientByID(ctx context.Context, idClient int) (*models.Client, error) {
	queryClient := `
	select first_name, last_name, phone_number, telegram_client
	from clients
	where id_client = $1
	`

	var client models.Client
	err := m.DB.QueryRowContext(ctx, queryClient, idClient).Scan(&client.FirstName,
		&client.LastName,
		&client.PhoneNumber,
		&client.Telegram)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (m *postgresDBRepo) getChildByID(ctx context.Context, idChild int) (*models.Child, error) {
	queryChild := `
	select name_child, date_of_birthday_child, id_gender_child
	from client_child
	where id_client_child = $1
	`

	var child models.Child

	err := m.DB.QueryRowContext(ctx, queryChild, idChild).Scan(
		&child.Name,
		&child.DateOfBirthDay,
		&child.Gender)
	if err != nil {
		return nil, err
	}

	return &child, nil
}

func (m *postgresDBRepo) getHeroesByID(ctx context.Context, idLead int) ([]models.LeadHero, error) {
	queryLeadHeroes := `
	select id_hero, id_artist
	from lead_heroes
	where id_lead = $1
	`

	rows, err := m.DB.QueryContext(ctx, queryLeadHeroes, idLead)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var heroes []models.LeadHero
	for rows.Next() {
		var hero models.Hero
		var artist models.Artist
		var heroLead models.LeadHero
		var artistID sql.NullInt64
		var artistPhone sql.NullString
		err := rows.Scan(&hero.ID, &artistID)
		if err != nil {
			return nil, err
		}

		if artistID.Valid {
			artist.ID = int(artistID.Int64)

			queryArtist := `
			select first_name, last_name, phone_number
			from artists
			where id_artist=$1
			`
			err := m.DB.QueryRowContext(ctx, queryArtist, artist.ID).Scan(
				&artist.FirstName,
				&artist.LastName,
				&artistPhone,
			)
			if err != nil {
				return nil, err
			}
			heroLead.ArtistFirstName = artist.FirstName
			heroLead.ArtistID = artist.ID
			heroLead.ArtistLastName = artist.LastName
			if artistPhone.Valid {
				heroLead.PhoneNumber = artistPhone.String
			}

		}

		queryHero := `
		select name_hero
		from heroes
		where id_hero=$1
		`
		err = m.DB.QueryRowContext(ctx, queryHero, hero.ID).Scan(
			&hero.Name)
		if err != nil {
			return nil, err
		}
		heroLead.HeroID = hero.ID
		heroLead.HeroName = hero.Name

		heroes = append(heroes, heroLead)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return heroes, nil
}

func (m *postgresDBRepo) getProgramsByLeadID(ctx context.Context, idLead int) (map[int][]interface{}, error) {
	querySelect := `
	select id_check_list, description
	from lead_programs
	where id_lead=$1
	`

	result := make(map[int][]interface{}, 5)

	rows, err := m.DB.QueryContext(ctx, querySelect, idLead)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var program models.Program
		var desc sql.NullString
		err := rows.Scan(&program.CheckListID, &desc)
		if err != nil {
			return nil, err
		}
		if desc.Valid {
			program.Description = desc.String
		}
		programs = append(programs, program)
	}

	var checkLists []models.CheckList
	for _, program := range programs {
		checkList, err := m.GetCheckListByID(program.CheckListID)
		if err != nil {
			return nil, err
		}
		checkLists = append(checkLists, *checkList)
		switch checkList.TypeOfList {
		case CHECK_LISTS_TYPE_OF_SHOW:
			var show models.Show
			show.ID = checkList.ID
			show.Description = program.Description
			show.Duration = checkList.Duration
			show.Name = checkList.Name
			result[CHECK_LISTS_TYPE_OF_SHOW] = append(result[CHECK_LISTS_TYPE_OF_SHOW], show)
		case CHECK_LISTS_TYPE_OF_MASTER_CLASS:
			var master_class models.MasterClass
			master_class.ID = checkList.ID
			master_class.Description = program.Description
			master_class.Duration = checkList.Duration
			master_class.Name = checkList.Name
			result[CHECK_LISTS_TYPE_OF_MASTER_CLASS] = append(result[CHECK_LISTS_TYPE_OF_MASTER_CLASS], master_class)
		case CHECK_LISTS_TYPE_OF_ANIMATION:
			var animation models.Animation
			animation.ID = checkList.ID
			animation.Description = program.Description
			animation.Duration = checkList.Duration
			animation.Name = checkList.Name
			result[CHECK_LISTS_TYPE_OF_ANIMATION] = append(result[CHECK_LISTS_TYPE_OF_ANIMATION], animation)
		case CHECK_LISTS_TYPE_OF_OTHER:
			var other models.Other
			other.ID = checkList.ID
			other.Description = program.Description
			other.Duration = checkList.Duration
			other.Name = checkList.Name
			result[CHECK_LISTS_TYPE_OF_OTHER] = append(result[CHECK_LISTS_TYPE_OF_OTHER], other)
		case CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS:
			var party models.PartyAndQuest
			party.ID = checkList.ID
			party.Description = program.Description
			party.Duration = checkList.Duration
			party.Name = checkList.Name
			result[CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS] = append(result[CHECK_LISTS_TYPE_OF_PARTIES_AND_QUESTS], party)
		}
	}
	return result, nil
}

func (m *postgresDBRepo) getLeadAssistants(ctx context.Context, idLead int) ([]models.Assistant, error) {
	querySelect := `
	select id_assistant
	from lead_assistants
	where id_lead=$1
	`
	rows, err := m.DB.QueryContext(ctx, querySelect, idLead)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assistantIDs []int
	var assistants []models.Assistant
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		assistantIDs = append(assistantIDs, id)
	}
	for _, id := range assistantIDs {
		assistant, err := m.GetAssistantByID(id)
		if err != nil {
			return nil, err
		}
		assistants = append(assistants, *assistant)
	}
	return assistants, nil
}

//Set and delete confirmed -------------------------------------------------------------------------------------------------------------------------\

func (m *postgresDBRepo) SetConfirmedLeadByID(idLead int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsert := `
	update leads
	set confirmed=$1
	where id_lead=$2
	`

	_, err := m.DB.ExecContext(ctx, queryInsert, true, idLead)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteConfirmedLeadByID(idLead int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsert := `
	update leads
	set confirmed=$1
	where id_lead=$2
	`

	_, err := m.DB.ExecContext(ctx, queryInsert, false, idLead)
	if err != nil {
		return err
	}
	return nil
}

//Delete lead by id------------------------------------------------------------------------------------------------------------------------------------

func (m *postgresDBRepo) DeleteLeadByID(idLead int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryDeleteOrderStore := `
	delete 
	from order_store
	where id_lead=$1 and completed<>true
	`

	queryDisassembleBagOrderStore := `
	update order_store
	set disassemble_bag=true
	where id_lead=$1 and completed=true
	`

	queryDeleteLead := `
	delete 
	from leads 
	where id_lead=$1
	`
	queryDeleteAssistant := `
	delete 
	from lead_assistants 
	where id_lead=$1
	`

	queryDeleteHeroes := `
	delete 
	from lead_heroes 
	where id_lead=$1
	`

	queryDeletePrograms := `
	delete 
	from lead_programs 
	where id_lead=$1
	`

	_, err := m.DB.ExecContext(ctx, queryDeleteOrderStore, idLead)
	if err != nil {
		return err
	}
	_, err = m.DB.ExecContext(ctx, queryDisassembleBagOrderStore, idLead)
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(ctx, queryDeleteLead, idLead)
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(ctx, queryDeleteAssistant, idLead)
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(ctx, queryDeleteHeroes, idLead)
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(ctx, queryDeletePrograms, idLead)
	if err != nil {
		return err
	}

	return nil

}

// Update lead start______________________________________________________________________________________
func (m *postgresDBRepo) UpdateLead(lead *models.Lead) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.updateClient(ctx, &lead.Client)
	if err != nil {
		return err
	}

	err = m.updateChild(ctx, &lead.Child)
	if err != nil {
		return err
	}

	LeadInfo, err := m.GetLeadByID(lead.ID)
	if err != nil {
		return err
	}

	if LeadInfo.Address != lead.Address {
		for _, assistant := range LeadInfo.Assistants {
			chatID, err := m.getAssistantChatID(ctx, assistant.ID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новый адрес: " + "<strong>" + lead.Address + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
		for _, heroArtist := range LeadInfo.Heroes {
			chatID, err := m.getArtistChatID(ctx, heroArtist.ArtistID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новый адрес: " + "<strong>" + lead.Address + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
	}

	if LeadInfo.Date != lead.Date {
		for _, assistant := range LeadInfo.Assistants {
			chatID, err := m.getAssistantChatID(ctx, assistant.ID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новая дата: " + "<strong>" + string(lead.Date.Format("02-01-2006")) + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
		for _, heroArtist := range LeadInfo.Heroes {
			chatID, err := m.getArtistChatID(ctx, heroArtist.ArtistID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новая дата: " + "<strong>" + string(lead.Date.Format("02-01-2006")) + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
		chatID, err := m.getStoreChatID(ctx)
		if err != nil {
			return err
		}
		if chatID != 0 {
			text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новая дата: " + "<strong>" + string(lead.Date.Format("02-01-2006")) + "</strong>"

			m.App.MailChan <- modelsTelegram.MailData{
				ChatID: chatID,
				Text:   text,
			}
		}
	}

	if LeadInfo.Time != lead.Time {
		for _, assistant := range LeadInfo.Assistants {
			chatID, err := m.getAssistantChatID(ctx, assistant.ID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новое время: " + "<strong>" + string(lead.Time.Format("15:04")) + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
		for _, heroArtist := range LeadInfo.Heroes {
			chatID, err := m.getArtistChatID(ctx, heroArtist.ArtistID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новое время: " + "<strong>" + string(lead.Time.Format("15:04")) + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
		chatID, err := m.getStoreChatID(ctx)
		if err != nil {
			return err
		}
		if chatID != 0 {
			text := "<strong>Изменеиние в заказе №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Новое время: " + "<strong>" + string(lead.Time.Format("15:04")) + "</strong>"

			m.App.MailChan <- modelsTelegram.MailData{
				ChatID: chatID,
				Text:   text,
			}
		}
	}

	queryUpdateLead := `update leads
	set amount_of_children=$1, average_age_of_children=$2, address=$3, date=$4, time=$5, description=$6, duration=$7
	where id_lead=$8
	`

	_, err = m.DB.ExecContext(ctx, queryUpdateLead,
		lead.AmountOfChildren,
		lead.AverageAgeOfChildren,
		lead.Address,
		lead.Date,
		lead.Time,
		lead.Description,
		lead.Duration,
		lead.ID)
	if err != nil {
		return err
	}

	err = m.updatePrograms(ctx, lead)
	if err != nil {
		return err
	}
	err = m.updateHeroes(ctx, lead)
	if err != nil {
		return err
	}
	err = m.updateAssistants(ctx, lead)
	if err != nil {
		return err
	}

	checkArtist, err := m.checkArtist(ctx, lead.ID)
	if err != nil {
		return err
	}
	queryCheckArtist := `update leads
	set check_artists=$1
	where id_lead=$2
	`
	_, err = m.DB.ExecContext(ctx, queryCheckArtist, checkArtist, lead.ID)
	if err != nil {
		return err
	}

	checkAssistant, err := m.checkAssistant(ctx, lead.ID)
	if err != nil {
		return err
	}
	queryCheckAssistant := `update leads
	set check_assistants=$1
	where id_lead=$2
	`
	_, err = m.DB.ExecContext(ctx, queryCheckAssistant, checkAssistant, lead.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) updateChild(ctx context.Context, child *models.Child) error {
	queryUpdateChild := `
	update client_child
	set name_child=$1, date_of_birthday_child=$2, id_gender_child=$3
	where id_client_child=$4
	`
	_, err := m.DB.ExecContext(ctx, queryUpdateChild,
		child.Name,
		child.DateOfBirthDay,
		child.Gender,
		child.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) updateClient(ctx context.Context, client *models.Client) error {
	queryUpdateChild := `
	update clients
	set first_name=$1, last_name=$2, phone_number=$3, telegram_client=$4
	where id_client=$5
	`
	_, err := m.DB.ExecContext(ctx, queryUpdateChild,
		client.FirstName,
		client.LastName,
		client.PhoneNumber,
		client.Telegram,
		client.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) updateAssistants(ctx context.Context, lead *models.Lead) error {
	assistants, err := m.getAssistantByID(ctx, lead.ID)
	if err != nil {
		return err
	}

	for i := range lead.Assistants {
		lead.Assistants[i].NeedSendMessage = true
	}

	for i := range assistants {
		assistants[i].Canceled = true
	}

	//Сравниваем артистов которые уже были занесены в таблицу с новыми артистами
	for i := range assistants {
		for y := range lead.Assistants {
			if assistants[i].ID == lead.Assistants[y].ID {
				assistants[i].Canceled = false
				lead.Assistants[y].NeedSendMessage = false
			}
		}
	}

	for _, assistantTable := range assistants {
		if assistantTable.ID != 0 && assistantTable.Canceled {
			chatID, err := m.getAssistantChatID(ctx, assistantTable.ID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Вас сняли с заказа №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Дата: " + "<strong>" + string(lead.Date.Format("02-01-2006")) + "</strong>" + "\n" + "Время: " + "<strong>" + string(lead.Time.Format("15:04")) + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
	}

	queryDelete := `
	delete 
	from lead_assistants
	where id_lead=$1
	`
	_, err = m.DB.ExecContext(ctx, queryDelete, lead.ID)
	if err != nil {
		return err
	}

	err = m.insertAssistants(ctx, &lead.Assistants, lead)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) getAssistantChatID(ctx context.Context, assistantID int) (int64, error) {
	querySelect := `select id_telegram_chat
	from assistants
	where id_assistant=$1
	`

	var chatID sql.NullInt64
	var result int64
	err := m.DB.QueryRowContext(ctx, querySelect, assistantID).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	if chatID.Valid {
		result = chatID.Int64
	} else {
		return 0, nil
	}
	return result, nil
}

func (m *postgresDBRepo) getAssistantByID(ctx context.Context, leadID int) ([]models.Assistant, error) {
	querySelect := `
	select id_assistant
	from lead_assistants
	where id_lead = $1
	`

	rows, err := m.DB.QueryContext(ctx, querySelect, leadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assistants []models.Assistant
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		assistant, err := m.GetAssistantByID(id)
		if err != nil {
			return nil, err
		}

		assistants = append(assistants, *assistant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assistants, nil

}

func (m *postgresDBRepo) updateHeroes(ctx context.Context, lead *models.Lead) error {

	heroes, err := m.getHeroesByID(ctx, lead.ID)
	if err != nil {
		return err
	}

	for i := range lead.Heroes {
		lead.Heroes[i].NeedSendMessage = true
	}

	for i := range heroes {
		heroes[i].Canceled = true
	}

	//Сравниваем артистов которые уже были занесены в таблицу с новыми артистами
	for i := range heroes {
		for y := range lead.Heroes {
			if heroes[i].ArtistID == lead.Heroes[y].ArtistID {
				heroes[i].Canceled = false
				lead.Heroes[y].NeedSendMessage = false
			}
		}
	}

	for _, heroTable := range heroes {
		if heroTable.ArtistID != 0 && heroTable.Canceled {
			chatID, err := m.getArtistChatID(ctx, heroTable.ArtistID)
			if err != nil {
				return err
			}
			if chatID != 0 {
				text := "<strong>Вас сняли с заказа №</strong> " + "<strong>" + strconv.Itoa(lead.ID) + "</strong>" + "\n" + "Дата: " + "<strong>" + string(lead.Date.Format("02-01-2006")) + "</strong>" + "\n" + "Время: " + "<strong>" + string(lead.Time.Format("15:04")) + "</strong>"

				m.App.MailChan <- modelsTelegram.MailData{
					ChatID: chatID,
					Text:   text,
				}
			}
		}
	}

	queryDelete := `
	delete 
	from lead_heroes
	where id_lead=$1
	`
	_, err = m.DB.ExecContext(ctx, queryDelete, lead.ID)
	if err != nil {
		return err
	}

	err = m.insertHeroes(ctx, &lead.Heroes, lead)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) updatePrograms(ctx context.Context, lead *models.Lead) error {
	queryDeleteProgram := `
	delete 
	from lead_programs
	where id_lead=$1
	`

	//Удаляем программы которые не были собраны
	queryDeleteOrderStore := `
	delete 
	from order_store
	where id_lead=$1 and completed<>true
	`

	_, err := m.DB.ExecContext(ctx, queryDeleteProgram, lead.ID)
	if err != nil {
		return err
	}
	_, err = m.DB.ExecContext(ctx, queryDeleteOrderStore, lead.ID)
	if err != nil {
		return err
	}

	//выборка программ которые уже были собраны и произведен вычет расходных материалов
	queryCompletedOrderStore := `
	select id_check_list
	from order_store
	where id_lead=$1 and completed=true and disassemble_bag<>true
	`
	rows, err := m.DB.QueryContext(ctx, queryCompletedOrderStore, lead.ID)
	if err != nil {
		return err
	}

	var idsPrograms []int
	for rows.Next() {
		var idProgram int
		err := rows.Scan(&idProgram)
		if err != nil {
			return err
		}
		idsPrograms = append(idsPrograms, idProgram)
	}

	//выборка программ которые уже были собраны и произведен вычет расходных материалов
	queryChangeDescription := `
	update lead_programs
	set description=$1
	where id_lead=$2 and id_check_list=$3
	`

	//Выборка всех id которые есть в изменнем лиде(новом)
	for index, item := range lead.MasterClasses {
		for y, programDB := range idsPrograms {
			if item.ID == programDB {
				_, err = m.DB.ExecContext(ctx, queryChangeDescription, item.Description, lead.ID, item.ID)
				if err != nil {
					return err
				}
				lead.MasterClasses[index] = models.MasterClass{}
				idsPrograms[y] = 0
			}
		}
	}
	for index, item := range lead.Shows {
		for y, programDB := range idsPrograms {
			if item.ID == programDB {
				_, err = m.DB.ExecContext(ctx, queryChangeDescription, item.Description, lead.ID, item.ID)
				if err != nil {
					return err
				}
				lead.Shows[index] = models.Show{}
				idsPrograms[y] = 0
			}
		}
	}
	for index, item := range lead.Others {
		for y, programDB := range idsPrograms {
			if item.ID == programDB {
				_, err = m.DB.ExecContext(ctx, queryChangeDescription, item.Description, lead.ID, item.ID)
				if err != nil {
					return err
				}
				lead.Others[index] = models.Other{}
				idsPrograms[y] = 0
			}
		}
	}
	for index, item := range lead.Animations {
		for y, programDB := range idsPrograms {
			if item.ID == programDB {
				_, err = m.DB.ExecContext(ctx, queryChangeDescription, item.Description, lead.ID, item.ID)
				if err != nil {
					return err
				}
				lead.Animations[index] = models.Animation{}
				idsPrograms[y] = 0
			}
		}
	}
	for index, item := range lead.PartyAndQuests {
		for y, programDB := range idsPrograms {
			if item.ID == programDB {
				_, err = m.DB.ExecContext(ctx, queryChangeDescription, item.Description, lead.ID, item.ID)
				if err != nil {
					return err
				}
				lead.PartyAndQuests[index] = models.PartyAndQuest{}
				idsPrograms[y] = 0
			}
		}
	}

	//Сравниваем id из БД c id которые пришли в измененном лиде
	// for _, checkListID := range idCheckLists {

	// }
	//Программы которые были уже собраны и отменены в результате редактирвоания подлежат разбору
	queryDisassembleBagOrderStore := `
	update order_store
	set disassemble_bag=true, canceled=true
	where id_lead=$1 and completed=true and id_check_list=$2
	`

	for _, programID := range idsPrograms {
		if programID != 0 {
			_, err = m.DB.ExecContext(ctx, queryDisassembleBagOrderStore, lead.ID, programID)
			if err != nil {
				return err
			}
		}
	}

	err = m.insertPrograms(ctx, lead, lead.ID)
	if err != nil {
		return err
	}
	return nil
}

// Update lead end

//Get count leads start

func (m *postgresDBRepo) GetCountRawLeads() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryRawLeads := `
	select count(id_lead)
	from leads
	where (date + time)>(current_time+current_date)  AND (check_artists<>true OR confirmed<>true OR check_assistants<>true)
	`
	var count int
	err := m.DB.QueryRowContext(ctx, queryRawLeads).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetCountConfirmedLeads() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryRawLeads := `
	select count(id_lead)
	from leads
	where (date + time)>(current_time+ current_date) AND (check_artists=true AND confirmed=true AND check_assistants=true)
	`
	var count int
	err := m.DB.QueryRowContext(ctx, queryRawLeads).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetCountArchiveLeads() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryRawLeads := `
	select count(id_lead)
	from leads
	where (date + time)<(current_time+ current_date)
	`
	var count int
	err := m.DB.QueryRowContext(ctx, queryRawLeads).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

//Get ount leads end
