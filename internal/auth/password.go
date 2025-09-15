package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash is returned when the encoded hash is not in the correct format
	ErrInvalidHash = errors.New("invalid hash format")
	// ErrIncompatibleVersion is returned when the hash version is incompatible
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type params struct {
	memory      uint32
time        uint32
	threads    uint8
	keyLen     uint32
	saltLength uint32
}

// Default parameters for Argon2id hashing
var defaultParams = &params{
	memory:      64 * 1024, // 64 MB
	time:        3,
	threads:     4,
	saltLength:  16,
	keyLen:      32, // 32 bytes = 256 bits
}

// HashPassword hashes a password using Argon2id
func HashPassword(password string) (string, error) {
	// Generate a cryptographically secure random salt
	salt := make([]byte, defaultParams.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate the hash using Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.time,
		defaultParams.memory,
		defaultParams.threads,
		defaultParams.keyLen,
	)

	// Base64 encode the salt and hashed password
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return the encoded string with parameters
	// Format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.memory,
		defaultParams.time,
		defaultParams.threads,
		b64Salt,
		b64Hash,
	), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generate the hash using the same parameters
	hashToCompare := argon2.IDKey(
		[]byte(password),
		salt,
		p.time,
		p.memory,
		p.threads,
		p.keyLen,
	)

	// Compare the hashes in constant time
	if subtle.ConstantTimeCompare(hash, hashToCompare) == 1 {
		return true, nil
	}

	return false, nil
}

// decodeHash decodes the encoded hash into its components
func decodeHash(encodedHash string) (*params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidHash
	}

	// Parse version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Parse parameters
	p := &params{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.time, &p.threads); err != nil {
		return nil, nil, nil, err
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	// Decode hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLen = uint32(len(hash))

	return p, salt, hash, nil
}
