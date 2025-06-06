package authentication

import (
	"github.com/guuzaa/email-newsletter/internal/api/middleware"
	"github.com/guuzaa/email-newsletter/internal/database/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type Credentials struct {
	Username string
	Password string
}

func (cred *Credentials) Validate(c *gin.Context, db *gorm.DB) bool {
	log := middleware.GetContextLogger(c)
	var user = models.User{
		Password: `$argon2id$v=19$m=15000,t=2,p=1$gZiV/M1gPc22ElAH/Jh1Hw$CWOrkoo7oJBQ/iyh7uJ0LO2aLEfrHwTWllSAxT0zRno`,
	}
	if err := db.Where("username = ?", cred.Username).First(&user).Error; err != nil {
		log.Trace().Err(err).Str("username", cred.Username).Msg("failed to find user")
	}

	valid, err := VerifyPassword(cred.Password, user.Password)
	if err != nil {
		log.Trace().Err(err).Msg("failed to verify password")
		return false
	}
	return valid
}

// HashPassword creates an Argon2id hash in PHC format
func HashPassword(password string) (string, error) {
	memory := uint32(15_000)
	iterations := uint32(2)
	parallelism := uint8(1)
	saltLength := uint32(16)
	keyLength := uint32(32)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// PHC format: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	phcString := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism,
		b64Salt, b64Hash,
	)

	return phcString, nil
}

// VerifyPassword verifies a password against a PHC formatted hash
func VerifyPassword(password, phcHash string) (bool, error) {
	// Parse PHC format: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	parts := strings.Split(phcHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid PHC hash format")
	}

	if parts[1] != "argon2id" {
		return false, errors.New("not an argon2id hash")
	}

	// Parse parameters
	var memory, iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Decode salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	keyLength := uint32(len(decodedHash))

	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
