package models

import (
	"fmt"
	"gin-go/db"
	"gin-go/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to create user save query: %w", err)
	}
	defer stmt.Close()

	hashedPwd, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Email, hashedPwd)
	if err != nil {
		return fmt.Errorf("failed to execute user save statement: %w", err)
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to receive last insert id: %w", err)
	}
	u.ID = userId
	return nil
}

func (u *User) ValidateCreds() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retreivedPwd string
	err := row.Scan(&u.ID, &retreivedPwd)
	if err != nil {
		return fmt.Errorf("failed to scan a pwd row: %w", err)
	}

	isPwdValid, err := utils.CheckPwdHash(u.Password, retreivedPwd)
	if err != nil {
		return fmt.Errorf("failed to compare passwords: %w", err)
	}

	if !isPwdValid {
		return fmt.Errorf("password is not valid: %w", err)
	}

	return nil
}
