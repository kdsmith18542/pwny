// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import "errors"

var (
	ErrModuleNotFound     = errors.New("module not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrOptionNotFound     = errors.New("option not found")
	ErrMissingRequiredOpt = errors.New("missing required option")
	ErrNotImplemented     = errors.New("not implemented")
	ErrInvalidConnType    = errors.New("invalid connection type")
	ErrInteractiveUnavail = errors.New("interactive mode not available")
)
