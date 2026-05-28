package repository

import "errors"

// ErrNotFound is returned when a queried record does not exist.
// Callers should use errors.Is(err, repository.ErrNotFound) to check for this
// condition instead of comparing the returned pointer to nil.
var ErrNotFound = errors.New("record not found")
