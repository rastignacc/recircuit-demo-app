package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func ErrBadRequest(msg string) *AppError {
	return NewAppError(http.StatusBadRequest, msg)
}

func ErrNotFound(msg string) *AppError {
	return NewAppError(http.StatusNotFound, msg)
}

func ErrUnauthorized(msg string) *AppError {
	return NewAppError(http.StatusUnauthorized, msg)
}

func ErrForbidden(msg string) *AppError {
	return NewAppError(http.StatusForbidden, msg)
}

func ErrConflict(msg string) *AppError {
	return NewAppError(http.StatusConflict, msg)
}

func ErrInternal(msg string) *AppError {
	return NewAppError(http.StatusInternalServerError, msg)
}

func WriteError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		WriteJSON(w, appErr.Code, appErr)
		return
	}
	WriteJSON(w, http.StatusInternalServerError, &AppError{Message: "internal server error"})
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Fprintf(w, `{"error":"failed to encode response"}`)
	}
}
