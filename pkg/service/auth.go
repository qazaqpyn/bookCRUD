package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
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
	tokenTTL   = 12 * time.Minute
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

func (s *AuthService) SignIn(ctx context.Context, inp model.LoginInput) (string, string, error) {
	password := generatePasswordHash(inp.Password)

	user, err := s.repo.Authorization.GetUser(ctx, inp.Email, password)
	if err != nil {
		return "", "", nil
	}

	return s.GenerateTokens(ctx, user.Id)
}

func (s *AuthService) GenerateTokens(ctx context.Context, userId primitive.ObjectID) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId.Hex(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
	})

	accessToken, err := t.SignedString([]byte(salt))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.repo.Tokens.Create(ctx, model.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s AuthService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.repo.Tokens.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("refresh token expired")
	}

	return s.GenerateTokens(ctx, session.UserID)
}

func (s *AuthService) ParseToken(ctx context.Context, token string) (primitive.ObjectID, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return primitive.NilObjectID, err
	}

	if !t.Valid {
		return primitive.NilObjectID, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return primitive.NilObjectID, errors.New("token are not *TokenClaim types")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid object")
	}

	id, err := primitive.ObjectIDFromHex(subject)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid subject")
	}
	return id, nil

}
