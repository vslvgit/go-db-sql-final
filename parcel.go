package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

res, _ := s.db.Exec("INSERT INTO parcel VALUES (:Number :Client :Status :Address :CreatedAt)", sql.Named("Number", p.Number), sql.Named("Client", p.Client), sql.Named("Status", p.Status), sql.Named("Address", p.Address), sql.Named("CreatedAt", p.CreatedAt))

	id, err := res.LastInsertId()

	if err != nil {

return 0, err

	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

row := s.db.QueryRow("SELECT number client status address createdAt FROM parcel WHERE id = :id", sql.Named("id", number))
	p := Parcel{}
  err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
if err != nil {

return Parcel{}, err

}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

rows, err := s.db.Query(
		"SELECT number, client, status, address, createdAt FROM parcel WHERE client = :client", sql.Named("client", client),
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
_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE id = :id", sql.Named("status", status), sql.Named("id", number))

if err != nil {

return err

}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE id = :id", sql.Named("address", address), sql.Named("id", number))

if err != nil {

return err

}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var status string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE id = :id", sql.Named("id", number)).Scan(&status)

	if status == "registered" {
  _, err = s.db.Exec("DELETE FROM parcel WHERE id = :id", sql.Named("id", number))
    if err != nil {
        return err 
		}
	
}
return nil
}