package dbrepo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetAllArtists() (*[]models.Artist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAll := `
	select  id_artist, first_name, last_name, id_gender_type, growth, shoe_size, clothing_size, telegram_artist, phone_number, photo_artist, description, vk
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
			&artist.Photo,
			&artist.Description,
			&artist.VK)
		if err != nil {
			return nil, err
		}
		artist.PhoneNumber = helpers.ConvertNumberPhone(artist.PhoneNumber)

		artists = append(artists, artist)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &artists, nil
}

func (m *postgresDBRepo) GetArtistByID(id int) (*models.Artist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select  id_artist, first_name, last_name, id_gender_type, growth, shoe_size, clothing_size, telegram_artist, phone_number, photo_artist, description, vk
	from artists
	where id_artist = $1
	`

	var artist models.Artist

	row := m.DB.QueryRowContext(ctx, querySelect, id)

	err := row.Scan(&artist.ID,
		&artist.FirstName,
		&artist.LastName,
		&artist.Gender,
		&artist.Growth,
		&artist.ShoeSize,
		&artist.ClothingSize,
		&artist.Telegram,
		&artist.PhoneNumber,
		&artist.Photo,
		&artist.Description,
		&artist.VK,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(artist)

	if err = row.Err(); err != nil {
		return nil, err
	}

	return &artist, nil
}

func (m *postgresDBRepo) InsertArtist(artist *models.Artist) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsert := `
	insert into artists (first_name, last_name, id_gender_type, growth, shoe_size, clothing_size, telegram_artist, phone_number, photo_artist, description, vk)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := m.DB.ExecContext(ctx, queryInsert,
		artist.FirstName,
		artist.LastName,
		artist.Gender,
		artist.Growth,
		artist.ShoeSize,
		artist.ClothingSize,
		artist.Telegram,
		artist.PhoneNumber,
		artist.Photo,
		artist.Description,
		artist.VK,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) UpdateArtist(artist *models.Artist) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select photo_artist 
	from artists
	where id_artist = $1
	`
	var photo string
	row := m.DB.QueryRowContext(ctx, querySelect, artist.ID)
	row.Scan(&photo)

	if artist.Photo != "" {

		filePath := "static/img/animators/" + photo
		err := os.Remove(filePath)
		if err != nil {
			return err
		}

	} else {
		artist.Photo = photo
	}

	queryUpdate := `update artists 
		set first_name=$1, last_name=$2, id_gender_type=$3, growth=$4, shoe_size=$5, clothing_size=$6, telegram_artist=$7, phone_number=$8, photo_artist=$9, description=$10, vk=$11
		where id_artist = $12
		`

	_, err := m.DB.ExecContext(ctx, queryUpdate,
		artist.FirstName,
		artist.LastName,
		artist.Gender,
		artist.Growth,
		artist.ShoeSize,
		artist.ClothingSize,
		artist.Telegram,
		artist.PhoneNumber,
		artist.Photo,
		artist.Description,
		artist.VK,
		artist.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) DeleteArtistByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select photo_artist 
	from artists
	where id_artist = $1
	`
	var photo string
	row := m.DB.QueryRowContext(ctx, querySelect, id)
	row.Scan(&photo)

	if photo != "" {
		filePath := "static/img/animators/" + photo
		fmt.Println(filePath)
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	deleteItem := `
	delete
	from artists
	where id_artist = $1
	`

	_, err := m.DB.ExecContext(ctx, deleteItem, id)
	if err != nil {
		return err
	}

	return nil
}
