package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
)

func (m *postgresDBRepoTelegram) GetAllLeadsOfAssistantByChatID(chatID int64) ([]models.LeadList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var leads []models.LeadList

	querySelect := `
	select distinct l.date, l.time, l.address, l.id_lead, a.id_assistant
from lead_assistants la
left join assistants a
on la.id_assistant = a.id_assistant 
left join leads l 
on la.id_lead = l.id_lead 
where a.id_telegram_chat=$1 and (l.date + l.time) > (current_date + current_time)
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
		leads = append(leads, lead)
	}

	return leads, nil
}

func (m *postgresDBRepoTelegram) GetAllLeadsOfAssistantTodayByChatID(chatID int64) ([]models.LeadList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var leads []models.LeadList

	querySelect := `
	select distinct l.date, l.time, l.address, l.id_lead, a.id_assistant
from lead_assistants la
left join assistants a
on la.id_assistant = a.id_assistant 
left join leads l 
on la.id_lead = l.id_lead 
	where a.id_telegram_chat=$1 and l.date=current_date
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
		leads = append(leads, lead)
	}

	return leads, nil
}
