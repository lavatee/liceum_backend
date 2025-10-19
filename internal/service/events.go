package service

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/liceum_backend/internal/model"
	"github.com/lavatee/liceum_backend/internal/repository"
	"github.com/sirupsen/logrus"
)

var getRequestsCounter = 0

const (
	tokenKey   = "miron_huesos_123"
	accessTTL  = 15 * time.Minute
	refreshTTL = 15 * 24 * time.Hour
)

type CodeStore struct {
	mu    sync.RWMutex
	codes map[string][2]interface{}
}

func NewCodeStore() *CodeStore {
	return &CodeStore{
		codes: make(map[string][2]interface{}),
	}
}

func (cs *CodeStore) SetCode(userEmail string, code string) {
	expiration := time.Now().Add(60 * time.Second)
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.codes[userEmail] = [2]interface{}{code, expiration}
}

func (cs *CodeStore) VerifyCode(userEmail string, code string) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	data, exists := cs.codes[userEmail]
	if !exists {
		return false
	}

	storedCode := data[0].(string)
	expiration := data[1].(time.Time)
	if time.Now().After(expiration) {
		delete(cs.codes, userEmail)
		return false
	}
	return storedCode == code
}

type EventsService struct {
	repo      *repository.Repository
	smtpAuth  smtp.Auth
	gmail     string
	smtpHost  string
	smtpPort  string
	codeStore *CodeStore
}

func NewEventsService(repo *repository.Repository, smtpAuth smtp.Auth, gmail string, smtpHost string, smtpPort string) *EventsService {
	return &EventsService{
		repo:      repo,
		smtpAuth:  smtpAuth,
		gmail:     gmail,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
		codeStore: NewCodeStore(),
	}
}

var adminsMap = map[string]bool{
	"gorodilow.aleksandr@gmail.com": true,
	"aleksgraznov0@gmail.com":       true,
}

func (s *EventsService) CheckIsAdmin(email string) bool {
	_, ok := adminsMap[email]
	return ok
}

func generateRandomCode() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(900000) + 100000
}

func (s *EventsService) SendAuthCode(email string) error {
	if !s.CheckIsAdmin(email) {
		return fmt.Errorf("notadmin")
	}
	code := strconv.Itoa(generateRandomCode())
	if err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, s.smtpAuth, s.gmail, []string{email}, []byte(fmt.Sprintf("Subject: Ваш код для входа в аккаунт: %s", code))); err != nil {
		return err
	}
	s.codeStore.SetCode(email, code)
	return nil
}

func (s *EventsService) VerifyCode(code string, email string) (string, string, error) {
	if !s.codeStore.VerifyCode(email, code) {
		return "", "", fmt.Errorf("wrongcode")
	}
	accessClaims := jwt.MapClaims{
		"exp":   time.Now().Add(accessTTL).Unix(),
		"email": email,
	}
	refreshClaims := jwt.MapClaims{
		"exp":   time.Now().Add(refreshTTL).Unix(),
		"email": email,
	}
	accessToken, err := s.NewToken(accessClaims)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.NewToken(refreshClaims)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *EventsService) ParseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("expired token")
}

func (s *EventsService) NewToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}
	return stringToken, nil
}

func (s *EventsService) CreateEvent(event model.Event) (int, error) {
	return s.repo.Events.CreateEvent(event)
}

func (s *EventsService) DeleteEvent(eventId int) error {
	return s.repo.Events.DeleteEvent(eventId)
}

func (s *EventsService) CreateEventBlocks(blocks []model.EventBlock) error {
	if len(blocks) == 0 {
		return fmt.Errorf("blocks is empty")
	}
	eventId := blocks[0].EventID
	return s.repo.Events.CreateEventBlocks(blocks, eventId)
}

func (s *EventsService) DeleteEventBlock(blockId int) error {
	return s.repo.Events.DeleteEventBlock(blockId)
}

func (s *EventsService) EditEventInfo(event model.Event) error {
	return s.repo.Events.EditEventInfo(event)
}

func (s *EventsService) EditBlockInfo(block model.EventBlock) error {
	return s.repo.Events.EditBlockInfo(block)
}

func (s *EventsService) GetCurrentEvents() ([]model.Event, error) {
	return s.repo.Events.GetCurrentEvents()
}

func (s *EventsService) GetAllEvents() ([]model.Event, error) {
	getRequestsCounter += 1
	if getRequestsCounter == 20 {
		getRequestsCounter = 0
		if err := s.repo.Events.CleanEvents(); err != nil {
			logrus.Errorf("CLEAN EVENTS ERROR: %s", err.Error())
		}
	}
	return s.repo.Events.GetAllEvents()
}

func (s *EventsService) GetOneEvent(eventId int) (model.Event, error) {
	return s.repo.Events.GetOneEvent(eventId)
}

func (s *EventsService) GetOneBlock(blockId int) (model.EventBlock, error) {
	return s.repo.Events.GetOneBlock(blockId)
}

func (s *EventsService) RefreshToken(refreshToken string) (string, string, error) {
	parsedToken, err := s.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	accessClaims := jwt.MapClaims{
		"exp":   time.Now().Add(accessTTL).Unix(),
		"email": parsedToken["email"],
	}
	refreshClaims := jwt.MapClaims{
		"exp":   time.Now().Add(refreshTTL).Unix(),
		"email": parsedToken["email"],
	}
	accessToken, err := s.NewToken(accessClaims)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = s.NewToken(refreshClaims)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
