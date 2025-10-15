package service

import (
	"net/smtp"

	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/liceum_backend/internal/model"
	"github.com/lavatee/liceum_backend/internal/repository"
)

type Events interface {
	SendAuthCode(email string) error
	VerifyCode(code string, email string) (string, string, error)
	CheckIsAdmin(email string) bool
	CreateEvent(event model.Event) (int, error)
	DeleteEvent(eventId int) error
	CreateEventBlocks(blocks []model.EventBlock) error
	DeleteEventBlock(blockId int) error
	EditEventInfo(event model.Event) error
	EditBlockInfo(block model.EventBlock) error
	GetCurrentEvents() ([]model.Event, error)
	GetAllEvents() ([]model.Event, error)
	ParseToken(token string) (jwt.MapClaims, error)
	GetOneEvent(eventId int) (model.Event, error)
	GetOneBlock(blockId int) (model.EventBlock, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type Service struct {
	Events
}

func NewService(repo *repository.Repository, smtpAuth smtp.Auth, gmail string, smtpHost string, smtpPort string) *Service {
	return &Service{
		Events: NewEventsService(repo, smtpAuth, gmail, smtpHost, smtpPort),
	}
}
