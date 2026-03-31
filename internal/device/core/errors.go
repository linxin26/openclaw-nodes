package core

import "errors"

var (
	ErrCapabilityUnavailable = errors.New("capability unavailable")
	ErrNotSupported          = errors.New("capability not supported")
	ErrRestricted            = errors.New("capability restricted")
)
