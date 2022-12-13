package models

type Hero struct {
	ID              int    `db:"id_assistant"`
	Name            string `db:"name_hero"`
	Gender          int    `db:"gender_hero"`
	Photo           string `db:"photo_assistant"`
	ClothingSizeMin int    `db:"clothing_size_min"`
	ClothingSizeMax int    `db:"clothing_size_max"`
	Description     string
}
