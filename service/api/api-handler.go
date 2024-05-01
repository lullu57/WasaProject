package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)
	rt.router.GET("/bans", rt.wrap(handleGetBannedUsers))
	rt.router.GET("/users/username/:username", rt.wrap(HandleGetUserProfile))
	rt.router.GET("/users/id/:userID", rt.wrap(HandleGetUserProfileID))
	rt.router.GET("/photos", rt.wrap(handleGetPhotos))
	rt.router.GET("/users", rt.wrap(HandleGetAllUsers))
	rt.router.GET("/photos/:photoId/comment/", rt.wrap(handleGetComments))
	rt.router.GET("/stream", rt.wrap(handleGetMyStream))
	rt.router.GET("/users/followers/:username", rt.wrap(handleGetFollowers))
	rt.router.GET("/photos/:photoId", rt.wrap(handleGetPhoto))
	rt.router.GET("/username/:userId", rt.wrap(handleGetUsername))
	rt.router.GET("/likes/:photoId", rt.wrap(HandleIsLiked))
	rt.router.GET("/follows/:userId", rt.wrap(handleIsUserFollowed))
	rt.router.GET("/bans/:userId", rt.wrap(handleIsUserBanned))
	rt.router.POST("/users", rt.wrap(HandleAddUser))
	rt.router.POST("/photos/:photoId/comments", rt.wrap(handleCommentPhoto))
	rt.router.POST("/photos/:photoId/likes", rt.wrap(HandleLikePhoto))
	rt.router.POST("/users/bans/:userId", rt.wrap(handleBanUser))
	rt.router.POST("/users/follows/:userId", rt.wrap(HandleFollowUser))
	rt.router.POST("/session", rt.wrap(doLogin))
	rt.router.POST("/photos", rt.wrap(handleUploadPhoto))
	rt.router.PATCH("/users/:username", rt.wrap(HandleSetUsername))
	rt.router.DELETE("/photos/:photoId/likes", rt.wrap(HandleUnlikePhoto))
	rt.router.DELETE("/photos", rt.wrap(handleDeletePhoto))
	rt.router.DELETE("/users/bans/:userId", rt.wrap(handleUnbanUser))
	rt.router.DELETE("/users/follows/:userId", rt.wrap(HandleUnfollowUser))
	rt.router.DELETE("/comments/:commentId", rt.wrap(handleUncommentPhoto))

	return rt.router
}
