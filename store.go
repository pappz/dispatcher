package dispatcher

import (
	"sync"
)

type Store struct {
	devices map[string]*Device

	lock sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		devices: make(map[string]*Device),
	}
}

func (s *Store) AddDevice(device *Device) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.devices[device.id] = device
}

func (s *Store) Device(id string) (*Device, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	p, ok := s.devices[id]
	return p, ok
}

func (s *Store) RemoveDevice(d *Device) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.devices, d.id)
}
