package task

import (
	"fmt"
	"sync"
	"time"
)

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

type TaskInfo struct {
	ID         string     `json:"id"`
	Status     TaskStatus `json:"status"`
	Message    string     `json:"message"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	TemplateID string     `json:"template_id,omitempty"`
}

type Service struct {
	tasks map[string]*TaskInfo
	mu    sync.RWMutex
}

func NewService() *Service {
	return &Service{
		tasks: make(map[string]*TaskInfo),
	}
}

func (s *Service) CreateTask(templateID string) *TaskInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &TaskInfo{
		ID:         generateTaskID(),
		Status:     TaskPending,
		Message:    "任务已创建，等待执行",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		TemplateID: templateID,
	}

	s.tasks[task.ID] = task
	return task
}

func (s *Service) UpdateTask(taskID string, status TaskStatus, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.Status = status
		task.Message = message
		task.UpdatedAt = time.Now()
	}
}

func (s *Service) GetTask(taskID string) (*TaskInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, false
	}

	taskCopy := *task
	return &taskCopy, true
}

func (s *Service) GetAllTasks() []*TaskInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*TaskInfo, 0, len(s.tasks))
	for _, task := range s.tasks {
		taskCopy := *task
		tasks = append(tasks, &taskCopy)
	}
	return tasks
}

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}
