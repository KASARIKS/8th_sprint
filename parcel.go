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
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at);",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	lastInd, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(lastInd), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number;", sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client;", sql.Named("client", client))
	if err != nil {
		return []Parcel{}, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		tmpParcel := Parcel{}
		if err = rows.Scan(&tmpParcel.Number, &tmpParcel.Client, &tmpParcel.Status, &tmpParcel.Address, &tmpParcel.CreatedAt); err != nil {
			return []Parcel{}, err
		}

		res = append(res, tmpParcel)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number;",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	isReg, err := s.isNumberRegistered(number)
	if err == nil && isReg {
		_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number;",
			sql.Named("address", address),
			sql.Named("number", number))
	}

	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	isReg, err := s.isNumberRegistered(number)
	if err == nil && isReg {
		_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number;", sql.Named("number", number))
	}

	return err
}

func (s ParcelStore) isNumberRegistered(number int) (bool, error) {
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number;", sql.Named("number", number))
	var res string
	err := row.Scan(&res)
	if err != nil {
		return false, err
	}

	if res == ParcelStatusRegistered {
		return true, nil
	}

	return false, nil
}
