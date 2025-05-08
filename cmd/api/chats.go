package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/data"
	"github.com/mf751/gocha/internal/validator"
)

func (app *application) createChatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OwnerID   uuid.UUID `json:"-"`
		Name      string    `json:"name"`
		IsPrivate bool      `json:"is_private"`
		UserID    uuid.UUID `json:"user_id,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	requestUser := app.contextGetUser(r)

	chat := &data.Chat{
		OwnerID:   requestUser.ID,
		Name:      input.Name,
		IsPrivate: input.IsPrivate,
	}

	vdtr := validator.New()
	if data.ValidateChatName(vdtr, chat.Name); !vdtr.Valid() {
		app.failedValidationResponse(w, r, vdtr.Errors)
		return
	}
	if input.IsPrivate {
		user, err := app.models.Users.GetByID(input.UserID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				vdtr.AddError("user_id", "no user exists with this id")
				app.failedValidationResponse(w, r, vdtr.Errors)
				return
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}
		if requestUser.ID == user.ID {
			vdtr.AddError("user_id", "cannot create chat with self")
			app.failedValidationResponse(w, r, vdtr.Errors)
			return
		}
		chat.ID = input.UserID
	} else {
		chat.ID = uuid.New()
	}

	chat.CreatedAt, err = app.models.Chats.Insert(chat)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateChat):
			vdtr.AddError("chat", "duplicate private chat")
			app.failedValidationResponse(w, r, vdtr.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.models.Chats.Join(chat.ID, requestUser.ID, true)

	err = app.writeJSON(w, http.StatusOK, envelope{"chat": chat}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteChatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatId uuid.UUID `json:"chat_id"`
	}

	app.readJSON(w, r, &input)

	user := app.contextGetUser(r)

	err := app.models.Chats.Delete(user.ID, input.ChatId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDeletionFailed):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "Unable to delete chat")
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "deleted successfully!"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatId uuid.UUID `json:"chat_id"`
	}

	app.readJSON(w, r, &input)
	chatUsers, err := app.models.Chats.GetUsers(input.ChatId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrChatNotFound):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "Chat not found")
			return
		case errors.Is(err, data.ErrPrivateChat):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "Private chat")
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"users": chatUsers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) joinChatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatId uuid.UUID `json:"chat_id"`
	}

	app.readJSON(w, r, &input)
	user := app.contextGetUser(r)
	err := app.models.Chats.Join(input.ChatId, user.ID, false)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrChatNotFound):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "Chat not found")
			return
		case errors.Is(err, data.ErrPrivateChat):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "Private chat")
			return
		case errors.Is(err, data.ErrAlreadyMember):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		default:
			app.serverErrorResponse(w, r, err)
			return

		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "ok"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) leaveChatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatId uuid.UUID `json:"chat_id"`
	}

	app.readJSON(w, r, &input)
	user := app.contextGetUser(r)
	err := app.models.Chats.Leave(input.ChatId, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotInChat):
			app.errorResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "ok"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
