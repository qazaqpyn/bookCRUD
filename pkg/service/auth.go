package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	repo *repository.Repository
}

type TokenClaim struct {
	jwt.StandardClaims
	UserID primitive.ObjectID `json:"user_id"`
}

const (
	tokenTTL   = 12 * time.Hour
	salt       = "sdf12341asd_3423resdf1"
	signingKey = "asdjfji12#$fdo13__34123joisdf"
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

func (s *AuthService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	password = generatePasswordHash(password)

	user, err := s.repo.GetUser(ctx, email, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(ctx context.Context, token string) (primitive.ObjectID, error) {
	tokens, err := jwt.ParseWithClaims(token, &TokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	claims, ok := tokens.Claims.(*TokenClaim)
	if !ok {
		return primitive.NilObjectID, errors.New("token are not *TokenClaim types")
	}

	return claims.UserID, nil
}
