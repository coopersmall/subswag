package api

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/http/server"
	"github.com/coopersmall/subswag/utils"
)

type UsersHandler struct {
	server.IHandler
}

func NewUsersHandler(env env.IEnv) server.IHandler {
	resource := "/users"
	return &UsersHandler{
		IHandler: server.NewHandler(
			resource,
			env,
			[]domain.Permission{domain.APIPermission},
			[]server.Middleware{},
			server.APIGetRoute("", GetAllUsersRoute),
			server.APIGetRoute("/{userId}", GetUserRoute),
			server.APIPutRoute("", UpdateUserRoute),
		),
	}
}

func GetAllUsersRoute(r server.IRequest) (any, error) {
	users, err := r.GetServices().UsersService().GetAllUsers(r.Ctx())
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserRoute(r server.IRequest) (any, error) {
	userId, err := r.Param("userId")
	if err != nil {
		return nil, err
	}
	parsed, err := utils.ParseID(userId)
	if err != nil {
		return nil, err
	}
	user, err := r.GetServices().UsersService().GetUser(r.Ctx(), user.UserID(parsed))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserRoute(r server.IRequest) (any, error) {
	body, err := r.Body()
	if err != nil {
		return nil, err
	}
	var user *user.User
	err = utils.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}
	user, err = r.GetServices().UsersService().UpdateUser(r.Ctx(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
