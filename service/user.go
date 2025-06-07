package service

import (
	"context"

	dao "github.com/sztu/mutli-table/DAO"
)

func ListUsersService(ctx context.Context) ([]string, error) {
	userList, err := dao.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, user := range userList {
		userName := user.Username
		result = append(result, userName)
	}
	return result, nil
}
