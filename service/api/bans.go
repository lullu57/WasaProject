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
	var err = ctx.Database.BanUser(bannedBy, userId)
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
	w.Write([]byte("User successfully banned"))
}

// Handler for unbanning a user
func handleUnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	ctx.Logger.Infof("Unbanning user %s", userId)
	if userId == "" {
		ctx.Logger.Infof("Invalid parameters")
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Fetch user IDs by username
	bannerUser := ctx.User.ID

	var err = ctx.Database.UnbanUser(bannerUser, userId)
	if err != nil {
		ctx.Logger.Infof("Internal server error " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User %s unbanned by %s", userId, ctx.User.Username)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully unbanned"))
}

// Handler for getting all banned users
func handleGetBannedUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	bannedUsers, err := ctx.Database.GetBans()
	if err != nil {
		ctx.Logger.Infof("Internal server error1 " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(bannedUsers)
	ctx.Logger.Infof("Banned users fetched")
	if err != nil {
		ctx.Logger.Infof("Internal server error " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
	json.NewEncoder(w).Encode(response)
}
