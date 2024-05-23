package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.wrap(doLogin))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// User routes
	rt.router.GET("/users", rt.wrap(HandleGetAllUsers))
	rt.router.GET("/users/:userId/username", rt.wrap(handleGetUsername))
	rt.router.GET("/users/:userId", rt.wrap(HandleGetUserProfileID))
	rt.router.POST("/users", rt.wrap(HandleAddUser))
	rt.router.PATCH("/users/username", rt.wrap(HandleSetUsername))

	// Photo routes
	rt.router.GET("/photos", rt.wrap(handleGetPhotos))
	rt.router.GET("/photos/:photoId", rt.wrap(handleGetPhoto))
	rt.router.POST("/photos", rt.wrap(handleUploadPhoto))
	rt.router.DELETE("/photos/:photoId", rt.wrap(handleDeletePhoto))
	rt.router.GET("/stream", rt.wrap(handleGetMyStream))

	// likes routes
	rt.router.GET("/photos/:photoId/likes", rt.wrap(HandleIsLiked))
	rt.router.POST("/photos/:photoId/likes", rt.wrap(HandleLikePhoto))
	rt.router.DELETE("/photos/:photoId/likes", rt.wrap(HandleUnlikePhoto))

	// Comments routes
	rt.router.POST("/photos/:photoId/comments", rt.wrap(handleCommentPhoto))
	rt.router.GET("/photos/:photoId/comments", rt.wrap(handleGetComments))
	rt.router.DELETE("/comments/:commentId", rt.wrap(handleUncommentPhoto))

	// follow routes
	rt.router.GET("/follows/:userId", rt.wrap(handleIsUserFollowed))
	rt.router.DELETE("/users/:userId/follows", rt.wrap(HandleUnfollowUser))
	rt.router.POST("/users/:userId/follows", rt.wrap(HandleFollowUser))

	// ban routes
	rt.router.GET("/bans/:userId", rt.wrap(handleIsUserBanned))
	rt.router.DELETE("/users/:userId/bans", rt.wrap(handleUnbanUser))
	rt.router.POST("/users/:userId/bans", rt.wrap(handleBanUser))

	return rt.router
}
