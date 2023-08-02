package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// connect to a database

	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=test_connect user=postgres password=postgres")

	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect: %v\n", err))
	}
	defer conn.Close() // it will close when the main() function finishes executing
	log.Println("Connected to database!")

	// test my connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot ping database!")
	}
	log.Println("Pinged database")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// insert a row
	query := `insert into users (first_name, last_name) values ($1, $2)`
	_, err = conn.Exec(query, "Jack", "Brown")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a row!")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// update a row
	// stmt := `update users set first_name = $1 where first_name = $2`
	// _, err = conn.Exec(stmt, "Jackie", "Jack")
	stmt := `update users set first_name = $1 where id = $2`
	_, err = conn.Exec(stmt, "Jackie", 7)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Updated one or more rows")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// get one row by id
	query = `select id, first_name, last_name from users where id = $1`
	var firstName, lastName string
	var id int

	row := conn.QueryRow(query, 1) // returns at most one row
	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("QueryRow returns", id, firstName, lastName)

	// delete a row
	query = `delete from users where id = $1`
	_, err = conn.Exec(query, 6)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deleted a row!")

	// get rows again
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close() // we have to close rows after a query!

	var firstName, lastName string
	var id int

	for rows.Next() { // it ranges over rows
		err := rows.Scan(&id, &firstName, &lastName) // Scan places the values under the variables in the order of select
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("record is", id, firstName, lastName)
	}
	// we might get errors while scanning through
	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows", err)
	}
	fmt.Println("-------------------------")
	return nil
}
