package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type createNoteRequest struct {
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// notes — входная точка для:
// POST /notes
// GET  /notes/{id}?user_id=1
// PUT  /notes/{id}
// DELETE /notes/{id}?user_id=1
func (h *Handler) notes(w http.ResponseWriter, r *http.Request) {
	// POST /notes
	if r.URL.Path == "/notes" {
		if r.Method == http.MethodPost {
			h.createNote(w, r)
			return
		}
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// /notes/{id}
	if strings.HasPrefix(r.URL.Path, "/notes/") {
		switch r.Method {
		case http.MethodGet:
			h.getNote(w, r)
			return
		case http.MethodPut:
			h.updateNote(w, r)
			return
		case http.MethodDelete:
			h.deleteNote(w, r)
			return
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
	}

	http.NotFound(w, r)
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	var req createNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	if req.UserID <= 0 || strings.TrimSpace(req.Title) == "" {
		http.Error(w, `{"error":"user_id and title required"}`, http.StatusBadRequest)
		return
	}

	note, err := h.db.CreateNote(r.Context(), req.UserID, req.Title, req.Content)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(note)
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	noteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || noteID <= 0 {
		http.Error(w, `{"error":"bad note id"}`, http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, `{"error":"user_id query required"}`, http.StatusBadRequest)
		return
	}

	note, err := h.db.GetNoteByID(
		r.Context(),
		userID,
		noteID,
	)
	if err != nil {
		// если у тебя ErrNotFound — можно сравнить и вернуть 404
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	noteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || noteID <= 0 {
		http.Error(w, `{"error":"bad note id"}`, http.StatusBadRequest)
		return
	}

	var req updateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
		return
	}
	if req.UserID <= 0 || strings.TrimSpace(req.Title) == "" {
		http.Error(w, `{"error":"user_id and title required"}`, http.StatusBadRequest)
		return
	}

	note, err := h.db.CreateNote(r.Context(), req.UserID, req.Title, req.Content)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	noteID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || noteID <= 0 {
		http.Error(w, `{"error":"bad note id"}`, http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, `{"error":"user_id query required"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteNote(
		r.Context(),
		userID,
		noteID,
	); err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
