package task

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
	"sync"
)

var UnsafeManager = &Store{
	tasks: make(map[string]*Definition),
	mu:    &sync.RWMutex{},
}

type Store struct {
	tasks map[string]*Definition
	mu    *sync.RWMutex
}

func (s *Store) setDefinition(name string, tt *Definition) {
	// check that this task does not exist yet
	if s.UnsafeGetDefinition(name) != nil {
		panic(fmt.Sprint("Instance definition with name '", name, "' already exists"))
	}

	s.mu.Lock()
	s.tasks[name] = tt
	s.mu.Unlock()

}

func (s *Store) UnsafeGetDefinition(name string) *Definition {
	s.mu.RLock()
	t := s.tasks[name]
	s.mu.RUnlock()
	return t
}

func (s *Store) UnsafeGetInstance(name, encodedData string) *Instance {

	// get task class
	tt := s.UnsafeGetDefinition(name)

	// unserialize data
	err := serializer.Decode(encodedData, tt.defaultTask.Handler)
	if err != nil {
		panic(err)
	}

	return tt.defaultTask
}
