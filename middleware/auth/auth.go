package authMiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Blackmamoth/fileforte/config"
	userModel "github.com/Blackmamoth/fileforte/models/user"
	authService "github.com/Blackmamoth/fileforte/services/auth"
	"github.com/Blackmamoth/fileforte/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userId"

func WithJWTToken(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(config.JWTConfig.ACCESS_TOKEN_HEADER_NAME)
		if accessToken == "" {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("access token not present in the request"))
			return
		}

		jwtAccessToken, err := authService.ValidateJWTTOken(accessToken, authService.Access)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		if !jwtAccessToken.Valid {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid access token"))
			return
		}

		claims := jwtAccessToken.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)
		userId, err := strconv.Atoi(str)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("failed to convert userId to int: %v", err))
			return
		}

		user, err := userModel.GetUserById(userId)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.Id)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}
