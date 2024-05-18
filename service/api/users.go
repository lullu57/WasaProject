package api

import (
	"encoding/json"
	"net/http"

	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func HandleAddUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user database.User
	db := ctx.Database

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx.Logger.Info("Adding user to the database")
	err := db.AddUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "username already exists") {
			http.Error(w, "Username already exists", http.StatusConflict) // Use HTTP 409 Conflict for username conflicts
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	responseMessage := "User created successfully."
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(responseMessage)); err != nil {
		ctx.Logger.Errorf("Failed to write response: %v", err)
	}
}

func HandleSetUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the new username
	var reqBody struct {
		NewUsername string `json:"newUsername"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.NewUsername == "" {
		http.Error(w, "New username must be provided", http.StatusBadRequest)
		return
	}

	ctx.Logger.Info("CurrentID: ", ctx.User.ID)

	currentUserID := ctx.User.ID // Ensure that ctx.User is populated correctly in the middleware

	ctx.Logger.Info("Setting new username for user ID: ", currentUserID)
	err = ctx.Database.SetUsername(currentUserID, reqBody.NewUsername)
	if err != nil {
		ctx.Logger.Error("Failed to update username: ", err)
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}

func HandleGetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username") // Assuming username is the URL parameter

	ctx.Logger.Info("Retrieving user profile for username: ", username)
	user, err := ctx.Database.GetUserProfile(username)
	if err != nil {
		ctx.Logger.Error("User not found: ", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func HandleGetUserProfileID(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := ps.ByName("userId") // Assuming userID is the URL parameter

	ctx.Logger.Info("Retrieving user profile for userID: ", userID)
	user, err := ctx.Database.GetUserProfileByID(userID)
	if err != nil {
		ctx.Logger.Error("User not found: ", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check if user exists
	user, err := ctx.Database.GetUserByUsername(req.Name)
	if err != nil {
		ctx.Logger.Error("Error retrieving user: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		// User does not exist, create new one
		user = &database.User{Username: req.Name}
		err = ctx.Database.AddUser(user) // Directly call AddUser now
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				http.Error(w, "Username already exists", http.StatusConflict)
				return
			}
			ctx.Logger.Error("Failed to create user: ", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Return the userId as a token
	response := struct {
		Token string `json:"token"`
	}{
		Token: user.ID, // Use userId as the token directly
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Error("Error encoding response: ", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func HandleFollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId") // Extracting username from URL parameter

	followerID := ctx.User.ID

	var err = ctx.Database.FollowUser(followerID, userId)
	if err != nil {
		ctx.Logger.Errorf("Error following user: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User %s followed %s", ctx.User.Username, userId)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User Followed successfully"})
}

func HandleUnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId") // Extracting username from URL parameter

	ctx.Logger.Infof("Unfollowing user: %s", userId)

	followerID := ctx.User.ID
	var err = ctx.Database.UnfollowUser(followerID, userId)
	if err != nil {
		ctx.Logger.Errorf("Error unfollowing user: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User %s unfollowed %s", ctx.User.Username, userId)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User Unfollowed successfully"})
}

// get all users
func HandleGetAllUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	users, err := ctx.Database.GetAllUsers()
	if err != nil {
		ctx.Logger.Errorf("Failed to get all users: %v", err)
		http.Error(w, "Failed to get all users", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("Fetched all users")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

/*
	func handleGetFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
		username := ps.ByName("username")
		if username == "" {
			http.Error(w, "Invalid username parameter", http.StatusBadRequest)
			return
		}

		followers, err := ctx.Database.GetFollowersByUsername(username)
		if err != nil {
			ctx.Logger.Error("Failed to retrieve followers: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		ctx.Logger.Infof("Followers fetched for user: %s", username)
		response := map[string][]string{"followers": followers}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
*/
func handleGetUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ctx.Logger.Infof("Fetching username for userId")
	userId := ps.ByName("userId")
	if userId == "" {
		http.Error(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}
	username, err := ctx.Database.GetUsername(userId)
	if err != nil {
		ctx.Logger.Error("Failed to retrieve username: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("Username fetched for userID: %s", userId)
	response := map[string]string{"username": username}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleIsUserFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	followerId := ctx.User.ID

	isFollowed, err := ctx.Database.IsUserFollowed(userId, followerId)
	if err != nil {
		ctx.Logger.Error("Failed to check if user is followed: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User follow status checked")
	response := map[string]bool{"isFollowed": isFollowed}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
