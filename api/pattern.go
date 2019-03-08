package api

import uuid "github.com/satori/go.uuid"

// Types of patterns possible
const (
	JSSource PatternType = "js_source"
	HTML     PatternType = "html"
	Resource PatternType = "resource"
)

type PatternType string

type Pattern struct {
	ID        uuid.UUID
	Type      PatternType
	Value     string
	Mandatory bool
}
