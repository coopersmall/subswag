package userauthentication

import (
	"context"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
)

const (
	MaxClockSkew = 5 * time.Minute
)

type AuthenticationService struct {
	logger          utils.ILogger
	tracer          apm.ITracer
	jwtService      jwtService
	rsaService      rsaService
	macPub          *rsa.PublicKey
	macPriv         *rsa.PrivateKey
	tokenSigningKey []byte
}

func NewAuthenticationService(
	logger utils.ILogger,
	tracer apm.ITracer,
	rsaService rsaService,
	jwtService jwtService,
	getMACPrivateKey func() (*rsa.PrivateKey, error),
	getMACPublicKey func() (*rsa.PublicKey, error),
	getTokenSigningKey func() ([]byte, error),
) *AuthenticationService {
	macPub, err := getMACPublicKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get MAC public key", err))
	}
	macPriv, err := getMACPrivateKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get MAC private key", err))
	}
	tokenSigningKey, err := getTokenSigningKey()
	if err != nil {
		panic(utils.NewInternalError("failed to get token signing key", err))
	}
	return &AuthenticationService{
		logger:          logger,
		tracer:          tracer,
		jwtService:      jwtService,
		rsaService:      rsaService,
		macPub:          macPub,
		macPriv:         macPriv,
		tokenSigningKey: tokenSigningKey,
	}
}

func (u *AuthenticationService) AuthenticateMAC(
	ctx context.Context,
	timestamp time.Time,
	signature string,
) error {
	var err error
	u.tracer.Trace(ctx, "authenticate-mac", func(ctx context.Context, span apm.ISpan) error {
		now := time.Now().UTC()
		if timestamp.Before(now.Add(-MaxClockSkew)) || timestamp.After(now.Add(MaxClockSkew)) {
			return errors.New("timestamp outside acceptable range")
		}

		var (
			signatureBytes []byte
			decrypted      []byte
		)
		signatureBytes, err = hex.DecodeString(signature)
		if err != nil {
			span.AddEvent("unable to decode signature")
			return utils.NewUnauthenticatedError("invalid signature")
		}

		timestampBytes := []byte(timestamp.Format(time.RFC3339))
		decrypted, err = u.rsaService.Decrypt(ctx, signatureBytes, u.macPriv)
		if err != nil {
			span.AddEvent("unable to decrypt MAC")
			return utils.NewInternalError("failed to decrypt signature", err)
		}

		if string(decrypted) != string(timestampBytes) {
			span.AddEvent("unable to authenticate MAC")
			return utils.NewUnauthenticatedError("invalid signature")
		}

		span.AddEvent("MAC authenticated")
		return nil
	})
	return err

}

func (u *AuthenticationService) AuthenticateToken(
	ctx context.Context,
	token string,
	now time.Time,
	permissions ...domain.Permission,
) (user.UserID, apitoken.APITokenID, error) {
	var (
		userId user.UserID
		rawId  utils.ID
		err    error
	)
	u.tracer.Trace(ctx, "authenticate-token", func(ctx context.Context, span apm.ISpan) error {
		split := strings.Split(token, "Bearer ")
		if len(split) != 2 {
			span.AddEvent("incorrect token format")
			return utils.NewUnauthenticatedError("invalid token")
		}
		userId, rawId, err = u.jwtService.ValidateToken(ctx, split[1], u.tokenSigningKey)
		if err != nil {
			span.AddEvent("invalid token")
			return utils.NewUnauthenticatedError("invalid token")
		}
		span.AddEvent("token authenticated")
		span.SetAttribute("user_id", userId)
		return nil
	})
	return userId, apitoken.APITokenID(rawId), err
}

func (u *AuthenticationService) AuthenticateSession(
	ctx context.Context,
	session string,
) (user.UserID, error) {
	var (
		converted int
		err       error
	)
	u.tracer.Trace(ctx, "authenticate-session", func(ctx context.Context, span apm.ISpan) error {
		var decrypted []byte
		decrypted, err = u.rsaService.Decrypt(ctx, []byte(session), u.macPriv)
		if err != nil {
			span.AddEvent("unable to decrypt session")
			return utils.NewInternalError("failed to decrypt session", err)
		}
		converted, err = strconv.Atoi(string(decrypted))
		if err != nil {
			span.AddEvent("invalid session")
			return utils.NewUnauthenticatedError("invalid session", err)
		}
		span.AddEvent("session authenticated")
		span.SetAttribute("user_id", converted)
		return nil
	})
	return user.UserID(converted), nil
}

type rsaService interface {
	Decrypt(ctx context.Context, data []byte, privateKey *rsa.PrivateKey) ([]byte, error)
}

type jwtService interface {
	ValidateToken(ctx context.Context, token string, signingKey []byte) (user.UserID, utils.ID, error)
}
