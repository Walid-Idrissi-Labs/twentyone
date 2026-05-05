package anim

type FlashID string

type Flash struct {
	ID        FlashID
	Active    bool
	TicksLeft int
}

type Manager struct {
	flashes map[FlashID]*Flash
}

func NewManager() *Manager {
	return &Manager{
		flashes: make(map[FlashID]*Flash),
	}
}

func (m *Manager) Trigger(id FlashID, ticks int) {
	m.flashes[id] = &Flash{
		ID:        id,
		Active:    true,
		TicksLeft: ticks,
	}
}

func (m *Manager) Tick() []FlashID {
	var expired []FlashID
	for id, flash := range m.flashes {
		flash.TicksLeft--
		if flash.TicksLeft <= 0 {
			flash.Active = false
			expired = append(expired, id)
			delete(m.flashes, id)
		}
	}
	return expired
}

func (m *Manager) IsActive(id FlashID) bool {
	if f, ok := m.flashes[id]; ok {
		return f.Active
	}
	return false
}

func (m *Manager) SetActive(id FlashID, active bool) {
	if f, ok := m.flashes[id]; ok {
		f.Active = active
	}
}