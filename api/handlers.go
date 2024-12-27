package users

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	// Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		sendJSON(
			w,
			Response{Message: "something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write json data", "error", err)
		return
	}
}

func HandleFindAll(ua *UsersApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := ua.FindAll()
		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func HandleFindById(ua *UsersApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		user := ua.FindById(ID(id))
		if user.ID == "" {
			sendJSON(
				w,
				Response{Message: "The user with the specified ID does not exist", Data: nil},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

func HandleInsert(ua *UsersApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			sendJSON(
				w,
				Response{Message: "There was an error while saving the user to the database"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if user.FirstName == "" ||
			user.LastName == "" ||
			user.Biography == "" ||
			len(user.FirstName) < 2 ||
			len(user.FirstName) > 20 ||
			len(user.LastName) < 2 ||
			len(user.LastName) > 20 ||
			len(user.Biography) < 20 ||
			len(user.Biography) > 450 {
			sendJSON(
				w,
				Response{Message: "Please provide FirstName LastName and bio for the user"},
				http.StatusBadRequest,
			)
			return
		}

		userWithID := ua.Insert(user)
		sendJSON(w, Response{Data: userWithID}, http.StatusCreated)
	}
}

func HandleUpdate(ua *UsersApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			sendJSON(
				w,
				Response{Message: "The user with the specified ID does not exist"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if user.FirstName == "" ||
			user.LastName == "" ||
			user.Biography == "" ||
			len(user.FirstName) < 2 ||
			len(user.FirstName) > 20 ||
			len(user.LastName) < 2 ||
			len(user.LastName) > 20 ||
			len(user.Biography) < 20 ||
			len(user.Biography) > 450 {
			sendJSON(
				w,
				Response{Message: "Please provide FirstName LastName and bio for the user"},
				http.StatusBadRequest,
			)
			return
		}

		userWithID, err := ua.Update(ID(id), user)
		if err != nil {
			sendJSON(
				w,
				Response{Message: "The user information could not be modified"},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(w, Response{Data: userWithID}, http.StatusOK)
	}
}

func HandleDelete(ua *UsersApplication) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		userWithID, err := ua.Delete(ID(id))
		if err != nil {
			sendJSON(
				w,
				Response{Message: "The user could not be removed"},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(w, Response{Data: userWithID}, http.StatusNoContent)
	}
}
