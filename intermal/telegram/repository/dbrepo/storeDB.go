package dbrepo

import (
	"context"
	"time"
)

func (m *postgresDBRepoTelegram) GetOrderStoreIDByLeadID(leadID int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	querySelect := `
	select os.id_order_store 
from order_store os 
left join leads l 
on os.id_lead = l.id_lead 
where l.id_lead = $1
	`

	var idStoreOrder []int

	rows, err := m.DB.QueryContext(ctx, querySelect, leadID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		idStoreOrder = append(idStoreOrder, id)
	}
	return idStoreOrder, nil
}
