package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Owner struct {
	Id       string
	Username string
	Password string
	Email    string
	Phone    string
}

func (a *Owner) create() (resErr error, id string, role string) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT id::varchar, role::varchar FROM oauth2.create_owner('%s', '%s', '%s')",
		a.Email, a.Phone, a.Password)
	err = db.QueryRow(query).Scan(&id, &role)

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("looks like owner already exists")
	default:
		log.Printf("User created. ID: %s, ROLE: %s", id, role)
	}

	resErr = err
	return resErr, id, role
}

func (a *Owner) check() (resErr error, id string, role string) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT id::varchar, role::varchar FROM oauth2.check_owner('%s', '%s')",
		a.Username, a.Password)
	var uId, uRole sql.NullString
	err = db.QueryRow(query).Scan(&uId, &uRole)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("something bad happened")
	} else if uId.Valid && uRole.Valid {
		id, role = uId.String, uRole.String
	} else {
		err = fmt.Errorf("wrong login or password")
	}

	resErr = err
	return resErr, id, role
}

func (a *Owner) getOwnerRoleById() (resErr error, role string) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT role::text FROM oauth2.owner_role_by_id('%s')", a.Id)
	var uRole sql.NullString
	err = db.QueryRow(query).Scan(&uRole)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("something bad happened. Owner ID: '%s'", a.Id)
	} else if uRole.Valid {
		role = uRole.String
	} else {
		err = fmt.Errorf("wrong owner id '%s'", a.Id)
	}

	resErr = err
	return resErr, role
}

type Client struct {
	Id     string
	Secret string
}

func (c *Client) check() (resErr error, redirectUri string) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT redirect_uri::text FROM oauth2.check_client('%s', '%s')",
		c.Id, c.Secret)
	var uRedirectUri sql.NullString
	err = db.QueryRow(query).Scan(&uRedirectUri)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("something bad happened. Client ID: '%s'", c.Id)
	} else if uRedirectUri.Valid {
		redirectUri = uRedirectUri.String
	} else {
		err = fmt.Errorf("wrong client id '%s'", c.Id)
	}

	resErr = err
	return resErr, redirectUri
}

func dbConnect() (*sql.DB, error) {
	return sql.Open("postgres", *dbConnString)
}
