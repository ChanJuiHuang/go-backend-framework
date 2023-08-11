package argon2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgon2(t *testing.T) {
	password := "password"
	hash := MakeArgon2IdHash(password)
	assert.True(t, VerifyArgon2IdHash(password, hash))
}
