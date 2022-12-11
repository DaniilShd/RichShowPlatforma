package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetAllArtists() (*[]models.Artist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAll := `
	select  id_artist, first_name, last_name, id_gender_type, growth, shoe_size, clothing_size, telegram_artist, phone_number, photo_artist
	from artists
	`

	var artists []models.Artist

	rows, err := m.DB.QueryContext(ctx, querySelectAll)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var artist models.Artist
		err = rows.Scan(&artist.ID,
			&artist.FirstName,
			&artist.LastName,
			&artist.Gender,
			&artist.Growth,
			&artist.ShoeSize,
			&artist.ClothingSize,
			&artist.Telegram,
			&artist.PhoneNumber,
			&artist.Photo)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	fmt.Println(artists)

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &artists, nil
}
