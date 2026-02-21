package otp

import (
	"sync"
	"time"
)

type Entry struct {
	OTP        string
	TelegramID int64
	ExpiresAt  time.Time
}

type Store struct {
	mu      sync.Mutex
	entries map[string]Entry // key: email
}

func NewStore() *Store {
	store := &Store{
		entries: make(map[string]Entry),
	}
	go store.cleanupExpired()
	return store
}

func (s *Store) Set(email string, entry Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[email] = entry
}

func (s *Store) Get(email string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.entries[email]
	if !ok {
		return Entry{}, false
	}
	if time.Now().After(entry.ExpiresAt) {
		delete(s.entries, email)
		return Entry{}, false
	}
	return entry, true
}

func (s *Store) Delete(email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, email)
}

func (s *Store) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for email, entry := range s.entries {
			if now.After(entry.ExpiresAt) {
				delete(s.entries, email)
			}
		}
		s.mu.Unlock()
	}
}
