package api

import (
	"blog/internal/middleware"
	"blog/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (api *api) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["post_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commentID, err := api.db.CreateComment(comment.Content, userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment.ID = commentID
	comment.UserId = userID
	comment.PostId = postID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
