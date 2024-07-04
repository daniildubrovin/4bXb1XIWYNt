package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Day struct {
	ID       int
	Calories int
	Proteins int
	Fats     int
	Carbs    int
	Fibre    int
	Date     time.Time
}

type DayModel struct {
	Pg  *pgxpool.Pool
	Ctx context.Context
}

func (m *DayModel) Insert(calories int, proteins int, fats int, carbs int, fibre int, date string) (int, error) {
	query := "INSERT INTO days (calories, proteins, fats, carbs, fibre, date) VALUES (@calories, @proteins, @fats, @carbs, @fibre, @date) RETURNING id"
	args := pgx.NamedArgs{
		"calories": calories,
		"proteins": proteins,
		"fats":     fats,
		"carbs":    carbs,
		"fibre":    fibre,
		"date":     date,
	}

	var id int
	err := m.Pg.QueryRow(m.Ctx, query, args).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("unable to insert row: %w", err)
	}

	return id, nil
}

func (m *DayModel) Get(id int) (*Day, error) {
	query := `SELECT id, calories, proteins, fats, carbs, fibre, date FROM days WHERE id = $1`
	row := m.Pg.QueryRow(m.Ctx, query, id)

	d := &Day{}
	err := row.Scan(&d.ID, &d.Calories, &d.Proteins, &d.Fats, &d.Carbs, &d.Fibre, &d.Date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return d, nil
}

func (m *DayModel) Latest() ([]*Day, error) {
	query := `SELECT id, calories, proteins, fats, carbs, fibre, date FROM days ORDER BY id DESC LIMIT 10;`
	rows, err := m.Pg.Query(m.Ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	days := []*Day{}

	for rows.Next() {
		d := &Day{}
		err := rows.Scan(&d.ID, &d.Calories, &d.Proteins, &d.Fats, &d.Carbs, &d.Fibre, &d.Date)
		if err != nil {
			return nil, err
		}

		days = append(days, d)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return days, nil
}
