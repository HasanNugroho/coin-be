package session

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	TelegramID int64
	UserID     primitive.ObjectID
	State      string
	TempData   map[string]string
}

type Store struct {
	mu       sync.RWMutex
	sessions map[int64]*UserSession
}

func NewStore() *Store {
	return &Store{
		sessions: make(map[int64]*UserSession),
	}
}

func (s *Store) Get(telegramID int64) *UserSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[telegramID]
}

func (s *Store) GetOrCreate(telegramID int64) *UserSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sess, ok := s.sessions[telegramID]; ok {
		return sess
	}

	sess := &UserSession{
		TelegramID: telegramID,
		TempData:   make(map[string]string),
	}
	s.sessions[telegramID] = sess
	return sess
}

func (s *Store) Set(telegramID int64, sess *UserSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[telegramID] = sess
}

func (s *Store) ClearState(telegramID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[telegramID]; ok {
		sess.State = ""
		sess.TempData = make(map[string]string)
	}
}
