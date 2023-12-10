package main

import "sync"

type Mode struct {
	isRunning bool
	sync.Mutex
}

func (m *Mode) GetIsRunning() bool {
	m.Lock()
	defer m.Unlock()
	return m.isRunning
}

func (m *Mode) SetIsRunning(newBool bool) {
	m.Lock()
	defer m.Unlock()
	m.isRunning = newBool
}

func newMode() *Mode {
	return &Mode{isRunning: false}
}
