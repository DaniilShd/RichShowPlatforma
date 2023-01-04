package models

type Artist struct {
	ID             int    `db:"id_artist"`
	FirstName      string `db:"first_name"`
	LastName       string `db:"last_name"`
	Gender         int    `db:"id_gender_type"`
	Growth         int    `db:"growth"`
	ShoeSize       int    `db:"shoe_size"`
	ClothingSize   int    `db:"clothing_size"`
	Telegram       string `db:"telegram_artist"`
	PhoneNumber    string `db:"phone_number"`
	Photo          string `db:"photo_artist"`
	Description    string
	VK             string
	TelegramChatID int
}
