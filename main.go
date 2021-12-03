package main

import (

	//"fmt"
	"log"
	"net/http"

	//"database/sql"

	"library/handler"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)



func main() { 

	var createLibrary = `
	CREATE TABLE IF NOT EXISTS books (
		id serial,
		catid int,
    	bookName text,
		Image text,
    	is_available boolean,

		primary key(id)
);

	CREATE TABLE IF NOT EXISTS categories (
		id serial,
		categoryName text,
		
		primary key(id)
);

	CREATE TABLE IF NOT EXISTS bookings (
		id serial,
		user_id integer,
		book_id integer,
		start_time timestamp,
		end_time   timestamp,


		primary key(id)

);`
var createUserTable = `
CREATE TABLE IF NOT EXISTS users (
	id serial,
	first_name text,
	last_name text,
	email text,
	username text,
	password text,
	is_verified boolean,

	primary key(id)
);`



	db, err := sqlx.Connect("postgres", "user=akib password=password dbname=library sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(createLibrary)
	db.MustExec(createUserTable)

	decoder := schema.NewDecoder()

	decoder.IgnoreUnknownKeys(true)
	
	store := sessions.NewCookieStore([]byte("ffhfhfsfdfsfddsgfgfgfg"))
	r := handler.New(db, decoder, store)

	
	log.Println("Server starting....")
	if err := http.ListenAndServe("127.0.0.1:4000", r); err != nil {
		log.Fatal(err)
	}
}



