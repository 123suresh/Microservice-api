package repository

import (
	"fmt"

	"example.com/dynamicWordpressBuilding/internal/model"
	"github.com/google/uuid"
)

type UserInterface interface {
	CreateUser(req *model.User) (*model.User, error)
	GetAllUser() ([]model.User, error)
	GetUser(uid uuid.UUID) (*model.User, error)
	DeleteUser(uid uuid.UUID) error
	LoginUser(email string) (*model.User, error)
	EmailExistCheck(email string) bool
	ResetPassword(resetPass *model.ResetPassword) (*model.ResetPassword, error)
	FindByToken(token string) (*model.ResetPassword, error)
	UpdatePassword(resetDetail *model.ResetPassword, resetPassReq *model.ForgetPassword) error
}

func (r *Repo) CreateUser(data *model.User) (*model.User, error) {
	err := r.db.Model(&model.User{}).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Repo) GetAllUser() ([]model.User, error) {
	users := []model.User{}
	err := r.db.Model(&model.User{}).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repo) GetUser(uid uuid.UUID) (*model.User, error) {
	data := &model.User{}
	err := r.db.Model(&model.User{}).Where("id=?", uid).Take(&data).Error
	if err != nil {
		return nil, fmt.Errorf("user doesn't exists %v ", err)
	}
	return data, nil
}

func (r *Repo) DeleteUser(uid uuid.UUID) error {
	data := &model.User{}
	err := r.db.Model(&model.User{}).Where("id=?", uid).Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) LoginUser(email string) (*model.User, error) {
	data := &model.User{}
	err := r.db.Model(&model.User{}).Where("email=?", email).Take(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Repo) EmailExistCheck(email string) bool {
	data := &model.User{}
	err := r.db.Model(&model.User{}).Where("email = ?", email).Take(&data).Error
	if err == nil {
		return true
	}
	return false
}

func (r *Repo) ResetPassword(resetPass *model.ResetPassword) (*model.ResetPassword, error) {
	err := r.db.Model(&model.ResetPassword{}).Create(resetPass).Error
	if err != nil {
		return nil, fmt.Errorf("error while doing reset password %v ", err)
	}
	return resetPass, nil
}

func (r *Repo) FindByToken(token string) (*model.ResetPassword, error) {
	data := &model.ResetPassword{}
	err := r.db.Model(&model.ResetPassword{}).Where("token = ?", token).Take(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Repo) UpdatePassword(resetDetail *model.ResetPassword, resetPassReq *model.ForgetPassword) error {
	err := r.db.Model(&model.User{}).Where("email = ?", resetDetail.Email).UpdateColumns(
		map[string]interface{}{
			"password": resetPassReq.NewPassword,
		}).Error
	if err != nil {
		return err
	}
	return nil
}
