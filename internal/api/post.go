package api

import (
	"blog/internal/middleware"
	"blog/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (api *api) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID, err := api.db.CreatePost(post, userID)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	post.ID = postID
	post.UserId = userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (api *api) GetPostById(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	postID, ok := vars["id"]
	if ok {
		id, err := strconv.Atoi(postID)
		if err != nil {
			http.Error(w, "Invalid post id", http.StatusBadRequest)
			return
		}
		data, err := api.db.GetPostByID(id, userID)
		if err != nil {
			http.Error(w, "Error getting post", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

		return
	}
}

func (api *api) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	data, err := api.db.GetAllPosts(userID)
	if err != nil {
		http.Error(w, "Error getting posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (api *api) DeletePostById(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	postID, ok := vars["id"]
	if ok {
		id, err := strconv.Atoi(postID)
		if err != nil {
			http.Error(w, "Invalid post id", http.StatusBadRequest)
			return
		}
		err = api.db.DeletePost(id, userID)
		if err != nil {
			http.Error(w, "Error deleting post", http.StatusInternalServerError)
			return
		}

		return
	}
}
