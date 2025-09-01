package api

import (
	"blog/internal/middleware"
	"blog/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	// Проверяем rate limit для комментариев
	ctx := r.Context()
	err = api.commentLimiter.CheckCommentLimit(ctx, userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ВСЕГДА записываем попытку комментирования
	err = api.commentLimiter.RecordCommentAttempt(ctx, userID, postID)
	if err != nil {
		// Логируем ошибку, но продолжаем выполнение
		// log.Printf("Failed to record comment attempt: %v", err)
	}

	commentID, err := api.db.CreateComment(comment.Content, userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// НЕ сбрасываем счетчик при успехе - он должен накапливаться!
	// _ = api.commentLimiter.ResetCommentAttempts(ctx, userID, postID)

	comment.ID = commentID
	comment.UserId = userID
	comment.PostId = postID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
