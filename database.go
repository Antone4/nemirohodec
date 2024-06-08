package main

import (
	"github.com/jackc/pgx"
)

var config = pgx.ConnConfig{
	Host:                 "db",
	User:                 "admin",
	Password:             "postgres",
	Database:             "postgres",
	PreferSimpleProtocol: true,
}

var conn *pgx.Conn

func init() {
	var DBErr error
	conn, DBErr = pgx.Connect(config)
	Printfln("HIIIIIII")
	if DBErr != nil {
		Printfln("Unable to connect to database %v\n", DBErr)
	}

	//Ввод товаров в БД

	// for _, tovar := range tovars {
	// 	_, err = conn.Exec("INSERT INTO tovars (imgref, name, description, price) VALUES ($1, $2, $3, $4)", tovar.ImgRef, tovar.Name, tovar.Description, tovar.Price)
	// 	if err != nil {
	// 		Printfln("Error: %v", err.Error())
	// 		break
	// 	}
	// }

}

func getTovarsFromPGX(conn *pgx.Conn) *[]Tovar {
	tovars := make([]Tovar, 0, 10)
	rows, err := conn.Query("SELECT id, imgref, name, description, price FROM tovars")
	defer rows.Close()
	if err != nil {
		Printfln("Query error: %v", err.Error())
	}
	for rows.Next() {
		var newEl Tovar
		err = rows.Scan(&newEl.Id, &newEl.ImgRef, &newEl.Name, &newEl.Description, &newEl.Price)
		if err != nil {
			Printfln("Scan error: %v", err.Error())
		}
		tovars = append(tovars, newEl)
	}
	return &tovars
}
