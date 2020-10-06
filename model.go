package main

import (
	"database/sql"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

type book struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type User struct {
	ID        uint32    `json:"id"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *book) getBook(db *sql.DB) error {
	return db.QueryRow(`SELECT "Name", "Price" FROM "Book" WHERE "id"=$1`,
		p.ID).Scan(&p.Name, &p.Price)
}

func (p *book) updateBook(db *sql.DB) error {
	_, err :=
		db.Exec(`UPDATE public."Book" SET "Name"=$1, "Price"=$2 WHERE "id"=$3`,
			p.Name, p.Price, p.ID)

	return err
}

func (p *book) deleteBook(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM public."Book" WHERE "id"=$1`, p.ID)

	return err
}

func (p *book) createBook(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO book(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getBookss(db *sql.DB, start, count int) ([]book, error) {
	rows, err := db.Query(
		`SELECT "id", "Name",  "Price" FROM "Book" LIMIT ALL OFFSET 0`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []book{}

	for rows.Next() {
		var p book
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		books = append(books, p)
	}

	return books, nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

/*
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
*/
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
