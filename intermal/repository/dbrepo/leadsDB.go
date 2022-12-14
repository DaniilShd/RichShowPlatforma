package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	fmt.Println(idClient)
	idChild, err := m.insertChild(ctx, &lead.Child, idClient)
	if err != nil {
		return err
	}
	fmt.Println(idChild)

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

	err = m.insertPrograms(ctx, lead, idLead)
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
	for _, program := range programs.MasterClasses {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
		if err != nil {
			return err
		}
	}
	for _, program := range programs.Shows {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
		if err != nil {
			return err
		}
	}
	for _, program := range programs.PartyAndQuests {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
		if err != nil {
			return err
		}
	}
	for _, program := range programs.Animations {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
		if err != nil {
			return err
		}
	}
	for _, program := range programs.Others {
		_, err := m.DB.ExecContext(ctx, queryInsertProgram, program.ID, program.Description, idLead)
		if err != nil {
			return err
		}
	}
	return nil
}

// Пакет функций для добавления рограмм в заказ (Лид) End-------------------------------------------------------------------------

func (m *postgresDBRepo) insertHeroes(ctx context.Context, heroes *[]models.LeadHero, idLead int) error {

	queryInsertHeroesArtist0 := `insert into lead_heroes
	(id_hero, id_lead)
	VALUES ($1, $2)
	`

	queryInsertHeroes := `insert into lead_heroes
	(id_hero, id_lead, id_artist)
	VALUES ($1, $2, $3)
	`

	for _, hero := range *heroes {
		if hero.ArtistID == 0 {
			_, err := m.DB.ExecContext(ctx, queryInsertHeroesArtist0, hero.HeroID, idLead)
			if err != nil {
				return err
			}
		} else {
			_, err := m.DB.ExecContext(ctx, queryInsertHeroes, hero.HeroID, idLead, hero.ArtistID)
			if err != nil {
				return err
			}
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
	select id_artist
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

	for rows.Next() {
		var idArtist sql.NullInt64
		err := rows.Scan(&idArtist)
		if err != nil {
			return false, err
		}
		if int(idArtist.Int64) == 0 {
			return false, nil
		}
		idArtists = append(idArtists, int(idArtist.Int64))
	}
	if len(idArtists) == 0 {
		return false, nil
	}
	return true, nil
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
	if len(idAssistants) == 0 {
		return false, nil
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	return true, nil
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
		err := rows.Scan(&hero.ID, &artistID)
		if err != nil {
			return nil, err
		}

		if artistID.Valid {
			artist.ID = int(artistID.Int64)

			queryArtist := `
			select first_name, last_name
			from artists
			where id_artist=$1
			`
			err := m.DB.QueryRowContext(ctx, queryArtist, artist.ID).Scan(
				&artist.FirstName,
				&artist.LastName)
			if err != nil {
				return nil, err
			}
			heroLead.ArtistFirstName = artist.FirstName
			heroLead.ArtistID = artist.ID
			heroLead.ArtistLastName = artist.LastName
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
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!Здесь остановился!!!!!!!!!!
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

	_, err := m.DB.ExecContext(ctx, queryDeleteLead, idLead)
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

// Update lead start
func (m *postgresDBRepo) UpdateLead(lead *models.Lead) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	fmt.Println(1223)

	err := m.updateClient(ctx, &lead.Client)
	if err != nil {
		return err
	}

	fmt.Println(1223)
	err = m.updateChild(ctx, &lead.Child)
	if err != nil {
		return err
	}

	queryUpdateLead := `update leads
	set amount_of_children=$1, average_age_of_children=$2, address=$3, date=$4, time=$5, description=$6, duration=$7
	where id_lead=$8
	`
	fmt.Println(1223)
	var idLead int
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

	fmt.Println(1223)

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
	queryDelete := `
	delete 
	from lead_assistants
	where id_lead=$1
	`
	_, err := m.DB.ExecContext(ctx, queryDelete, lead.ID)
	if err != nil {
		return err
	}

	err = m.insertAssistants(ctx, &lead.Assistants, lead.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) updateHeroes(ctx context.Context, lead *models.Lead) error {
	queryDelete := `
	delete 
	from lead_heroes
	where id_lead=$1
	`
	_, err := m.DB.ExecContext(ctx, queryDelete, lead.ID)
	if err != nil {
		return err
	}

	err = m.insertHeroes(ctx, &lead.Heroes, lead.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) updatePrograms(ctx context.Context, lead *models.Lead) error {
	queryDelete := `
	delete 
	from lead_programs
	where id_lead=$1
	`
	_, err := m.DB.ExecContext(ctx, queryDelete, lead.ID)
	if err != nil {
		return err
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
