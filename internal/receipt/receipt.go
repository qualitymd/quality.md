// Package receipt contains shared agent-facing receipt elements.
package receipt

// Action is a suggested next action.
type Action struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Command string `json:"command,omitempty"`
}
