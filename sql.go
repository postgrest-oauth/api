package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env"
	_ "github.com/lib/pq"
)

type Owner struct {
	Id                string
	FacebookId        string
	FacebookJson      string
	Username          string
	Password          string
	Email             string
	Phone             string
	VerificationCode  string
	VerificationRoute string
}

var sqlConfig struct {
	DbConnString string `env:"OAUTH_DB_CONN_STRING" envDefault:"postgres://user:pass@localhost:5432/test?sslmode=disable"`
}

func init() {
	err := env.Parse(&sqlConfig)
	if err != nil {
		log.Printf("%+v\n", err)
	}
}

func (a *Owner) create() (id string, err error) {
	db, err := dbConnect()
	defer db.Close()
	query := fmt.Sprintf("SELECT id::varchar FROM oauth2.create_owner('%s', '%s', '%s', '%s', '%s')",
		a.Email, a.Phone, a.Password, a.VerificationCode, a.VerificationRoute)
	err = db.QueryRow(query).Scan(&id)

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("looks like owner already exists")
	default:
		log.Printf("User created. ID: %s\n", id)
	}

	return id, err
}

func (a *Owner) create_or_update_facebook(json string, phone string, lang string) (id string, role string, jti string, err error) {
	db, err := dbConnect()
	defer db.Close()
	query := fmt.Sprintf("SELECT id::varchar, role::varchar, jti::varchar FROM oauth2.create_or_update_facebook_owner('%s'::json,'%s', '%s')",
		json, phone, lang)
	var uId, uRole, uJti sql.NullString
	err = db.QueryRow(query).Scan(&uId, &uRole, &uJti)

	id = uId.String
	role = uRole.String
	jti = uJti.String

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("issue with database")
	default:
		log.Printf("User created or updated. ID: %s\n", id)
	}

	return id, role, jti, err
}

func (a *Owner) reVerify() (id string, err error) {
	db, err := dbConnect()
	defer db.Close()
	query := fmt.Sprintf("SELECT id::varchar FROM oauth2.re_verify('%s', '%s', '%s')",
		a.Username, a.VerificationCode, a.VerificationRoute)
	err = db.QueryRow(query).Scan(&id)

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("looks like owner '%s' doesn't exists or verified", a.Username)
	default:
		log.Printf("Verification code re-sent. ID: %s\n", id)
	}

	return id, err
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

func (a *Owner) check_facebook() (id string, role string, jti string, err error) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT id::varchar, role::varchar, jti::varchar FROM oauth2.check_owner_facebook('%s')",
		a.FacebookId)
	var uId, uRole, uJti sql.NullString
	err = db.QueryRow(query).Scan(&uId, &uRole, &uJti)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("something bad happened")
	} else if uId.Valid && uRole.Valid && uJti.Valid {
		id, role, jti = uId.String, uRole.String, uJti.String
	} else {
		err = fmt.Errorf("wrong facebook id")
	}

	return id, role, jti, err
}

func (a *Owner) verify() (resErr error) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT oauth2.verify_owner('%s')",
		a.Id)
	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("owner with id '%s' doesn't exist", a.Id)
	}

	resErr = err
	return resErr
}

func (a *Owner) requestPassword() (id string, resErr error) {
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
	return id, resErr
}

func (a *Owner) resetPassword() (resErr error) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT oauth2.password_reset('%s', '%s')",
		a.Id, a.Password)
	rows, err := db.Query(query)
	defer rows.Close()

	switch {
	case err != nil:
		log.Print(err)
		err = fmt.Errorf("password reset error. USER ID: '%s'", a.Id)
	default:
		log.Printf("password reseted. USER ID: %s", a.Id)
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

type Client struct {
	Id     string
	Secret string
}

func (c *Client) check() (resErr error, redirectUri string) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT redirect_uri::text FROM oauth2.check_client('%s')",
		c.Id)
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

func (c *Client) check_secret() (resErr error, ctype string) {
	db, err := dbConnect()
	defer db.Close()

	query := fmt.Sprintf("SELECT type::varchar FROM oauth2.check_client_secret('%s', '%s')",
		c.Id, c.Secret)
	var uType sql.NullString
	err = db.QueryRow(query).Scan(&uType)

	if err != nil {
		log.Print(err)
		err = fmt.Errorf("something bad happened. Client ID: '%s'", c.Id)
	} else if uType.Valid {
		ctype = uType.String
	} else {
		err = fmt.Errorf("wrong client id '%s'", c.Id)
	}

	resErr = err
	return resErr, ctype
}

func dbConnect() (*sql.DB, error) {
	return sql.Open("postgres", sqlConfig.DbConnString)
}
