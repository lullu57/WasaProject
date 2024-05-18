package api

import (
	"encoding/json"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func handleCommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if ctx.User == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	photoId := ps.ByName("photoId")
	if photoId == "" {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	comment := database.Comment{
		ID:        uuid.Must(uuid.NewV4()).String(), // Using a UUID library to generate the comment ID
		UserID:    ctx.User.ID,
		PhotoID:   photoId,
		Content:   req.Content,
		Timestamp: time.Now(),
	}

	err := ctx.Database.AddComment(comment)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("Comment added by %s", ctx.User.Username)
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("Comment added successfully")); err != nil {
		ctx.Logger.Errorf("Failed to write response: %v", err)
	}
}

func handleUncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commentID := ps.ByName("commentId")
	if commentID == "" {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err := ctx.Database.DeleteComment(commentID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("Comment deleted by %s", ctx.User.Username)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Comment deleted successfully")); err != nil {
		ctx.Logger.Errorf("Failed to write response: %v", err)
	}
}

func handleGetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoId := ps.ByName("photoId")
	if photoId == "" {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	comments, err := ctx.Database.GetCommentsByPhotoId(photoId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("Comments fetched")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
