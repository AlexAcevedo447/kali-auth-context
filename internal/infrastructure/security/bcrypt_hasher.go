package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"kali-auth-context/internal/ports"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

type Argon2idHasher struct{}

var _ ports.IPasswordHasher = (*Argon2idHasher)(nil)

const (
	argonVersion = argon2.Version
	argonMemory  = 64 * 1024
	argonTime    = 3
	argonThreads = 1
	argonKeyLen  = 32
	argonSaltLen = 16
)

func NewArgon2idHasher() *Argon2idHasher {
	return &Argon2idHasher{}
}

func (h *Argon2idHasher) Hash(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argonVersion, argonMemory, argonTime, argonThreads, b64Salt, b64Hash)
	return encoded, nil
}

func (h *Argon2idHasher) Compare(hash string, password string) error {
	if strings.HasPrefix(hash, "$2") {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	}

	params, salt, expectedHash, err := parseArgon2idHash(hash)
	if err != nil {
		return err
	}

	calculatedHash := argon2.IDKey([]byte(password), salt, params.time, params.memory, params.threads, uint32(len(expectedHash)))
	if subtle.ConstantTimeCompare(expectedHash, calculatedHash) != 1 {
		return errors.New("invalid credentials")
	}

	return nil
}

func (h *Argon2idHasher) NeedsRehash(hash string) bool {
	if strings.HasPrefix(hash, "$2") {
		return true
	}

	params, _, _, err := parseArgon2idHash(hash)
	if err != nil {
		return true
	}

	return params.memory != argonMemory || params.time != argonTime || params.threads != argonThreads || params.version != argonVersion
}

type argonParams struct {
	memory  uint32
	time    uint32
	threads uint8
	version int
}

func parseArgon2idHash(encodedHash string) (*argonParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid argon2 hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported hash algorithm")
	}

	versionField := strings.TrimPrefix(parts[2], "v=")
	version, err := strconv.Atoi(versionField)
	if err != nil {
		return nil, nil, nil, err
	}

	var memory uint64
	var time uint64
	var threads uint64
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	params := &argonParams{
		memory:  uint32(memory),
		time:    uint32(time),
		threads: uint8(threads),
		version: version,
	}

	return params, salt, hash, nil
}