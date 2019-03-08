package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string
	Password string
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := r.Context().Value(uid)
	ID, ok := id.(uuid.UUID)
	if !ok {
		h.responseErr(w, errors.New("Error getting context value"))
		return
	}

	user, err := h.UserStorage.ByID(ID)
	if err != nil {
		h.responseErr(w, err)
		return
	}
	spew.Dump(user)

	h.responseJSON(w, user)
	return
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	expireAuthCookie(w)
	h.responseJSON(w, "bye")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		h.responseErr(w, err)
		return
	}

	user, err := h.UserStorage.ByEmail(creds.Email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err == nil {
		token, err := h.Auth.Sign(user.ID)
		if err != nil {
			h.responseErr(w, err)
			return
		}
		setAuthCookie(w, token)
		h.responseJSON(w, token)
		return
	}

	h.responseErr(w, errors.New("Invalid credentials"))
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
		Expires:  time.Now().Add(24 * 30 * time.Hour),
	})
}

func expireAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
		Expires:  time.Now().Add(-1 * time.Second),
	})
}
