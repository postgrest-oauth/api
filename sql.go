package main

import (
	"database/sql"
	"fmt"
	"log"

	"flag"
	_ "github.com/lib/pq"
)

type Owner struct {
	Id                string
	Username          string
	Password          string
	Email             string
	Phone             string
	VerificationCode  string
	VerificationRoute string
}

var dbConnString = flag.String("dbConnString", "postgres://user:pass@localhost:5432/test?sslmode=disable",
	"Database connection string")

func (a *Owner) create() (id string, role string, jti string, err error) {
	db, err := dbConnect()
	defer db.Close()
	query := fmt.Sprintf("SELECT id::varchar, role::varchar, jti::varchar FROM oauth2.create_owner('%s', '%s', '%s', '%s', '%s')",
		a.Email, a.Phone, a.Password, a.VerificationCode, a.VerificationRoute)
	err = db.QueryRow(query).Scan(&id, &role, &jti)

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("looks like owner already exists")
	default:
		log.Printf("User created. ID: %s, ROLE: %s\n", id, role)
	}

	return id, role, jti, err
}

func (a *Owner) check() (id string, role string, jti string, err error) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT id::varchar, role::varchar, jti::varchar FROM oauth2.check_owner('%s', '%s')",
		a.Username, a.Password)
	var uId, uRole, uJti sql.NullString
	err = db.QueryRow(query).Scan(&uId, &uRole, &uJti)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("something bad happened")
	} else if uId.Valid && uRole.Valid && uJti.Valid {
		id, role, jti = uId.String, uRole.String, uJti.String
	} else {
		err = fmt.Errorf("wrong login or password")
	}

	return id, role, jti, err
}

func (a *Owner) verify() (resErr error) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT oauth2.verify_owner('%s')",
		a.Id)
	_, err = db.Query(query)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("owner with id '%s' doesn't exist", a.Id)
	}

	resErr = err
	return resErr
}

func (a *Owner) requestPassword() (resErr error, id string) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT id::varchar FROM oauth2.password_request('%s', '%s', '%s')",
		a.Username, a.VerificationCode, a.VerificationRoute)
	err = db.QueryRow(query).Scan(&id)

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("looks like owner doesn't exist")
	default:
		log.Printf("User exist. ID: %s", id)
	}

	resErr = err
	return resErr, id
}

func (a *Owner) resetPassword() (resErr error) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT oauth2.password_reset('%s', '%s')",
		a.Id, a.Password)
	_, err = db.Query(query)

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("password reset error. USER ID: '%s'", a.Id)
	default:
		log.Printf("Password reseted. USER ID: %s", a.Id)
	}

	resErr = err
	return resErr
}

func (a *Owner) getOwnerRoleAndJtiById() (role string, jti string, err error) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT role::text, jti::text FROM oauth2.owner_role_and_jti_by_id('%s')", a.Id)
	var uRole, uJti sql.NullString
	err = db.QueryRow(query).Scan(&uRole, &uJti)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("something bad happened. Owner ID: '%s'", a.Id)
	} else if uRole.Valid && uJti.Valid {
		role = uRole.String
		jti = uJti.String
	} else {
		err = fmt.Errorf("wrong owner id '%s'", a.Id)
	}

	return role, jti, err
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
