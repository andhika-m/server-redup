package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetHistory(c echo.Context) error {
	userID := c.Get("userID").(string)
	history, err := h.redupUsecase.GetHistory(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, history)
}

func (h *handler) AddHistory(c echo.Context) error {
	userID := c.Get("userID").(string)
	videoID := c.Param("id")
	err := h.redupUsecase.AddHistory(userID, videoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "History updated"})
}

func (h *handler) RemoveHistory(c echo.Context) error {
	userID := c.Get("userID").(string)
	historyID := c.Param("id")
	err := h.redupUsecase.RemoveHistory(userID, historyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "History removed"})
}
