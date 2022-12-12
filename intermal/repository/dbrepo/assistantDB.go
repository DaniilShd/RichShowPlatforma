package dbrepo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
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

		assistant.PhoneNumber = helpers.ConvertNumberPhone(assistant.PhoneNumber)

		assistants = append(assistants, assistant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &assistants, nil
}

func (m *postgresDBRepo) GetAssistantByID(id int) (*models.Assistant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select  id_assistant, first_name, last_name, id_gender_type, telegram_assistant, phone_number, photo_assistant,
description, vk
	from assistants
	where id_assistant = $1
	`

	var assistant models.Assistant

	row := m.DB.QueryRowContext(ctx, querySelect, id)

	err := row.Scan(&assistant.ID,
		&assistant.FirstName,
		&assistant.LastName,
		&assistant.Gender,
		&assistant.Telegram,
		&assistant.PhoneNumber,
		&assistant.Photo,
		&assistant.Description,
		&assistant.VK,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(assistant)

	if err = row.Err(); err != nil {
		return nil, err
	}

	return &assistant, nil
}

func (m *postgresDBRepo) InsertAssistant(assistant *models.Assistant) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsert := `
	insert into assistants 
	(first_name, last_name, id_gender_type, telegram_assistant, phone_number, photo_assistant, description, vk)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := m.DB.ExecContext(ctx, queryInsert,
		&assistant.FirstName,
		&assistant.LastName,
		&assistant.Gender,
		&assistant.Telegram,
		&assistant.PhoneNumber,
		&assistant.Photo,
		&assistant.Description,
		&assistant.VK,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) UpdateAssistant(assistant *models.Assistant) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select photo_assistant  
	from assistants
	where id_assistant = $1
	`
	var photo string
	row := m.DB.QueryRowContext(ctx, querySelect, assistant.ID)
	row.Scan(&photo)

	if assistant.Photo != "" {

		filePath := "static/img/assistants/" + photo
		err := os.Remove(filePath)
		if err != nil {
			return err
		}

	} else {
		assistant.Photo = photo
	}

	queryUpdate := `update assistants 
	set first_name=$1, last_name=$2, id_gender_type=$3, telegram_assistant=$4, phone_number=$5, photo_assistant=$6, description=$7, vk=$8
	where id_assistant = $9
	`

	_, err := m.DB.ExecContext(ctx, queryUpdate,
		&assistant.FirstName,
		&assistant.LastName,
		&assistant.Gender,
		&assistant.Telegram,
		&assistant.PhoneNumber,
		&assistant.Photo,
		&assistant.Description,
		&assistant.VK,
		&assistant.ID)
	fmt.Println(assistant.Gender)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteAssistantByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select photo_assistant 
	from assistants
	where id_assistant = $1
	`
	var photo string
	row := m.DB.QueryRowContext(ctx, querySelect, id)
	row.Scan(&photo)

	if photo != "" {
		filePath := "static/img/assistants/" + photo
		fmt.Println(filePath)
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	deleteItem := `
	delete
	from assistants
	where id_assistant = $1
	`

	_, err := m.DB.ExecContext(ctx, deleteItem, id)
	if err != nil {
		return err
	}

	return nil
}
