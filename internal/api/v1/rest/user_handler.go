package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redup/internal/model"

	"github.com/labstack/echo/v4"
)

func (h *handler) RegisterUser(c echo.Context) error {
	var request model.RegisterRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	userData, err := h.redupUsecase.RegisterUser(request)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": userData,
	})
}

func (h *handler) Login(c echo.Context) error {
	var request model.LoginRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	sessionData, err := h.redupUsecase.Login(request)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": sessionData,
	})
}

func (h *handler) GetUserByID(c echo.Context) error {
	userID := c.Param("userID")
	user, err := h.redupUsecase.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (h *handler) UpdateUser(c echo.Context) error {
	loggedInUserID := c.Get("userID").(string)
	userID := c.Param("userID")

	if loggedInUserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You can only edit your own profile"})
	}

	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user.ID = userID
	updatedUser, err := h.redupUsecase.UpdateUser(user.ID, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func (h *handler) Logout(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	user := c.Param("userID")

	if userID != user {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You can't access this user"})
	}
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	token, ok := c.Get("token").(string)
	if !ok || token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	err := h.redupUsecase.Logout(userID, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}

func (h *handler) DeleteUser(c echo.Context) error {
	authenticatedUserID := c.Get("userID").(string)
	userID := c.Param("userID")

	if authenticatedUserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You can only delete your own account"})
	}

	err := h.redupUsecase.DeleteUser(userID)
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}
