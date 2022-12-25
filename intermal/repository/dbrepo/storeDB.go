package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

// Store start
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

//Store end

//StoreOrder start------------------------------------------------------------------------------------------------------------------------------

func (m *postgresDBRepo) GetAllNewStoreOrder() ([]models.StoreLead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectJoinNewStoreOrder := `
	Select l.id_lead,l.date, l.time, l.description, 
	os.id_order_store, 
	cl.name_check_list, cl.id_check_list, 
	clt.name_type
	from leads l 
	RIGHT JOIN order_store os
	ON l.id_lead = os.id_lead
	LEFT JOIN check_lists cl
	ON os.id_check_list = cl.id_check_list 
	LEFT JOIN check_list_type clt
	ON cl.id_type_of_list = clt.id_type_of_list
	where disassemble_bag<>true and completed<>true
	`

	rows, err := m.DB.QueryContext(ctx, querySelectJoinNewStoreOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var StoreOrders []models.StoreLead
	var tmeString string
	for rows.Next() {
		var StoreOrder models.StoreLead
		err := rows.Scan(&StoreOrder.LeadID,
			&StoreOrder.Date,
			&tmeString,
			&StoreOrder.LeadDescription,
			&StoreOrder.ID,
			&StoreOrder.Name,
			&StoreOrder.CheckListID,
			&StoreOrder.ProgramType)

		if err != nil {
			return nil, err
		}

		StoreOrder.Time, err = time.Parse("15:04", tmeString[:5])
		if err != nil {
			return nil, err
		}

		StoreOrders = append(StoreOrders, StoreOrder)
	}
	return StoreOrders, nil
}

func (m *postgresDBRepo) GetAllCompleteStoreOrder() ([]models.StoreLead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectJoinNewStoreOrder := `
	Select l.id_lead,l.date, l.time, l.description, 
	os.id_order_store, 
	cl.name_check_list, cl.id_check_list, 
	clt.name_type
	from leads l 
	RIGHT JOIN order_store os
	ON l.id_lead = os.id_lead
	LEFT JOIN check_lists cl
	ON os.id_check_list = cl.id_check_list 
	LEFT JOIN check_list_type clt
	ON cl.id_type_of_list = clt.id_type_of_list
	where disassemble_bag<>true and completed=true
	`

	rows, err := m.DB.QueryContext(ctx, querySelectJoinNewStoreOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var StoreOrders []models.StoreLead
	var tmeString string
	for rows.Next() {
		var StoreOrder models.StoreLead
		err := rows.Scan(&StoreOrder.LeadID,
			&StoreOrder.Date,
			&tmeString,
			&StoreOrder.LeadDescription,
			&StoreOrder.ID,
			&StoreOrder.Name,
			&StoreOrder.CheckListID,
			&StoreOrder.ProgramType)

		if err != nil {
			return nil, err
		}
		StoreOrder.Time, err = time.Parse("15:04", tmeString[:5])
		if err != nil {
			return nil, err
		}
		StoreOrders = append(StoreOrders, StoreOrder)
	}
	return StoreOrders, nil
}

func (m *postgresDBRepo) GetAllToDestroyStoreOrder() ([]models.StoreLead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectJoinNewStoreOrder := `
	Select l.id_lead,l.date, l.time, l.description, 
	os.id_order_store, 
	cl.name_check_list, cl.id_check_list, 
	clt.name_type
	from leads l 
	RIGHT JOIN order_store os
	ON l.id_lead = os.id_lead
	LEFT JOIN check_lists cl
	ON os.id_check_list = cl.id_check_list 
	LEFT JOIN check_list_type clt
	ON cl.id_type_of_list = clt.id_type_of_list
	where disassemble_bag=true or (l.date+l.time)<(current_date + current_time)
	`

	rows, err := m.DB.QueryContext(ctx, querySelectJoinNewStoreOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var StoreOrders []models.StoreLead
	var tmeString string
	for rows.Next() {
		var StoreOrder models.StoreLead
		err := rows.Scan(&StoreOrder.LeadID,
			&StoreOrder.Date,
			&tmeString,
			&StoreOrder.LeadDescription,
			&StoreOrder.ID,
			&StoreOrder.Name,
			&StoreOrder.CheckListID,
			&StoreOrder.ProgramType)

		if err != nil {
			return nil, err
		}
		StoreOrder.Time, err = time.Parse("15:04", tmeString[:5])
		if err != nil {
			return nil, err
		}
		StoreOrders = append(StoreOrders, StoreOrder)
	}
	return StoreOrders, nil
}

func (m *postgresDBRepo) GetStoreOrderByID(id int) (*models.StoreLead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelectJoinNewStoreOrder := `
	Select l.id_lead, l.date, l.time, l.description, l.amount_of_children,
	os.id_order_store, os.photo, os.description, os.canceled,
	cl.name_check_list, cl.id_check_list, 
	clt.name_type
	from leads l 
	RIGHT JOIN order_store os
	ON l.id_lead = os.id_lead
	LEFT JOIN check_lists cl
	ON os.id_check_list = cl.id_check_list 
	LEFT JOIN check_list_type clt
	ON cl.id_type_of_list = clt.id_type_of_list
	where os.id_order_store = $1
	`

	var StoreOrder models.StoreLead
	var photo sql.NullString
	var description sql.NullString
	var timeString string

	err := m.DB.QueryRowContext(ctx, querySelectJoinNewStoreOrder, id).Scan(&StoreOrder.LeadID,
		&StoreOrder.Date,
		&timeString,
		&StoreOrder.LeadDescription,
		&StoreOrder.AmountOfChilds,
		&StoreOrder.ID,
		&photo,
		&description,
		&StoreOrder.Canceled,
		&StoreOrder.Name,
		&StoreOrder.CheckListID,
		&StoreOrder.ProgramType)
	if err != nil {
		return nil, err
	}
	StoreOrder.Time, err = time.Parse("15:04", timeString[:5])
	if err != nil {
		return nil, err
	}

	if photo.Valid {
		StoreOrder.Photo = photo.String
	}
	if description.Valid {
		StoreOrder.StoreDescription = description.String
	}
	return &StoreOrder, nil
}

func (m *postgresDBRepo) SetCompleteStoreOrder(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsert := `
		update order_store
		set completed=$1
		where id_order_store=$2
		`

	_, err := m.DB.ExecContext(ctx, queryInsert, true, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) UnSetCompleteStoreOrder(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsert := `
		update order_store
		set completed=$1
		where id_order_store=$2
		`

	_, err := m.DB.ExecContext(ctx, queryInsert, false, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteStoreOrderByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	storeOrder, err := m.GetStoreOrderByID(id)
	if err != nil {
		return err
	}

	if storeOrder.Canceled {
		m.addItemToStore(ctx, storeOrder.CheckListID, storeOrder.AmountOfChilds)
	}

	querySelect := `
	select photo 
	from order_store
	where id_order_store = $1
	`
	var photo string
	row := m.DB.QueryRowContext(ctx, querySelect, id)
	row.Scan(&photo)

	if photo != "" {
		filePath := "static/img/store-leads/" + photo
		fmt.Println(filePath)
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	queryDelete := `
		Delete 
		from order_store
		where id_order_store=$1
		`

	_, err = m.DB.ExecContext(ctx, queryDelete, id)
	if err != nil {
		return err
	}
	return nil

}

func (m *postgresDBRepo) InsertStoreOrder(storeOrder *models.StoreLead) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryInsertTest := `
	update order_store
		set photo=$1, description=$2
		where id_order_store=$3
	`
	_, err := m.DB.ExecContext(ctx, queryInsertTest, storeOrder.Photo, storeOrder.StoreDescription, storeOrder.ID)
	if err != nil {
		return err
	}

	m.SetCompleteStoreOrder(storeOrder.ID)
	m.removeItemFromStore(ctx, storeOrder.CheckListID, storeOrder.AmountOfChilds)

	fmt.Println("complete")
	fmt.Println(storeOrder.ID)

	return nil
}

// Функция по удаленею расходников из склада
func (m *postgresDBRepo) removeItemFromStore(ctx context.Context, idCheckList int, countChild int) error {
	checkList, err := m.GetCheckListByID(idCheckList)
	if err != nil {
		return err
	}

	var amount float64
	for _, item := range checkList.Items {
		amount = item.AmountItemOnce * float64(countChild)
		m.updateAmountStoreItemMinus(ctx, item.ID, amount)
	}

	return nil
}

func (m *postgresDBRepo) addItemToStore(ctx context.Context, idCheckList int, countChild int) error {
	checkList, err := m.GetCheckListByID(idCheckList)
	if err != nil {
		return err
	}

	var amount float64
	for _, item := range checkList.Items {
		amount = item.AmountItemOnce * float64(countChild)
		m.updateAmountStoreItemPlus(ctx, item.ID, amount)
	}

	return nil
}

func (m *postgresDBRepo) updateAmountStoreItemMinus(ctx context.Context, idItemStore int, amount float64) error {
	queryUpdateMinus := `
	update store
	set current_amount=current_amount-$1
	where id_item=$2
	`

	fmt.Println("idItemStore ", idItemStore)

	_, err := m.DB.ExecContext(ctx, queryUpdateMinus, amount, idItemStore)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) updateAmountStoreItemPlus(ctx context.Context, idItemStore int, amount float64) error {
	queryUpdatePlus := `
	update store
	set current_amount=current_amount+$1
	where id_item=$2
	`
	_, err := m.DB.ExecContext(ctx, queryUpdatePlus, amount, idItemStore)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) UpdateStoreOrder(storeOrder *models.StoreLead) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryUpdate := `
	update order_store
	set photo=$1, description=$2
	where id_order_store=$3
	`
	_, err := m.DB.ExecContext(ctx, queryUpdate, storeOrder.Photo, storeOrder.StoreDescription, storeOrder.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetCountNewStoreOrder() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `
	select count(id_order_store)
	from order_store
	where disassemble_bag<>true and completed<>true
	`
	var count int
	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetCountCompleteStoreOrder() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `
	select count(id_order_store)
	from order_store
	where disassemble_bag<>true and completed=true
	`
	var count int
	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetCountToDestroyStoreOrder() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	//Назначает строки на разбор если они были собраны и праздник прошел (по времени)
	queryUpdate := `
	update order_store
set disassemble_bag=true
from 
(select id_lead 
from leads where (date + time) < (current_date + current_time))  as testTable
where (testTable.id_lead = order_store.id_lead and order_store.completed=true)
	`

	//Удаляет те строчки, которые не были собраны и время праздника прошло
	queryToDelete := `
	delete
from order_store
where id_order_store in (select id_order_store from order_store os 
left join leads l 
on l.id_lead = os.id_lead 
where (l.date + l.time) < (current_date + current_time) and os.completed=false)
		`

	_, err := m.DB.ExecContext(ctx, queryUpdate)
	if err != nil {
		return 0, err
	}

	_, err = m.DB.ExecContext(ctx, queryToDelete)
	if err != nil {
		return 0, err
	}

	query := `
	select count(id_order_store)
	from order_store os
	where (os.disassemble_bag=true or (select (date + time) from leads l where os.id_lead = l.id_lead)<(current_date + current_time))
	`
	var count int
	err = m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

//StoreOrder end
