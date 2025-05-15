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
	message := &data.Message{
		UserID: requestUser.ID,
		ChatID: chat.ID,
		ID:     uuid.New(),
		Content: data.Content{
			NullString: sql.NullString{
				Valid:  true,
				String: requestUser.Name + " Created the chat.",
			},
		},
		Type: data.Int32{
			Int: sql.NullInt32{
				Valid: true,
				Int32: data.MessageJoined,
			},
		},
	}
	err = app.models.Messages.SendMessage(message)
	if err != nil {
		app.logger.PrintError(
			err,
			map[string]string{"error sending created chat message": err.Error()},
		)
		return
	}

	if client, ok := app.manager.connectionClients[requestUser.ID]; ok {
		var broadMessage NewMessageEvent
		broadMessage.Message = message.Content.NullString.String
		broadMessage.From = message.UserID
		broadMessage.ChatID = message.ChatID
		broadMessage.Sent = message.Sent.Sent.Time

		sendData, err := json.Marshal(broadMessage)
		if err != nil {
			app.logger.PrintError(
				err,
				map[string]string{"error marshaling created chat message": err.Error()},
			)
			return
		}

		outGoingEvent := Event{
			Payload: sendData,
			Type:    EventJoinedMessage,
		}

		app.manager.Lock()
		app.manager.clients[message.ChatID] = make(map[uuid.UUID]*Client)
		app.manager.clients[message.ChatID][message.UserID] = client
		app.manager.clients[message.ChatID][message.UserID].chatsID = append(
			app.manager.clients[message.ChatID][message.UserID].chatsID,
			message.ChatID,
		)
		app.manager.Unlock()
		client.egress <- outGoingEvent
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
	if _, ok := app.manager.clients[input.ChatId]; ok {
		app.manager.Lock()
		for userID := range app.manager.clients[input.ChatId] {
			app.manager.connectionClients[userID].chatsID = removeFromSliceByValue(
				app.manager.connectionClients[userID].chatsID,
				input.ChatId,
			)
		}
		delete(app.manager.clients, input.ChatId)
		app.manager.Unlock()
	}
}

func (app *application) getChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatId uuid.UUID `json:"chat_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
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

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	err = app.models.Chats.Join(input.ChatId, user.ID, false)
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
		return
	}

	message := &data.Message{
		UserID: user.ID,
		ChatID: input.ChatId,
		ID:     uuid.New(),
		Content: data.Content{
			NullString: sql.NullString{
				Valid:  true,
				String: user.Name + " Joined the chat.",
			},
		},
		Type: data.Int32{
			Int: sql.NullInt32{
				Valid: true,
				Int32: data.MessageJoined,
			},
		},
	}
	err = app.models.Messages.SendMessage(message)
	if err != nil {
		app.logger.PrintError(
			err,
			map[string]string{"error sending joined chat message": err.Error()},
		)
		return
	}

	var broadMessage NewMessageEvent
	broadMessage.Message = message.Content.NullString.String
	broadMessage.From = message.UserID
	broadMessage.ChatID = message.ChatID
	broadMessage.Sent = message.Sent.Sent.Time

	sendData, err := json.Marshal(broadMessage)
	if err != nil {
		app.logger.PrintError(
			err,
			map[string]string{"error marshaling joined chat message": err.Error()},
		)
		return
	}

	outGoingEvent := Event{
		Payload: sendData,
		Type:    EventJoinedMessage,
	}

	app.manager.Lock()
	if _, ok := app.manager.clients[message.ChatID]; !ok {
		app.manager.clients[message.ChatID] = make(map[uuid.UUID]*Client)
	}
	if _, ok := app.manager.connectionClients[message.UserID]; ok {
		app.manager.clients[message.ChatID][message.UserID] = app.manager.connectionClients[message.UserID]
		app.manager.clients[message.ChatID][message.UserID].chatsID = append(
			app.manager.clients[message.ChatID][message.UserID].chatsID,
			message.ChatID,
		)
	}
	app.manager.Unlock()

	for _, client := range app.manager.clients[message.ChatID] {
		client.egress <- outGoingEvent
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

	message := &data.Message{
		UserID: user.ID,
		ChatID: input.ChatId,
		ID:     uuid.New(),
		Content: data.Content{
			NullString: sql.NullString{
				Valid:  true,
				String: user.Name + " Left the chat.",
			},
		},
		Type: data.Int32{
			Int: sql.NullInt32{
				Valid: true,
				Int32: data.MessageLeft,
			},
		},
	}
	err = app.models.Messages.SendMessage(message)
	if err != nil {
		app.logger.PrintError(
			err,
			map[string]string{"error sending Left chat message": err.Error()},
		)
		return
	}

	var broadMessage NewMessageEvent
	broadMessage.Message = message.Content.NullString.String
	broadMessage.From = message.UserID
	broadMessage.ChatID = message.ChatID
	broadMessage.Sent = message.Sent.Sent.Time

	sendData, err := json.Marshal(broadMessage)
	if err != nil {
		app.logger.PrintError(
			err,
			map[string]string{"error marshaling left chat message": err.Error()},
		)
		return
	}

	outGoingEvent := Event{
		Payload: sendData,
		Type:    EventLeftMessage,
	}

	if _, ok := app.manager.clients[message.ChatID][message.UserID]; ok {
		app.manager.Lock()
		app.manager.connectionClients[message.UserID].chatsID = removeFromSliceByValue(
			app.manager.connectionClients[message.UserID].chatsID,
			message.ChatID,
		)
		delete(app.manager.clients[message.ChatID], user.ID)
		app.manager.Unlock()
	}

	for _, client := range app.manager.clients[message.ChatID] {
		client.egress <- outGoingEvent
	}
}

func (app *application) getChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatId uuid.UUID `json:"chat_id"`
		Size   int       `json:"size"`
		Start  int       `json:"start"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	vdtr := validator.New()
	vdtr.Check(input.Size != 0, "size", "must be provided and more than 0")
	if !vdtr.Valid() {
		app.failedValidationResponse(w, r, vdtr.Errors)
		return
	}

	messages, err := app.models.Chats.GetChatMessage(input.ChatId, input.Size, input.Start)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": messages}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
