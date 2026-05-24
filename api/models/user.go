package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName   string `binding:"required"`
	Password   string
	NickName   string
	IsRemember bool
}

var ErrUserExists = errors.New("user already exists")

func (um *User) GetUser(userName string) (u User) {
	r := DB.Limit(1).Where("user_name=?", userName).Find(&u)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			fmt.Println("user does not exist!")
		}
		// 数据库链接错误记录日志
		fmt.Println(r.Error)
	}
	return
}

// SetPassword hashes the raw password and stores it in the User struct
func (um *User) SetPassword(raw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	um.Password = string(hash)
	return nil
}

// CheckPassword compares a raw password with the stored hash
func (um *User) CheckPassword(raw string) bool {
	if um.Password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(um.Password), []byte(raw))
	return err == nil
}

// CreateUser creates a new user with hashed password
func CreateUser(userName, rawPassword string) (*User, error) {
	// check if username already exists
	var count int64
	DB.Model(&User{}).Where("user_name = ?", userName).Count(&count)
	if count > 0 {
		return nil, ErrUserExists
	}

	u := User{UserName: userName}
	if err := u.SetPassword(rawPassword); err != nil {
		return nil, err
	}
	r := DB.Create(&u)
	return &u, r.Error
}

// DeleteUserByID removes a user by ID
func DeleteUserByID(id uint) error {
	r := DB.Delete(&User{}, id)
	return r.Error
}

// GetAllUsers returns all users (excluding password field for safety is left to caller)
func GetAllUsers() ([]User, error) {
	var users []User
	r := DB.Find(&users)
	return users, r.Error
}
