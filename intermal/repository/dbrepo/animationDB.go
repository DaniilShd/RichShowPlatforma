package dbrepo

import (
	"context"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
)

func (m *postgresDBRepo) GetAnimationByID(id int) (*models.Animation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var animation models.Animation

	queryanimation := `
	select id_animation, name_animation, description, id_check_list 
	from animation
	where id_animation = $1
	`

	rows := m.DB.QueryRowContext(ctx, queryanimation, id)

	var idChekList int
	err := rows.Scan(&animation.ID, &animation.Name, &animation.Description, &idChekList)
	if err != nil {
		return nil, err
	}

	var checkList *models.CheckList
	checkList, err = m.GetCheckListByID(idChekList)
	animation.CheckList = *checkList

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &animation, nil
}

func (m *postgresDBRepo) GetAllAnimation() (*[]models.Animation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var animations []models.Animation

	queryanimation := `
	select id_animation, name_animation, description, id_check_list 
	from animation
	`

	rows, err := m.DB.QueryContext(ctx, queryanimation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var animation models.Animation
		var idChekList int
		err := rows.Scan(&animation.ID, &animation.Name, &animation.Description, &idChekList)
		if err != nil {
			return nil, err
		}

		var checkList *models.CheckList
		checkList, err = m.GetCheckListByID(idChekList)
		animation.CheckList = *checkList

		animations = append(animations, animation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &animations, nil
}

// Update animation in database
func (m *postgresDBRepo) UpdateAnimation(animation *models.Animation) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `update animation set name_animation = $1, description = $2, id_check_list = $3
	where id_animation = $4
	`

	_, err := m.DB.ExecContext(ctx, query,
		animation.Name,
		animation.Description,
		animation.CheckList.ID,
		animation.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// delete master class by id
func (m *postgresDBRepo) DeleteAnimationByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `delete from animation
	where id_animation = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil

} //Add new master-class
func (m *postgresDBRepo) InsertAnimation(animation *models.Animation) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `insert into animation (name_animation, description, id_check_list)
	VALUES ($1, $2, $3)
	`

	_, err := m.DB.ExecContext(ctx, query,
		animation.Name,
		animation.Description,
		animation.CheckList.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
