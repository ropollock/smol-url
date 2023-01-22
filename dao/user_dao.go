package dao

import (
	"fmt"
	"github.com/google/uuid"
	"server/data"
	"server/model"
)

type userDao struct {
	databaseProvider data.PostgresDBProviderInterface
}

type UserDaoInterface interface {
	CreateUser(user *model.User) (*model.User, error)
	DeleteUser(user *model.User) error
	FindUserById(id string) (model.User, error)
	FindUserByUsername(username string) (model.User, error)
	GetUsers() ([]model.User, error)
}

func UserDao(databaseProvider data.PostgresDBProviderInterface) *userDao {
	return &userDao{databaseProvider}
}

func (dao *userDao) CreateUser(user *model.User) (*model.User, error) {
	db := dao.databaseProvider.GetDB()
	insertResult := db.Create(user)
	result, _ := dao.FindUserById(user.ID.String())
	return &result, insertResult.Error
}

func (dao *userDao) DeleteUser(user *model.User) error {
	deleteResult := dao.databaseProvider.GetDB().Delete(user)
	return deleteResult.Error
}

func (dao *userDao) FindUserById(id string) (model.User, error) {
	result := model.User{ID: uuid.MustParse(id)}
	findResult := dao.databaseProvider.GetDB().First(&result)
	if findResult.Error != nil {
		fmt.Println(findResult.Error)
		return result, fmt.Errorf("an error occurred while decoding record : %v", findResult.Error)
	}
	return result, nil
}

func (dao *userDao) FindUserByUsername(username string) (model.User, error) {
	result := model.User{}
	findResult := dao.databaseProvider.GetDB().Where("username = ?", username).First(&result)
	if findResult.Error != nil {
		fmt.Println(findResult.Error)
		return result, fmt.Errorf("an error occurred while decoding record : %v", findResult.Error)
	}
	return result, nil
}

func (dao *userDao) GetUsers() ([]model.User, error) {
	var results []model.User
	findResult := dao.databaseProvider.GetDB().Find(&results)

	if findResult.Error != nil {
		fmt.Println("Finding all users ERROR:", findResult.Error)
		return results, findResult.Error
	}

	return results, nil
}
