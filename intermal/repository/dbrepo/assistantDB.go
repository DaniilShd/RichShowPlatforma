package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetAllAssistants() (*[]models.Assistant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAll := `
	select  id_assistant, first_name, last_name, id_gender_type, telegram_assistant, phone_number, photo_assistant
	from assistants
	`

	var assistants []models.Assistant

	rows, err := m.DB.QueryContext(ctx, querySelectAll)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var assistant models.Assistant
		err = rows.Scan(&assistant.ID,
			&assistant.FirstName,
			&assistant.LastName,
			&assistant.Gender,
			&assistant.Telegram,
			&assistant.PhoneNumber,
			&assistant.Photo)
		if err != nil {
			return nil, err
		}
		assistants = append(assistants, assistant)
	}
	fmt.Println(assistants)

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &assistants, nil
}
