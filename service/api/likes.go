package api

import (
	"fmt"
	"net/http"

	"encoding/json"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// HandleLikePhoto processes the request to like a photo
func HandleLikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoID := ps.ByName("photoId") // Assuming you're using httprouter and path parameter named "photoId"
	userID := ctx.User.ID           // Assuming `ctx` has a User object with ID field

	// Log the action
	ctx.Logger.Info("Liking photo", "userID", userID, "photoID", photoID)

	// Call LikePhoto method of the database object
	err := ctx.Database.LikePhoto(userID, photoID)
	if err != nil {
		ctx.Logger.Error("Error liking photo", "error", err)
		http.Error(w, "Failed to like photo", http.StatusInternalServerError)
		return
	}

	// Successfully liked the photo
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Photo liked successfully")
}

// HandleUnlikePhoto processes the request to unlike a photo
func HandleUnlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Use a map to hold the JSON payload for simplicity
	var data map[string]string

	// Decode the JSON body into the map
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID := ctx.User.ID
	photoID, ok := data["photoId"]
	if !ok {
		http.Error(w, "Photo ID is required", http.StatusBadRequest)
		return
	}

	// Log the action
	ctx.Logger.Info("Unliking photo", "userID", userID, "photoID", photoID)

	// Call UnlikePhoto method of the database object
	err := ctx.Database.UnlikePhoto(userID, photoID)
	if err != nil {
		ctx.Logger.Error("Error unliking photo", "error", err)
		http.Error(w, "Failed to unlike photo", http.StatusInternalServerError)
		return
	}

	// Successfully unliked the photo
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Photo unliked successfully")
}

func HandleIsLiked(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoID := ps.ByName("photoId")
	userID := ctx.User.ID

	ctx.Logger.Info("Checking if photo is liked", "userID", userID, "photoID", photoID)

	liked, err := ctx.Database.IsLiked(userID, photoID)
	if err != nil {
		ctx.Logger.Error("Error checking if photo is liked", "error", err)
		http.Error(w, "Failed to check if photo is liked", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"liked": liked}); err != nil {
		ctx.Logger.Errorf("Failed to write response: %v", err)
	}
}
