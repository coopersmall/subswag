package api

import (
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/http/server"
	"github.com/coopersmall/subswag/utils"
)

type ChatSessionsHandler struct {
	server.IHandler
}

func NewChatSessionsHandler(env env.IEnv) *ChatSessionsHandler {
	return &ChatSessionsHandler{
		IHandler: server.NewHandler(
			"/chatsessions",
			env,
			[]domain.Permission{domain.APIPermission},
			[]server.Middleware{},
			server.APIGetRoute("", GetAllChatSessionsRoute),
			server.APIGetRoute("/{sessionId}", GetChatSessionRoute),
			server.APIPostRoute("", CreateChatSessionRoute),
			server.APIPutRoute("/{sessionId}", UpdateChatSessionRoute),
			server.APIDeleteRoute("/{sessionId}", DeleteChatSessionRoute),
		),
	}
}

func GetAllChatSessionsRoute(r server.IRequest) (any, error) {
	return r.GetServices().ChatSessionsService(r.UserID()).GetAllChatSessions(r.Ctx())

}

func GetChatSessionRoute(r server.IRequest) (any, error) {
	sessionId, err := r.Param("sessionId")
	if err != nil {
		return nil, err
	}
	parsed, err := utils.ParseID(sessionId)
	if err != nil {
		return nil, err
	}
	return r.GetServices().ChatSessionsService(r.UserID()).GetChatSession(r.Ctx(), chatsession.ChatSessionID(parsed))
}

func CreateChatSessionRoute(r server.IRequest) (any, error) {
	body, err := r.Body()
	if err != nil {
		return nil, err
	}
	var data chatsession.ChatSessionData
	err = utils.Unmarshal(body, &data)
	return nil, r.GetServices().ChatSessionsService(r.UserID()).CreateChatSession(r.Ctx(), data)
}

func UpdateChatSessionRoute(r server.IRequest) (any, error) {
	body, err := r.Body()
	if err != nil {
		return nil, err
	}
	var data chatsession.ChatSession
	err = utils.Unmarshal(body, &data)
	return nil, r.GetServices().ChatSessionsService(r.UserID()).UpdateChatSession(r.Ctx(), &data)
}

func DeleteChatSessionRoute(r server.IRequest) (any, error) {
	sessionId, err := r.Param("sessionId")
	if err != nil {
		return nil, err
	}
	parsed, err := utils.ParseID(sessionId)
	return nil, r.GetServices().ChatSessionsService(r.UserID()).DeleteChatSession(r.Ctx(), chatsession.ChatSessionID(parsed))
}
