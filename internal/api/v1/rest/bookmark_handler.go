package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetBookmarks(c echo.Context) error {
	userID := c.Get("userID").(string)
	bookmarks, err := h.redupUsecase.GetBookmarks(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bookmarks)
}

func (h *handler) AddToBookmark(c echo.Context) error {
	userID := c.Get("userID").(string)
	videoID := c.Param("id")

	err := h.redupUsecase.AddBookmark(userID, videoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "video bookmarked"})
}

func (h *handler) RemoveBookmark(c echo.Context) error {
	userID := c.Get("userID").(string)
	bookmarkID := c.Param("id")

	err := h.redupUsecase.RemoveBookmark(userID, bookmarkID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if bookmarkID == "" {
		return c.JSON(http.StatusNoContent, map[string]string{"message": "no content"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "bookmark deleted"})
}
