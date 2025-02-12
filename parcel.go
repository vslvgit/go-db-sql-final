package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)", sql.Named("number", p.Number), sql.Named("client", p.Client), sql.Named("status", p.Status), sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))

	if err != nil {

		return 0, err

	}

	id, err := res.LastInsertId()

	if err != nil {

		return 0, err

	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	row := s.db.QueryRow("SELECT number, client, status,address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {

		return Parcel{}, err

	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client),
	)
	if err != nil {
		return nil, err
	}

	var parcels []Parcel

	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number", sql.Named("status", status), sql.Named("number", number))

	if err != nil {

		return err

	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	var status string
	err := s.db.QueryRow(
		"SELECT status FROM parcel WHERE number = :number",
		sql.Named("number", number),
	).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("посылка с номером %d не найдена", number)
		}
		return err
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("можно менять адрес только если статус 'registered'")
	}

	_, err = s.db.Exec(
		"UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address), sql.Named("number", number),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var status string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number)).Scan(&status)

	if err != nil {
		return fmt.Errorf("ошибка при проверке статуса: %w", err)
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("можно удалять только посылки со статусом registered")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))

	if err != nil {

		return err
	}

	return nil
}
