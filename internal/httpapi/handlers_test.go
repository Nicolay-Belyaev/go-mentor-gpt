package httpapi

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-mentor-gpt/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeRepo struct {
	goal   *domain.Goal
	topics []domain.Topic
	tasks  []domain.Task

	createTopicID int64
	createTaskID  int64

	getGoalErr     error
	getTopicsErr   error
	getTasksErr    error
	createTopicErr error
	updateTopicErr error
	createTaskErr  error
	updateTaskErr  error
	setGoalErr     error

	createdTopic *domain.Topic
	updatedTopic *domain.Topic
	createdTask  *domain.Task
	updatedTask  *domain.Task
	setGoal      *domain.Goal
}

func (f *fakeRepo) GetGoal(ctx context.Context) (*domain.Goal, error) {
	if f.getGoalErr != nil {
		return nil, f.getGoalErr
	}
	return f.goal, nil
}

func (f *fakeRepo) SetGoal(ctx context.Context, goal *domain.Goal) error {
	if f.setGoalErr != nil {
		return f.setGoalErr
	}
	f.setGoal = goal
	return nil
}

func (f *fakeRepo) GetAllTopics(ctx context.Context) ([]domain.Topic, error) {
	if f.getTopicsErr != nil {
		return nil, f.getTopicsErr
	}
	return f.topics, nil
}

func (f *fakeRepo) CreateTopic(ctx context.Context, topic domain.Topic) (int64, error) {
	if f.createTopicErr != nil {
		return 0, f.createTopicErr
	}
	f.createdTopic = &topic
	return f.createTopicID, nil
}

func (f *fakeRepo) UpdateTopic(ctx context.Context, topic domain.Topic) error {
	if f.updateTopicErr != nil {
		return f.updateTopicErr
	}
	f.updatedTopic = &topic
	return nil
}

func (f *fakeRepo) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	if f.getTasksErr != nil {
		return nil, f.getTasksErr
	}
	return f.tasks, nil
}

func (f *fakeRepo) CreateTask(ctx context.Context, task domain.Task) (int64, error) {
	if f.createTaskErr != nil {
		return 0, f.createTaskErr
	}
	f.createdTask = &task
	return f.createTaskID, nil
}

func (f *fakeRepo) UpdateTask(ctx context.Context, task domain.Task) error {
	if f.updateTaskErr != nil {
		return f.updateTaskErr
	}
	f.updatedTask = &task
	return nil
}

func TestGetState(t *testing.T) {
	repo := &fakeRepo{
		goal: &domain.Goal{
			ID:      1,
			Text:    "Senior Go",
			Comment: "System design",
		},
		topics: []domain.Topic{
			{
				ID:      1,
				Key:     "go.channels",
				Title:   "Channels",
				Status:  domain.REVIEWED,
				Comment: "Повторять cancellation",
			},
		},
		tasks: []domain.Task{
			{
				ID:      1,
				Title:   "Worker pool",
				Status:  domain.LEARNED,
				Comment: "База понятна",
			},
		},
	}

	h := NewHandlers(repo)

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	rec := httptest.NewRecorder()

	h.GetState(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	var resp stateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Goal.Text != "Senior Go" {
		t.Fatalf("unexpected goal: %+v", resp.Goal)
	}

	if len(resp.Topics) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(resp.Topics))
	}

	if len(resp.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(resp.Tasks))
	}
}

func TestGetStateRepositoryError(t *testing.T) {
	repo := &fakeRepo{
		getGoalErr: errors.New("boom"),
	}

	h := NewHandlers(repo)

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	rec := httptest.NewRecorder()

	h.GetState(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestCreateTopic(t *testing.T) {
	repo := &fakeRepo{
		createTopicID: 42,
	}

	h := NewHandlers(repo)

	body := []byte(`{
		"key":"go.channels",
		"title":"Channels",
		"comment":"test",
		"status":"familiar",
		"last_reviewed":"2026-08-22"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/topics", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateTopic(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
	}

	if repo.createdTopic == nil {
		t.Fatal("topic was not created")
	}

	if repo.createdTopic.Key != "go.channels" {
		t.Fatalf("unexpected topic key: %s", repo.createdTopic.Key)
	}

	var resp createResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.ID != 42 {
		t.Fatalf("expected id 42, got %d", resp.ID)
	}
}

func TestCreateTopicBadJSON(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPost,
		"/topics",
		bytes.NewBufferString(`{broken`),
	)
	rec := httptest.NewRecorder()

	h.CreateTopic(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateTopicRepositoryError(t *testing.T) {
	repo := &fakeRepo{
		createTopicErr: errors.New("boom"),
	}

	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPost,
		"/topics",
		bytes.NewBufferString(`{
			"key":"go.channels",
			"title":"Channels",
			"status":"familiar"
		}`),
	)
	rec := httptest.NewRecorder()

	h.CreateTopic(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestUpdateTopic(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandlers(repo)

	body := []byte(`{
		"comment":"updated",
		"status":"review",
		"last_reviewed":"2026-08-22"
	}`)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/topics/7",
		bytes.NewReader(body),
	)
	req.SetPathValue("id", "7")

	rec := httptest.NewRecorder()

	h.UpdateTopic(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rec.Code)
	}

	if repo.updatedTopic == nil {
		t.Fatal("topic was not updated")
	}

	if repo.updatedTopic.ID != 7 {
		t.Fatalf("expected id 7, got %d", repo.updatedTopic.ID)
	}

	if repo.updatedTopic.Status != domain.REVIEWED {
		t.Fatalf("unexpected status: %s", repo.updatedTopic.Status)
	}
}

func TestUpdateTopicBadID(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/topics/nope",
		bytes.NewBufferString(`{}`),
	)
	req.SetPathValue("id", "nope")

	rec := httptest.NewRecorder()

	h.UpdateTopic(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateTopicNotFound(t *testing.T) {
	repo := &fakeRepo{
		updateTopicErr: sql.ErrNoRows,
	}
	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/topics/99",
		bytes.NewBufferString(`{
			"comment":"x",
			"status":"review",
			"last_reviewed":"2026-08-22"
		}`),
	)
	req.SetPathValue("id", "99")

	rec := httptest.NewRecorder()

	h.UpdateTopic(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestCreateTask(t *testing.T) {
	repo := &fakeRepo{
		createTaskID: 11,
	}

	h := NewHandlers(repo)

	body := []byte(`{
		"title":"Worker pool",
		"comment":"test",
		"status":"learned",
		"last_reviewed":"2026-08-22"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
	}

	if repo.createdTask == nil {
		t.Fatal("task was not created")
	}

	if repo.createdTask.Title != "Worker pool" {
		t.Fatalf("unexpected title: %s", repo.createdTask.Title)
	}

	var resp createResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.ID != 11 {
		t.Fatalf("expected id 11, got %d", resp.ID)
	}
}

func TestCreateTaskBadJSON(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBufferString(`{broken`),
	)
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandlers(repo)

	body := []byte(`{
		"comment":"updated",
		"status":"review",
		"last_reviewed":"2026-08-22"
	}`)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/tasks/5",
		bytes.NewReader(body),
	)
	req.SetPathValue("id", "5")

	rec := httptest.NewRecorder()

	h.UpdateTask(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rec.Code)
	}

	if repo.updatedTask == nil {
		t.Fatal("task was not updated")
	}

	if repo.updatedTask.ID != 5 {
		t.Fatalf("expected id 5, got %d", repo.updatedTask.ID)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	repo := &fakeRepo{
		updateTaskErr: sql.ErrNoRows,
	}

	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/tasks/99",
		bytes.NewBufferString(`{
			"comment":"x",
			"status":"review",
			"last_reviewed":"2026-08-22"
		}`),
	)
	req.SetPathValue("id", "99")

	rec := httptest.NewRecorder()

	h.UpdateTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestSetGoal(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandlers(repo)

	body := []byte(`{
		"id":1,
		"text":"Подготовка к senior",
		"comment":"Фокус на system design"
	}`)

	req := httptest.NewRequest(http.MethodPatch, "/goal", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.SetGoal(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rec.Code)
	}

	if repo.setGoal == nil {
		t.Fatal("goal was not updated")
	}

	if repo.setGoal.Text != "Подготовка к senior" {
		t.Fatalf("unexpected goal: %+v", repo.setGoal)
	}
}

func TestSetGoalBadJSON(t *testing.T) {
	repo := &fakeRepo{}
	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/goal",
		bytes.NewBufferString(`{broken`),
	)
	rec := httptest.NewRecorder()

	h.SetGoal(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSetGoalRepositoryError(t *testing.T) {
	repo := &fakeRepo{
		setGoalErr: errors.New("boom"),
	}

	h := NewHandlers(repo)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/goal",
		bytes.NewBufferString(`{
			"text":"test",
			"comment":"test"
		}`),
	)
	rec := httptest.NewRecorder()

	h.SetGoal(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
