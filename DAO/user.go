package dao

import (
	"context"

	mysql "github.com/sztu/mutli-table/DAO/MySQL"
	"github.com/sztu/mutli-table/model"
)

func CreateUser(ctx context.Context, user *model.User) error {
	sqlStr := `INSERT INTO user (user_id, username, password, email) VALUES (?, ?, ?, ?)`
	return mysql.GetDB().WithContext(ctx).Exec(sqlStr, user.UserID, user.Username, user.Password, user.Email).Error
}

func FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	sqlStr := `SELECT user_id, username, password FROM user WHERE username = ? AND delete_time = 0`
	result := mysql.GetDB().WithContext(ctx).Raw(sqlStr, username).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, nil
}

func FindUserByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	sqlStr := `SELECT user_id, username, password FROM user WHERE user_id = ? AND delete_time = 0`
	err := mysql.GetDB().WithContext(ctx).Raw(sqlStr, userID).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
