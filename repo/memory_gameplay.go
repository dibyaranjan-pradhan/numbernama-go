package repo

import (
	"sync"

	"numbernama-go/model"
)

// MemoryGameplay keeps one Board per WebSocket client id (RAM-only).
type MemoryGameplay struct {
	mu     sync.Mutex
	states map[string]*model.GameState
}

// NewMemoryGameplay constructs an empty in-memory session store.
func NewMemoryGameplay() *MemoryGameplay {
	return &MemoryGameplay{states: make(map[string]*model.GameState)}
}

// State returns gameplay state for clientID, creating it if missing.
func (m *MemoryGameplay) State(clientID string) *model.GameState {
	m.mu.Lock()
	defer m.mu.Unlock()
	state, ok := m.states[clientID]
	if !ok {
		state = &model.GameState{}
		m.states[clientID] = state
	}
	return state
}

// Remove deletes a session (optional; used on disconnect if desired).
func (m *MemoryGameplay) Remove(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, clientID)
}
