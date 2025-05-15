package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/mf751/gocha/internal/data"
	"github.com/mf751/gocha/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	vdtr := validator.New()

	if data.ValidateUser(vdtr, user); !vdtr.Valid() {
		app.failedValidationResponse(w, r, vdtr.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmial):
			vdtr.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, vdtr.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// --[ Activation Tokens ]--

	// token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }
	//
	// app.background(func() {
	// 	data := map[string]interface{}{
	// 		"activationToken": token.PlainText,
	// 		"userID":          user.ID,
	// 	}
	// 	err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
	// 	if err != nil {
	// 		app.logger.PrintError(err, nil)
	// 	}
	// })

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	vdtr := validator.New()

	data.ValidateEmail(vdtr, input.Email)
	data.ValidatePasswordPlainText(vdtr, input.Password)

	if !vdtr.Valid() {
		app.failedValidationResponse(w, r, vdtr.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeAuthentication, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 2*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusCreated,
		envelope{"authentication_token": token, "user": user},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	chats, err := app.models.Users.GetChats(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": chats}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
