package dbrepo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetAllHeroes() (*[]models.Hero, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAll := `
	select  id_hero, name_hero, id_gender_hero, clothing_size_min, clothing_size_max, photo, description
	from heroes
	`

	var heroes []models.Hero

	rows, err := m.DB.QueryContext(ctx, querySelectAll)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var hero models.Hero
		err = rows.Scan(&hero.ID,
			&hero.Name,
			&hero.Gender,
			&hero.ClothingSizeMin,
			&hero.ClothingSizeMax,
			&hero.Photo,
			&hero.Description,
		)
		if err != nil {
			return nil, err
		}
		heroes = append(heroes, hero)
	}
	fmt.Println(heroes)

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &heroes, nil
}

func (m *postgresDBRepo) GetHeroByID(id int) (*models.Hero, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select  id_hero, name_hero, id_gender_hero, clothing_size_min, clothing_size_max, photo, description
	from heroes
	where id_hero = $1
	`

	var hero models.Hero

	rows := m.DB.QueryRowContext(ctx, querySelect, id)

	err := rows.Scan(&hero.ID,
		&hero.Name,
		&hero.Gender,
		&hero.ClothingSizeMin,
		&hero.ClothingSizeMax,
		&hero.Photo,
		&hero.Description,
	)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &hero, nil
}

func (m *postgresDBRepo) InsertHero(hero *models.Hero) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsert := `
	insert into heroes 
	(name_hero, id_gender_hero, clothing_size_min, clothing_size_max, photo, description)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := m.DB.ExecContext(ctx, queryInsert,
		&hero.Name,
		&hero.Gender,
		&hero.ClothingSizeMin,
		&hero.ClothingSizeMax,
		&hero.Photo,
		&hero.Description,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) UpdateHero(hero *models.Hero) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select photo  
	from heroes
	where id_hero = $1
	`
	var photo string
	row := m.DB.QueryRowContext(ctx, querySelect, hero.ID)
	row.Scan(&photo)

	if hero.Photo != "" {

		filePath := "static/img/heroes/" + photo
		err := os.Remove(filePath)
		if err != nil {
			return err
		}

	} else {
		hero.Photo = photo
	}

	queryUpdate := `update heroes 
	set name_hero=$1, id_gender_hero=$2, clothing_size_min=$3, clothing_size_max=$4, photo=$5, description=$6
	where id_hero = $7
	`

	_, err := m.DB.ExecContext(ctx, queryUpdate,
		&hero.Name,
		&hero.Gender,
		&hero.ClothingSizeMin,
		&hero.ClothingSizeMax,
		&hero.Photo,
		&hero.Description,
		&hero.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteHeroByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select photo  
	from heroes
	where id_hero = $1
	`

	var photo string
	row := m.DB.QueryRowContext(ctx, querySelect, id)
	row.Scan(&photo)

	if photo != "" {
		filePath := "static/img/heroes/" + photo
		fmt.Println(filePath)
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	deleteItem := `
	delete
	from heroes
	where id_hero = $1
	`

	_, err := m.DB.ExecContext(ctx, deleteItem, id)
	if err != nil {
		return err
	}

	return nil
}
