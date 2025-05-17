package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/data"
	"github.com/mf751/gocha/internal/validator"
)

func (app *application) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatID  uuid.UUID `json:"chat_id"`
		Content string    `json:"content"`
		Type    int32     `json:"type"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	message := &data.Message{
		UserID: user.ID,
		ID:     uuid.New(),
		ChatID: input.ChatID,
		Content: data.Content{
			NullString: sql.NullString{
				Valid:  true,
				String: input.Content,
			},
		},
		Type: data.Int32{
			Int: sql.NullInt32{
				Valid: true,
				Int32: data.MessageNormal,
			},
		},
	}

	vdtr := validator.New()

	if data.ValidateMessage(vdtr, message, &app.models.Users); !vdtr.Valid() {
		app.failedValidationResponse(w, r, vdtr.Errors)
		return
	}

	err = app.models.Users.IsInChat(message.UserID, message.ChatID)
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

	err = app.models.Messages.SendMessage(message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var broadCastMessage NewMessageEvent
	broadCastMessage.ChatID = message.ChatID
	broadCastMessage.From = message.UserID
	broadCastMessage.Sent = message.Sent.Sent.Time
	broadCastMessage.Message = message.Content.NullString.String
	broadCastMessage.ID = message.ID

	sendData, err := json.Marshal(broadCastMessage)
	if err != nil {
		app.logger.PrintError(
			err,
			map[string]string{"error marshaling new chat message": err.Error()},
		)
		return
	}

	outGoingEvent := Event{
		Payload: sendData,
		Type:    EventNewMessage,
	}

	if _, ok := app.manager.clients[message.ChatID]; ok {
		for _, client := range app.manager.clients[message.ChatID] {
			client.egress <- outGoingEvent
		}
	}
}
