package agent

// Base represents the base agent structure with common fields
type Base struct {
	ID        string
	Name      string
	Status    string
	CreatedAt int64
	UpdatedAt int64
}

// NewBase creates a new base agent with default values
func NewBase(id, name string) *Base {
	now := time.Now().Unix()
	return &Base{
		ID:        id,
		Name:      name,
		Status:    "initialized",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GetID returns the agent's ID
func (b *Base) GetID() string {
	return b.ID
}

// GetStatus returns the agent's current status
func (b *Base) GetStatus() string {
	return b.Status
}

// UpdateStatus sets a new status and updates the UpdatedAt timestamp
func (b *Base) UpdateStatus(status string) {
	b.Status = status
	b.UpdatedAt = time.Now().Unix()
}