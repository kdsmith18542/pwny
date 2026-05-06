// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package encoder

import (
	"crypto/rand"
	"github.com/kdsmith18542/pwny/internal/payload"
)

type XOREncoder struct{}

func (e *XOREncoder) Encode(data []byte, opts payload.PayloadOptions) ([]byte, error) {
	key, err := GenerateRandomKey(4) // Default to 4-byte key
	if err != nil {
		return nil, err
	}
	
	// Prepend the key to the encoded data so it can be decoded
	encoded := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		encoded[i] = data[i] ^ key[i%len(key)]
	}
	
	// TODO: Add XOR decoder stub to the beginning of the payload
	
	return encoded, nil
}

// GenerateRandomKey generates a random XOR key of specified length
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
