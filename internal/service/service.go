package service

import (
	"example.com/dynamicWordpressBuilding/internal/repository"
	"example.com/dynamicWordpressBuilding/utils"
	"example.com/dynamicWordpressBuilding/utils/mailer"
)

type Service struct {
	repo       repository.RepoInterface
	tokenMaker utils.Maker
	mailJet    mailer.MailService
}

type ServiceInterface interface {
	IUser
}

func NewService(repo repository.RepoInterface) ServiceInterface {
	svc := &Service{}
	svc.repo = repo
	svc.repo = repository.NewRepo()
	svc.tokenMaker = utils.NewTokenMaker()
	svc.mailJet = mailer.NewMailService()
	return svc
}
