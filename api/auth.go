package api

import uuid "github.com/satori/go.uuid"

type Auth interface {
	Valid(string) (bool, uuid.UUID, error)
	Sign(uuid.UUID) (string, error)
}
