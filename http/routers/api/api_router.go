package api

import (
	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/http/server"
)

const PATH = "/api/"

func NewAPIRouter(env env.IEnv) server.IRouter {
	return server.NewRouter(
		PATH,
		env.GetLogger("api-router"),
		env.GetTracer("api-router"),
		[]server.Middleware{},
		NewChatSessionsHandler(env),
		NewUsersHandler(env),
		NewAnswerQuestionHandler(env),
	)
}
