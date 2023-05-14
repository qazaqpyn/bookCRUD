package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/logging"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	audit "github.com/qazaqpyn/crud-audit-log/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	repo        *repository.Repository
	auditClient AuditClient
}

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

const (
	sessionTtl = 12 * time.Minute
	salt       = "ahsjdf1412_0293@sifnsHIUfb123UHI"
)

func NewAuthService(repo *repository.Repository, auditClient AuditClient) *AuthService {
	return &AuthService{
		repo:        repo,
		auditClient: auditClient,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) error {
	user.Password = generatePasswordHash(user.Password)
	user.Id = primitive.NewObjectID()

	if err := s.repo.CreateUser(ctx, &user); err != nil {
		return err
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_CREATE,
		Entity:    audit.ENTITY_USER,
		EntityID:  user.Id.Hex(),
		Timestamp: time.Now(),
	}); err != nil {
		logging.LogError("Users.SignUp", err)
	}

	return nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) SignIn(ctx context.Context, inp model.LoginInput) (string, error) {
	password := generatePasswordHash(inp.Password)

	user, err := s.repo.Authorization.GetUser(ctx, inp.Email, password)
	if err != nil {
		return "", nil
	}

	session_id, err := generateSession()
	if err != nil {
		return "", nil
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_LOGIN,
		Entity:    audit.ENTITY_USER,
		EntityID:  user.Id.Hex(),
		Timestamp: time.Now(),
	}); err != nil {
		logging.LogError("Users.SignIn", err)
	}

	return session_id, nil
}

func (s *AuthService) CheckSessionID(ctx context.Context, session_id string) (string, error) {
	sess, err := s.repo.Tokens.Get(ctx, session_id)
	if err != nil {
		return "", err
	}

	if sess.ExpiresAt.Unix() < time.Now().Unix() {
		return "", errors.New("session_id is expired login again")
	}

	return sess.Id.Hex(), nil
}

func generateSession() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
