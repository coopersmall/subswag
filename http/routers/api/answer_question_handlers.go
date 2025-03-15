package api

import (
	"encoding/json"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/chatsession"
	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/http/server"
)

type AnswerQuestionHandler struct {
	server.IHandler
}

func NewAnswerQuestionHandler(env env.IEnv) server.IHandler {
	resource := "/answer-question"
	return &AnswerQuestionHandler{
		IHandler: server.NewHandler(
			resource,
			env,
			[]domain.Permission{domain.APIPermission},
			[]server.Middleware{},
			server.APIPostRoute("", AnswerQuestionRoute),
		),
	}
}

func AnswerQuestionRoute(r server.IRequest) (any, error) {
	body, err := r.Body()
	if err != nil {
		return nil, err
	}
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, err
	}
	item := chatsession.NewUserChatSessionItem(chatsession.NewChatSessionID(), req.Content)
	service := r.GetServices().AnswerAssistantService(r.UserID())
	response, err := service.AnswerQuestion(
		r.Ctx(),
		r.UserID(),
		[]chatsession.ChatSessionItem{item},
	)
	return response, err
}

type Request struct {
	Content string `json:"content"`
}
