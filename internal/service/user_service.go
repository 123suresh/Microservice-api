package service

import (
	"errors"
	"net/http"
	"os"
	"time"

	"example.com/dynamicWordpressBuilding/internal/model"
	"example.com/dynamicWordpressBuilding/utils"
	cerr "example.com/dynamicWordpressBuilding/utils/error"
	"example.com/dynamicWordpressBuilding/utils/security"
	"github.com/google/uuid"
)

type IUser interface {
	CreateUser(req *model.UserRequest) (*model.UserResponse, int, error)
	GetAllUser() ([]model.UserResponse, int, error)
	GetUser(uid uuid.UUID) (*model.UserResponse, int, error)
	DeleteUser(uid uuid.UUID) (int, error)
	LoginUser(loginReq *model.LoginRequest) (*model.LoginToken, int, error)
	ResetPassword(forgetPassReq *model.ResetPasswordReq) (*model.ResetPassword, int, error)
	ForgetPassword(resetPassReq *model.ForgetPassword) (int, error)
}

func (s Service) CreateUser(req *model.UserRequest) (*model.UserResponse, int, error) {
	user := model.NewUser(req)
	//check email
	emailExist := s.repo.EmailExistCheck(user.Email)
	if emailExist {
		return nil, http.StatusBadRequest, errors.New("email already taken")
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	user.Password = hashedPassword
	result, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	response := result.UserRes()
	return response, http.StatusCreated, nil
}

func (s Service) GetAllUser() ([]model.UserResponse, int, error) {
	result, err := s.repo.GetAllUser()
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	responses := []model.UserResponse{}
	for _, user := range result {
		responses = append(responses, *user.UserRes())
	}
	return responses, http.StatusOK, nil
}

func (s Service) GetUser(uid uuid.UUID) (*model.UserResponse, int, error) {
	result, err := s.repo.GetUser(uid)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	response := result.UserRes()
	return response, http.StatusOK, nil
}

func (s Service) DeleteUser(uid uuid.UUID) (int, error) {
	err := s.repo.DeleteUser(uid)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func (s Service) LoginUser(loginReq *model.LoginRequest) (*model.LoginToken, int, error) {
	userData, err := s.repo.LoginUser(loginReq.Email)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	err = utils.CheckPassword(loginReq.Password, userData.Password)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	accessTokenDuration := os.Getenv("ACCESS_TOKEN_DURATION")
	duration, err := time.ParseDuration(accessTokenDuration)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	//then we can use CreateToken function
	accessToken, err := s.tokenMaker.CreateToken(userData.ID, userData.Email, duration)

	//for using struct make model and use like this
	// rrr := utils.JWTMaker{}
	// accessToken, err := rrr.CreateToken(userData.Email, duration)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	response := &model.LoginToken{
		AccessToken: accessToken,
		User:        userData.UserRes(),
	}
	return response, http.StatusOK, nil
}

func (s Service) ResetPassword(forgetPassReq *model.ResetPasswordReq) (*model.ResetPassword, int, error) {
	userEmail := s.repo.EmailExistCheck(forgetPassReq.Email)
	if !userEmail {
		return nil, http.StatusBadRequest, cerr.ErrRequiredEmail
	}
	token := security.TokenHash(forgetPassReq.Email)
	resetPassword := &model.ResetPassword{}
	resetPassword.Email = forgetPassReq.Email
	resetPassword.Token = token
	resetDetails, err := s.repo.ResetPassword(resetPassword)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	//Sending mail for reset
	s.mailJet.SendResetPassword(resetDetails.Email, resetDetails.Token)
	return resetDetails, http.StatusOK, nil
}

func (s Service) ForgetPassword(resetPassReq *model.ForgetPassword) (int, error) {
	userDetails, err := s.repo.FindByToken(resetPassReq.Token)
	if err != nil {
		return http.StatusNotFound, nil
	}
	if !userDetails.DeletedAt.Time.IsZero() {
		return http.StatusUnprocessableEntity, cerr.ErrTokenExpired
	}
	if resetPassReq.NewPassword == "" || resetPassReq.RetypePassword == "" {
		return http.StatusUnprocessableEntity, cerr.ErrEmptyPassword
	}
	if resetPassReq.NewPassword != "" && resetPassReq.RetypePassword != "" {
		if len(resetPassReq.NewPassword) < 6 || len(resetPassReq.RetypePassword) < 6 {
			return http.StatusUnprocessableEntity, cerr.ErrPasswordLen
		}
		if resetPassReq.NewPassword != resetPassReq.RetypePassword {
			return http.StatusUnprocessableEntity, cerr.ErrPasswordMatch
		}
		hashedPassword, err := utils.HashPassword(resetPassReq.NewPassword)
		if err != nil {
			return http.StatusBadRequest, err
		}
		resetPassReq.NewPassword = hashedPassword
		err = s.repo.UpdatePassword(userDetails, resetPassReq)
		if err != nil {
			return http.StatusBadRequest, cerr.ErrUpdatePass
		}
	}
	return http.StatusOK, nil
}
