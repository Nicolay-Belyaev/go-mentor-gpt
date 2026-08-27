package domain

import "context"

type GoalRepository interface {
	GetGoal(ctx context.Context) (*Goal, error)
	SetGoal(ctx context.Context, goal *Goal) error
}

type TopicRepository interface {
	GetAllTopics(ctx context.Context) ([]Topic, error)
	CreateTopic(ctx context.Context, topic Topic) (int64, error)
	UpdateTopic(ctx context.Context, topic Topic) error
}

type TaskRepository interface {
	GetAllTasks(ctx context.Context) ([]Task, error)
	CreateTask(ctx context.Context, task Task) (int64, error)
	UpdateTask(ctx context.Context, task Task) error
}
