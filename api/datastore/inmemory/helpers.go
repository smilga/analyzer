package inmemory

import uuid "github.com/satori/go.uuid"

func inSlice(id uuid.UUID, ids []uuid.UUID) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}
