package auth

import (
	"fmt"
	"net/http"

	userModel "github.com/Blackmamoth/fileforte/models/user"
	"github.com/Blackmamoth/fileforte/types"
	"github.com/Blackmamoth/fileforte/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func AuthHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/register", registerUser)
	return router
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errors)
		return
	}

	_, err := userModel.GetUserByEmail(payload.Email)
	if err == nil {
		utils.SendAPIErrorResponse(w, http.StatusConflict, fmt.Errorf("user with email [%s] already exists", payload.Email))
		return
	}

	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	result, err := userModel.CreateUser(types.User{
		UserName: payload.UserName,
		Email:    payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	insertId, err := result.LastInsertId()

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	userDetails := map[string]any{
		"userDetails": map[string]any{
			"userId":   insertId,
			"username": payload.UserName,
			"email":    payload.Email,
		},
	}

	utils.SendAPIResponse(w, http.StatusCreated, userDetails, nil)
}
