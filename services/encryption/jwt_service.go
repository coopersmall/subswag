package encryption

import (
	"context"
	"errors"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
	"github.com/golang-jwt/jwt"
)

type JWTService struct {
	logger    utils.ILogger
	makeToken func(user.UserID, utils.ID, ...domain.Permission) *jwt.Token
}

func NewJWTService(logger utils.ILogger) *JWTService {
	signingMethod := jwt.GetSigningMethod(algorithm)
	makeToken := func(userId user.UserID, tokenId utils.ID, permissions ...domain.Permission) *jwt.Token {
		return jwt.NewWithClaims(signingMethod, jwt.MapClaims{
			userClaim:  userId,
			tokenClaim: tokenId,
		})
	}

	return &JWTService{
		logger:    logger,
		makeToken: makeToken,
	}
}

func (j *JWTService) CreateToken(
	ctx context.Context,
	userId user.UserID,
	signingKey []byte,
) (string, error) {
	return j.createToken(
		ctx,
		userId,
		utils.NewID(),
		signingKey,
	)
}

func (j *JWTService) CreateTokenWithID(
	ctx context.Context,
	userId user.UserID,
	tokenId utils.ID,
	signingKey []byte,
) (string, error) {
	return j.createToken(
		ctx,
		userId,
		tokenId,
		signingKey,
	)
}

func (j *JWTService) createToken(
	ctx context.Context,
	userId user.UserID,
	tokenId utils.ID,
	signingKey []byte,
) (string, error) {
	token := j.makeToken(
		userId,
		tokenId,
	)

	signed, err := token.SignedString(signingKey)
	if err != nil {
		j.logger.Error(ctx, "error signing token", err, nil)
		return "", err
	}

	j.logger.Debug(ctx, "token created", nil)
	return signed, nil
}

func (j *JWTService) ValidateToken(ctx context.Context, token string, signingKey []byte) (
	user.UserID,
	utils.ID,
	error,
) {
	return j.validateToken(ctx, token, signingKey)
}

func (j *JWTService) validateToken(ctx context.Context, token string, signingKey []byte) (user.UserID, utils.ID, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return signingKey, nil
	})
	if err != nil {
		return 0, 0, errors.New("error parsing token")
	}
	if err := t.Claims.(jwt.MapClaims).Valid(); err != nil {
		return 0, 0, err
	}
	userId, ok := t.Claims.(jwt.MapClaims)[userClaim]
	if !ok {
		return 0, 0, errors.New("user id not found in token")
	}

	tokenId, ok := t.Claims.(jwt.MapClaims)[tokenClaim]
	if !ok {
		return 0, 0, errors.New("token id not found in token")
	}
	return user.UserID(userId.(float64)), utils.ID(tokenId.(float64)), nil
}

const (
	algorithm  = "HS256"
	userClaim  = "uid"
	tokenClaim = "tid"
)
