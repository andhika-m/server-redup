package rest

import (
	"redup/internal/middleware"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, handler *handler) {
	authMiddleware := middleware.GetAuthMiddleware(handler.redupUsecase)

	// Video
	videoGroup := e.Group("/api/v1/videos")
	videoGroup.POST("", handler.CreateVideo, authMiddleware.CheckAuth)
	videoGroup.GET("", handler.GetVideos)
	videoGroup.GET("/:id", handler.GetVideoByID)
	videoGroup.PUT("/:id", handler.EditVideo, authMiddleware.CheckAuth)
	videoGroup.DELETE("/:id", handler.DeleteVideo, authMiddleware.CheckAuth)
	videoGroup.GET("/search", handler.SearchVideos)

	// Bookmark
	bookmarkGroup := e.Group("/api/v1/user/:userID/bookmarks")
	bookmarkGroup.GET("", handler.GetBookmarks, authMiddleware.CheckUsers)
	bookmarkGroup.POST("/:id", handler.AddToBookmark, authMiddleware.CheckUsers)
	bookmarkGroup.DELETE("/:id", handler.RemoveBookmark, authMiddleware.CheckUsers)

	// auth
	authGroup := e.Group("/api/v1")
	authGroup.POST("/register", handler.RegisterUser)
	authGroup.POST("/login", handler.Login)
	authGroup.DELETE("/logout/:userID", handler.Logout, authMiddleware.CheckUsers)

	// User
	userGroup := e.Group("/api/v1/user")
	userGroup.GET("/:userID", handler.GetUserByID, authMiddleware.CheckUsers)
	userGroup.PUT("/:userID", handler.UpdateUser, authMiddleware.CheckUsers)
	userGroup.DELETE("/:userID", handler.DeleteUser, authMiddleware.CheckUsers)

	// history
	historyGroup := e.Group("/api/v1/user/:userID/history")
	historyGroup.GET("", handler.GetHistory, authMiddleware.CheckUsers)
	historyGroup.POST("/:id", handler.AddHistory, authMiddleware.CheckUsers)
	historyGroup.DELETE("/:id", handler.RemoveHistory, authMiddleware.CheckUsers)
}
