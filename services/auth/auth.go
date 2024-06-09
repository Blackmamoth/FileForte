package authService

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Blackmamoth/fileforte/config"
	userModel "github.com/Blackmamoth/fileforte/models/user"
	"github.com/Blackmamoth/fileforte/types"
	"github.com/Blackmamoth/fileforte/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func AuthHandler() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/register", registerUser)
	router.Post("/login", loginUser)

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

func loginUser(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errors)
		return
	}

	user, err := userModel.GetUserByEmail(payload.Email)
	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusNotFound, fmt.Errorf("user with email [%s] does not exist", payload.Email))
		return
	}

	if !CompareHashedPassword(user.Password, payload.Password) {
		utils.SendAPIErrorResponse(w, http.StatusNotAcceptable, fmt.Errorf("invalid password, please try again with correct password"))
		return
	}

	accessToken, err := GenerateJWTToken(user.Id, r, Access, time.Now().Add(time.Minute*time.Duration(config.JWTConfig.JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS)))
	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	refreshTokenExpiration := time.Hour * 24 * time.Duration(config.JWTConfig.JWT_REFRESH_TOKEN_EXPIRATION_IN_DAYS)
	refreshToken, err := GenerateJWTToken(user.Id, r, Refresh, time.Now().Add(time.Hour*24*time.Duration(config.JWTConfig.JWT_REFRESH_TOKEN_EXPIRATION_IN_DAYS)))
	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{
		Name:     config.JWTConfig.REFRESH_TOKEN_COOKIE_NAME,
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(refreshTokenExpiration),
		Secure:   config.AppConfig.ENVIRONMENT != "DEVELOPMENT",
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   7,
		Path:     "/",
	}

	data := map[string]any{
		"userDetails": map[string]any{
			"userId":   user.Id,
			"username": user.UserName,
			"email":    user.Email,
		},
		"accessToken": accessToken,
	}

	utils.SendAPIResponse(w, http.StatusOK, data, &cookie)

}
