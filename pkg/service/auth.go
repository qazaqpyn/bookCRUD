package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	repo *repository.Repository
}

const (
	sessionTtl = 12 * time.Minute
	salt       = "ahsjdf1412_0293@sifnsHIUfb123UHI"
)

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) error {
	user.Password = generatePasswordHash(user.Password)
	user.Id = primitive.NewObjectID()

	return s.repo.CreateUser(ctx, &user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) SignIn(ctx context.Context, inp model.LoginInput) (string, error) {
	password := generatePasswordHash(inp.Password)

	_, err := s.repo.Authorization.GetUser(ctx, inp.Email, password)
	if err != nil {
		return "", nil
	}

	return generateSession()
}

func (s *AuthService) CheckSessionID(ctx context.Context, session_id string) error {
	sess, err := s.repo.Tokens.Get(ctx, session_id)
	if err != nil {
		return err
	}

	if sess.ExpiresAt.Unix() < time.Now().Unix() {
		return errors.New("session_id is expired login again")
	}

	return nil
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
