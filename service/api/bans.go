package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func handleBanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	bannedBy := ctx.User.ID

	// Check if user is trying to ban themselves
	if userId == bannedBy {
		http.Error(w, "Cannot ban yourself", http.StatusBadRequest)
		return
	}

	// Check if the banning user is banned by the user they are trying to ban
	isBannedByUser, err := ctx.Database.IsBannedBy(userId, bannedBy)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if isBannedByUser {
		http.Error(w, "Cannot ban a user who has banned you", http.StatusForbidden)
		return
	}

	err = ctx.Database.BanUser(bannedBy, userId)
	if err != nil {
		if err.Error() == "user is already banned" {
			http.Error(w, "User is already banned", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User %s banned by %s", userId, ctx.User.Username)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("User successfully banned")); err != nil {
		ctx.Logger.Errorf("Failed to write response: %v", err)
	}
}

func handleUnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	if userId == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Check if user is trying to unban themselves
	if userId == ctx.User.ID {
		http.Error(w, "Cannot unban yourself", http.StatusBadRequest)
		return
	}

	bannerUser := ctx.User.ID

	err := ctx.Database.UnbanUser(bannerUser, userId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User %s unbanned by %s", userId, ctx.User.Username)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("User successfully unbanned")); err != nil {
		ctx.Logger.Errorf("Failed to write response: %v", err)
	}
}

func handleIsUserBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var banner = ctx.User.ID
	userId := ps.ByName("userId")
	if userId == "" {
		http.Error(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	banned, err := ctx.Database.BanExists(banner, userId)
	if err != nil {
		ctx.Logger.Error("Failed to check if user is banned: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{"banned": banned}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Errorf("Failed to write response: %v", err)
	}
}
