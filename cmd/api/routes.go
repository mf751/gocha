package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	router.HandlerFunc(
		http.MethodPost,
		"/v1/users",
		app.registerUserHandler,
	)
	router.HandlerFunc(
		http.MethodPost,
		"/v1/tokens/authentication",
		app.createAuthenticationTokenHandler,
	)

	router.HandlerFunc(
		http.MethodPost,
		"/v1/chats",
		app.requireAuthentication(app.createChatHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/v1/chats",
		app.requireAuthentication(app.deleteChatHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/v1/chat/users",
		app.requireAuthentication(app.getChatUsersHandler),
	)
	router.HandlerFunc(
		http.MethodPost,
		"/v1/chat/join",
		app.requireAuthentication(app.joinChatHandler),
	)
	router.HandlerFunc(
		http.MethodPost,
		"/v1/chat/leave",
		app.requireAuthentication(app.leaveChatHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/v1/chats",
		app.requireAuthentication(app.getUserChatsHandler),
	)
	router.HandlerFunc(http.MethodGet, "/v1/ws", app.requireAuthentication(app.manager.serveWS))

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
