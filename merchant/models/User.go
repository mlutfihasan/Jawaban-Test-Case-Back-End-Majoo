package models

import (
	"errors"
	"html"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CrudResult struct {
	Status string `json:"status"`
	Note   error  `json:"note"`
}

// User is model
type User struct {
	Id       int    `gorm:"auto_increment;primary_key;" json:"id"`
	Nik      string `gorm:"size:30;not null;index:idx_users;" json:"nik"`
	Username string `gorm:"size:50;index:idx_users;" json:"username"`
	Phone    string `gorm:"size:20;index:idx_users;" json:"phone"`
	Email    string `gorm:"size:100;index:idx_users;" json:"email"`
	Password string `gorm:"size:200;" json:"password"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword is ...
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Nik = html.EscapeString(strings.TrimSpace(u.Nik))
	u.Phone = html.EscapeString(strings.TrimSpace(u.Phone))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("required username")
	}
	if u.Password == "" {
		return errors.New("required password")
	}
	if u.Nik == "" {
		return errors.New("required nik")
	}
	if u.Phone == "" {
		return errors.New("required phone")
	}
	if u.Email == "" {
		return errors.New("required email")
	}
	return nil
}

func (u *User) SaveUser(db *gorm.DB) CrudResult {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return CrudResult{
			Status: "0",
			Note:   err,
		}
	}

	return CrudResult{
		Status: "1",
		Note:   nil,
	}
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var users []User

	err := db.Debug().Model(&User{}).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}

	return &users, nil
}

func (u *User) FindUserByID(db *gorm.DB, nik string) (*User, error) {
	err := db.Debug().Model(&User{}).Where(&User{Nik: nik}).Take(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, nik string) CrudResult {
	db = db.Debug().Where(&User{Nik: nik}).Take(&User{}).Updates(
		User{
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
		},
	)
	if db.Error != nil {
		return CrudResult{
			Status: "0",
			Note:   db.Error,
		}
	}

	return CrudResult{
		Status: "1",
		Note:   nil,
	}
}

func (u *User) DeleteAUser(db *gorm.DB) CrudResult {
	db = db.Debug().Model(&User{}).Where(&User{Nik: u.Nik}).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return CrudResult{
			Status: "0",
			Note:   db.Error,
		}
	}

	return CrudResult{
		Status: "1",
		Note:   nil,
	}
}

func (u *User) AuthenticateLogin(db *gorm.DB) error {
	var ul User
	db = db.Debug().Model(&User{}).Where(&User{Nik: u.Nik}).Take(&ul)
	if db.Error != nil {
		return db.Error
	}

	HashedCurrPass := ul.Password

	if HashedCurrPass == "123" {
		HashedPass, err := bcrypt.GenerateFromPassword([]byte(HashedCurrPass), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("ERROR:", err.Error())
		}
		HashedCurrPass = string(HashedPass)
	}

	err := VerifyPassword(HashedCurrPass, u.Password)
	if err != nil {
		return err
	}

	return nil
}
