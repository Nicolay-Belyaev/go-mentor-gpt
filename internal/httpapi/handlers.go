package httpapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-mentor-gpt/domain"
	"log"
	"net/http"
	"strconv"
)

type StateRepository interface {
	domain.GoalRepository
	domain.TopicRepository
	domain.TaskRepository
}

type stateResponse struct {
	Goal   *domain.Goal   `json:"goal"`
	Topics []domain.Topic `json:"topics"`
	Tasks  []domain.Task  `json:"tasks"`
}

type updateRequest struct {
	Comment      string        `json:"comment"`
	Status       domain.Status `json:"status"`
	LastReviewed string        `json:"last_reviewed"`
}

type createResponse struct {
	ID int64 `json:"id"`
}

type Handlers struct {
	repo StateRepository
}

func NewHandlers(repo StateRepository) *Handlers {
	return &Handlers{repo: repo}
}

func (h *Handlers) GetState(w http.ResponseWriter, r *http.Request) {
	rContext := r.Context()

	goal, err := h.repo.GetGoal(rContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tasks, err := h.repo.GetAllTasks(rContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	topics, err := h.repo.GetAllTopics(rContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := stateResponse{
		Goal:   goal,
		Topics: topics,
		Tasks:  tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("encode state response: %v", err)
		return
	}
}
func (h *Handlers) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topic domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.CreateTopic(r.Context(), topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := createResponse{
		ID: id,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode topic response: %v", err)
		return
	}
}
func (h *Handlers) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topic := domain.Topic{
		ID:           id,
		Comment:      req.Comment,
		Status:       req.Status,
		LastReviewed: req.LastReviewed,
	}
	if err := h.repo.UpdateTopic(r.Context(), topic); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.CreateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := createResponse{
		ID: id,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode task response: %v", err)
		return
	}
}
func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := domain.Task{
		ID:           id,
		Comment:      req.Comment,
		Status:       req.Status,
		LastReviewed: req.LastReviewed,
	}
	if err := h.repo.UpdateTask(r.Context(), task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *Handlers) SetGoal(w http.ResponseWriter, r *http.Request) {
	var g domain.Goal
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.SetGoal(r.Context(), &g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
