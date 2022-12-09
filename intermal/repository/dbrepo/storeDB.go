package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetAllStoreItem() ([]models.StoreItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var storeItems []models.StoreItem

	queryStoreItems := `
	select id_item, name_item, current_amount, dimension, min_amount, description
	from store
	`

	rows, err := m.DB.QueryContext(ctx, queryStoreItems)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var storeItem models.StoreItem
		err := rows.Scan(&storeItem.ID, &storeItem.Name, &storeItem.CurrentAmount, &storeItem.Dimension, &storeItem.MinAmount, &storeItem.Description)
		if err != nil {
			return nil, err
		}

		storeItems = append(storeItems, storeItem)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return storeItems, nil
}

func (m *postgresDBRepo) GetStoreItemByID(id int) (*models.StoreItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var storeItem models.StoreItem

	queryStoreItems := `
	select id_item, name_item, current_amount, dimension, min_amount, description
	from store
	where id_item = $1
	`

	row := m.DB.QueryRowContext(ctx, queryStoreItems, id)

	err := row.Scan(&storeItem.ID, &storeItem.Name, &storeItem.CurrentAmount, &storeItem.Dimension, &storeItem.MinAmount, &storeItem.Description)
	if err != nil {
		return nil, err

	}
	if err = row.Err(); err != nil {
		return nil, err
	}

	return &storeItem, nil
}

func (m *postgresDBRepo) DeleteStoreItemByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	deleteItem := `
	delete
	from store
	where id_item = $1
	`

	_, err := m.DB.ExecContext(ctx, deleteItem, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) InsertStoreItem(storeItem *models.StoreItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	InsertStoreItem := `
	insert into store (name_item, current_amount, dimension, min_amount, description)
	values ($1, $2, $3, $4, $5)
	`

	_, err := m.DB.ExecContext(ctx, InsertStoreItem,
		&storeItem.Name,
		&storeItem.CurrentAmount,
		&storeItem.Dimension,
		&storeItem.MinAmount,
		&storeItem.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdateStoreItem(storeItem *models.StoreItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	InsertStoreItem := `
	update store 
	set name_item=$1, current_amount=$2, dimension=$3, min_amount=$4, description=$5
	where id_item = $6
	`

	_, err := m.DB.ExecContext(ctx, InsertStoreItem,
		&storeItem.Name,
		&storeItem.CurrentAmount,
		&storeItem.Dimension,
		&storeItem.MinAmount,
		&storeItem.Description,
		&storeItem.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
