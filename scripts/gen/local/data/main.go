package main

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/apitoken"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/utils"
	"github.com/coopersmall/subswag/utils/actions"

	"github.com/joho/godotenv"
)

const SCHEMA_FILE = "db/sql/schema.sql"

func main() {
	godotenv.Load()
	logger := utils.GetLogger("load_local_data")

	schema, err := actions.GetSchema()
	if err != nil {
		panic(err)
	}

	envVars := env.MustGetEnvVars()
	env := env.MustGetEnv(envVars)

	manager := env.GetDB()
	wrDB := manager.ReadWrite()
	_, err = wrDB.Exec(schema)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	userId := user.UserID(12345)
	services, close := env.GetServices()
	defer close()

	data := user.UserData{
		Email:     "fcoopersmall@gmail.com",
		FirstName: "Cooper",
		LastName:  "Small",
	}

	_, err = services.UsersService().CreateUserWithId(ctx, userId, data)
	if err != nil {
		panic(err)
	}

	logger.Info(ctx, "User created", map[string]any{"user_id": userId})

	tokenId := utils.ID(12345)
	expiry := time.Now().Add(time.Hour * 24 * 7)

	token, err := services.APITokenService(userId).CreateTokenWithId(
		ctx,
		tokenId,
		&apitoken.APITokenData{
			UserId: userId,
			Expiry: expiry,
			Permissions: []domain.Permission{
				domain.APIPermission,
			},
		},
	)

	if err != nil {
		panic(err)
	}

	logger.Info(ctx, "User API Token", map[string]any{
		"expiry":  expiry,
		"user_id": userId,
		"token":   token.APIToken,
		"secret":  token.Secret,
	})

	return
}
