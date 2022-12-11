package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetAllHeroes() (*[]models.Hero, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectAll := `
	select  id_hero, name_hero, id_gender_hero, clothing_size_min, clothing_size_max, photo
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
			&hero.Photo)
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
