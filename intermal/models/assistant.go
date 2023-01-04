package models

type Assistant struct {
	ID             int    `db:"id_assistant"`
	FirstName      string `db:"first_name"`
	LastName       string `db:"last_name"`
	PhoneNumber    string `db:"phone_number"`
	Telegram       string `db:"check_assistants"`
	Gender         int    `db:"gender"`
	Photo          string `db:"photo_assistant"`
	VK             string
	Description    string
	TelegramChatID int
}
