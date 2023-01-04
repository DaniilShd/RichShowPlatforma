package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
)

func (m *postgresDBRepoTelegram) GetAllLeadsOfArtistByChatID(chatID int64) ([]models.LeadList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var leads []models.LeadList

	querySelect := `
	select distinct l.date, l.time, l.address, l.id_lead, a.id_artist
from lead_heroes lh
left join artists a
on lh.id_artist = a.id_artist 
left join leads l 
on lh.id_lead = l.id_lead 
left join lead_heroes le
on le.id_lead = l.id_lead 
left join heroes h
on le.id_hero = h.id_hero 
where a.id_telegram_chat=$1 and (l.date + l.time) > (current_date + current_time)
	`

	querySelectHeroes := `
	select name_hero
	from lead_heroes lh
	left join heroes h
	on lh.id_hero = h.id_hero 
	where lh.id_lead = $1 and lh.id_artist=$2
	`

	var timeDate string
	rows, err := m.DB.QueryContext(ctx, querySelect, chatID)
	if err != nil {
		return nil, err
	}
	var lead models.LeadList
	for rows.Next() {
		err := rows.Scan(&lead.Date, &timeDate, &lead.Address, &lead.ID, &lead.ArtistID)
		if err != nil {
			return nil, err
		}

		for _, item := range leads {
			if lead.ID == item.ID {
				continue
			}
		}

		lead.Time, err = time.Parse("15:04", timeDate[:5])
		if err != nil {
			return nil, err
		}
		rowHero, err := m.DB.QueryContext(ctx, querySelectHeroes, lead.ID, lead.ArtistID)
		if err != nil {
			return nil, err
		}
		for rowHero.Next() {
			var hero string
			err := rowHero.Scan(&hero)
			if err != nil {
				return nil, err
			}
			lead.NameHeroes = append(lead.NameHeroes, hero)
		}
		leads = append(leads, lead)
	}

	return leads, nil
}

func (m *postgresDBRepoTelegram) GetAllLeadsOfArtistTodayByChatID(chatID int64) ([]models.LeadList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var leads []models.LeadList

	querySelect := `
	select distinct l.date, l.time, l.address, l.id_lead, a.id_artist
	from lead_heroes lh
	left join artists a
	on lh.id_artist = a.id_artist 
	left join leads l 
	on lh.id_lead = l.id_lead  
	left join lead_heroes le
	on le.id_lead = l.id_lead 
	left join heroes h
	on le.id_hero = h.id_hero 
	where a.id_telegram_chat=$1 and l.date=current_date
	`

	querySelectHeroes := `
	select name_hero
	from lead_heroes lh
	left join heroes h
	on lh.id_hero = h.id_hero 
	where lh.id_lead = $1 and lh.id_artist=$2
	`

	var timeDate string
	rows, err := m.DB.QueryContext(ctx, querySelect, chatID)
	if err != nil {
		return nil, err
	}
	var lead models.LeadList
	for rows.Next() {
		err := rows.Scan(&lead.Date, &timeDate, &lead.Address, &lead.ID, &lead.ArtistID)
		if err != nil {
			return nil, err
		}

		for _, item := range leads {
			if lead.ID == item.ID {
				continue
			}
		}

		lead.Time, err = time.Parse("15:04", timeDate[:5])
		if err != nil {
			return nil, err
		}
		rowHero, err := m.DB.QueryContext(ctx, querySelectHeroes, lead.ID, lead.ArtistID)
		if err != nil {
			return nil, err
		}
		for rowHero.Next() {
			var hero string
			err := rowHero.Scan(&hero)
			if err != nil {
				return nil, err
			}
			lead.NameHeroes = append(lead.NameHeroes, hero)
		}
		leads = append(leads, lead)
	}

	return leads, nil
}
