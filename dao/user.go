package dao

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id              int64
	Username        string
	Password        string
	Email           string
	App_secret      string
	Last_Login_Ip   string
	Last_Login_Time time.Time
	Status          int
	Updated_At      time.Time
	Created_At      time.Time
}

type UserRedis struct {
	Id       int64
	Username string
	Email    string
	Status   int
}

type UserDAO struct {
	db *gorm.DB
}

var tableName = "t_user_info"

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func ConvertToUserRedis(user *User) *UserRedis {
	return &UserRedis{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}
}

func (u *UserDAO) GetByUsername(username string) (*User, error) {
	user := User{}
	err := u.db.Table(tableName).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserDAO) Add(user interface{}) error {
	return u.db.Table(tableName).Create(user).Error
}

func (u *UserDAO) Update(id, info interface{}) error {
	return u.db.Table(tableName).Where("id = ?", id).Updates(info).Error
}
