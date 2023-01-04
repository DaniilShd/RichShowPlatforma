package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	myconstant "github.com/DaniilShd/RichShowPlatforma/intermal/telegram/constant"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/repository"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/service/accounts"
)

type postgresDBRepoTelegram struct {
	DB *sql.DB
}

func NewPostgresRepoTelegram(conn *sql.DB) repository.DatabaseRepoTelegram {
	return &postgresDBRepoTelegram{
		DB: conn,
	}
}

func (m *postgresDBRepoTelegram) GetArtistChatID() ([]accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectArtistChatID := `
	Select id_artist, first_name, last_name, id_telegram_chat, phone_number
	from artists
	`
	var chatID sql.NullInt64
	var phone sql.NullString
	var persons []accounts.Person

	rows, err := m.DB.QueryContext(ctx, querySelectArtistChatID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var person accounts.Person
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
		if err != nil {
			return nil, err
		}
		if chatID.Valid {
			person.ChatID = chatID.Int64
		} else {
			person.ChatID = 0
		}
		if phone.Valid {
			person.PhoneNumber = phone.String
		} else {
			person.PhoneNumber = ""
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (m *postgresDBRepoTelegram) GetStoreChatID() ([]accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectStoreChatID := `
	Select id_account, first_name, last_name, id_telegram_chat, phone_number
	from id_accounts 
	where access_level = 3
	`
	var chatID sql.NullInt64
	var phone sql.NullString
	var persons []accounts.Person

	rows, err := m.DB.QueryContext(ctx, querySelectStoreChatID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var person accounts.Person
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
		if err != nil {
			return nil, err
		}
		if chatID.Valid {
			person.ChatID = chatID.Int64
		} else {
			person.ChatID = 0
		}
		if phone.Valid {
			person.PhoneNumber = phone.String
		} else {
			person.PhoneNumber = ""
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (m *postgresDBRepoTelegram) GetAssistantsChatID() ([]accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAssistantsChatID := `
	Select id_assistant, first_name, last_name, id_telegram_chat, phone_number
	from assistants
	`
	var chatID sql.NullInt64
	var phone sql.NullString
	var persons []accounts.Person

	rows, err := m.DB.QueryContext(ctx, querySelectAssistantsChatID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var person accounts.Person
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
		if err != nil {
			return nil, err
		}
		if chatID.Valid {
			person.ChatID = chatID.Int64
		} else {
			person.ChatID = 0
		}
		if phone.Valid {
			person.PhoneNumber = phone.String
		} else {
			person.PhoneNumber = ""
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (m *postgresDBRepoTelegram) GetManagerChatID() ([]accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectManagerChatID := `
	Select id_account, first_name, last_name, id_telegram_chat, phone_number
	from id_accounts 
	where access_level = 2
	`
	var chatID sql.NullInt64
	var phone sql.NullString
	var persons []accounts.Person

	rows, err := m.DB.QueryContext(ctx, querySelectManagerChatID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var person accounts.Person
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
		if err != nil {
			return nil, err
		}
		if chatID.Valid {
			person.ChatID = chatID.Int64
		} else {
			person.ChatID = 0
		}
		if phone.Valid {
			person.PhoneNumber = phone.String
		} else {
			person.PhoneNumber = ""
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (m *postgresDBRepoTelegram) GetAdminChatID() ([]accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAdminChatID := `
	Select id_account, first_name, last_name, id_telegram_chat, phone_number
	from id_accounts 
	where access_level = 1
	`
	var chatID sql.NullInt64
	var phone sql.NullString
	var persons []accounts.Person

	rows, err := m.DB.QueryContext(ctx, querySelectAdminChatID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var person accounts.Person
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
		if err != nil {
			return nil, err
		}
		if chatID.Valid {
			person.ChatID = chatID.Int64
		} else {
			person.ChatID = 0
		}
		if phone.Valid {
			person.PhoneNumber = phone.String
		} else {
			person.PhoneNumber = ""
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (m *postgresDBRepoTelegram) GetArtistByChatID(id int64) (*accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectArtistChatID := `
	Select id_artist, first_name, last_name, id_telegram_chat, phone_number
	from artists
	where id_telegram_chat = $1
	`
	var chatID sql.NullInt64
	var phone sql.NullString

	rows := m.DB.QueryRowContext(ctx, querySelectArtistChatID, id)

	var person accounts.Person
	err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
	if err != nil {
		return nil, err
	}
	if chatID.Valid {
		person.ChatID = chatID.Int64
	} else {
		person.ChatID = 0
	}
	if phone.Valid {
		person.PhoneNumber = phone.String
	} else {
		person.PhoneNumber = ""
	}

	return &person, nil
}

func (m *postgresDBRepoTelegram) GetStoreByChatID(id int64) (*accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectStoreChatID := `
	Select id_account, first_name, last_name, id_telegram_chat, phone_number
	from id_accounts 
	where access_level = 3 and id_telegram_chat = $1
	`
	var chatID sql.NullInt64
	var phone sql.NullString

	rows := m.DB.QueryRowContext(ctx, querySelectStoreChatID, id)

	var person accounts.Person
	err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
	if err != nil {
		return nil, err
	}
	if chatID.Valid {
		person.ChatID = chatID.Int64
	} else {
		person.ChatID = 0
	}
	if phone.Valid {
		person.PhoneNumber = phone.String
	} else {
		person.PhoneNumber = ""
	}

	return &person, nil
}

func (m *postgresDBRepoTelegram) GetAssistantsByChatID(id int64) (*accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAssistantsChatID := `
	Select id_assistant, first_name, last_name, id_telegram_chat, phone_number
	from assistants
	where id_telegram_chat = $1
	`
	var chatID sql.NullInt64
	var phone sql.NullString

	rows := m.DB.QueryRowContext(ctx, querySelectAssistantsChatID, id)

	var person accounts.Person
	err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
	if err != nil {
		return nil, err
	}
	if chatID.Valid {
		person.ChatID = chatID.Int64
	} else {
		person.ChatID = 0
	}
	if phone.Valid {
		person.PhoneNumber = phone.String
	} else {
		person.PhoneNumber = ""
	}

	return &person, nil
}

func (m *postgresDBRepoTelegram) GetManagerByChatID(id int64) (*accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectManagerChatID := `
	Select id_account, first_name, last_name, id_telegram_chat, phone_number
	from id_accounts 
	where access_level = 2 and id_telegram_chat = $1
	`
	var chatID sql.NullInt64
	var phone sql.NullString

	rows := m.DB.QueryRowContext(ctx, querySelectManagerChatID, id)

	var person accounts.Person
	err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
	if err != nil {
		return nil, err
	}
	if chatID.Valid {
		person.ChatID = chatID.Int64
	} else {
		person.ChatID = 0
	}
	if phone.Valid {
		person.PhoneNumber = phone.String
	} else {
		person.PhoneNumber = ""
	}

	return &person, nil
}

func (m *postgresDBRepoTelegram) GetAdminByChatID(id int64) (*accounts.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAdminChatID := `
	Select id_account, first_name, last_name, id_telegram_chat, phone_number
	from id_accounts 
	where access_level = 1 and id_telegram_chat = $1
	`
	var chatID sql.NullInt64
	var phone sql.NullString

	rows := m.DB.QueryRowContext(ctx, querySelectAdminChatID, id)

	var person accounts.Person
	err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &chatID, &phone)
	if err != nil {
		return nil, err
	}
	if chatID.Valid {
		person.ChatID = chatID.Int64
	} else {
		person.ChatID = 0
	}
	if phone.Valid {
		person.PhoneNumber = phone.String
	} else {
		person.PhoneNumber = ""
	}

	return &person, nil
}

func (m *postgresDBRepoTelegram) SetChatIDByRoleAndID(role int, idAccount int, chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var table string
	var column string

	if role == myconstant.ADMIN || role == myconstant.MANAGER || role == myconstant.STORE {
		table = "id_accounts"
		column = "id_account"
	}

	if role == myconstant.ARTIST {
		table = "artists"
		column = "id_artist"
	}

	if role == myconstant.ASSISTANT {
		table = "assistants"
		column = "id_assistant"
	}

	queryUpdate := fmt.Sprintf(`
	Update %s
	set id_telegram_chat=$1
	where %s=$2
	`, table, column)

	_, err := m.DB.ExecContext(ctx, queryUpdate, chatID, idAccount)
	if err != nil {
		return err
	}

	return nil
}
