package agent

// Agent interface represents an internal agent implementation (for plugin pattern)
type Agent interface {
	Name() string
	CanHandle(action string) bool
	Handle(action string, payload map[string]interface{}) (map[string]interface{}, error)
}