package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2IdConfig struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

const argon2IdDefaultMemory = 64 * 1024

func MakeArgon2IdHash(password string) string {
	return MakeArgon2IdHashWithConfig(
		password,
		&Argon2IdConfig{
			Memory:  argon2IdDefaultMemory,
			Time:    1,
			Threads: 2,
			KeyLen:  32,
			SaltLen: 16,
		})
}

func MakeArgon2IdHashWithConfig(password string, config *Argon2IdConfig) string {
	salt := make([]byte, config.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	hash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.Memory,
		config.Time,
		config.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

func VerifyArgon2IdHash(password string, encodedHash string) bool {
	strParts := strings.Split(encodedHash, "$")
	var memory uint32
	var time uint32
	var threads uint8

	_, err := fmt.Sscanf(strParts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		panic(err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(strParts[4])
	if err != nil {
		panic(err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(strParts[5])
	if err != nil {
		panic(err)
	}
	i := subtle.ConstantTimeCompare(hash, argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash))))

	return i == 1
}
