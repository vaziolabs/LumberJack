package core

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// SetPassword sets the password for the user
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// VerifyPassword verifies the password for the user
func (u *User) VerifyPassword(password string) bool {
	fmt.Println("Verifying password for user:", u.Username)
	fmt.Println("Password:", u.Password)
	fmt.Println("Input password:", password)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
