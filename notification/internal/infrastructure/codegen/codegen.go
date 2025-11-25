package codegen

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/pkg/platform/config"
)

const maxCountDigits = 10
const minCountDigits = 4

// CodeGenerator type for generating verification codes
type CodeGenerator struct {
	expiresAt   time.Duration
	countDigits int
	maxN        int
	maxAttempts int
}

// NewCodeGenerator returns New CodeGenerator
func NewCodeGenerator(cfg *config.CodeGeneratorCfg) *CodeGenerator {
	count := cfg.CountDigits
	if count > maxCountDigits {
		count = maxCountDigits
	} else if count < minCountDigits {
		count = minCountDigits
	}

	maxN := int(math.Pow10(count))

	return &CodeGenerator{expiresAt: cfg.ExpiresAt, countDigits: count, maxN: maxN, maxAttempts: cfg.MaxAttempts}
}

// NewVerificationCode returns new domain.VerificationCode
func (c *CodeGenerator) NewVerificationCode(userID string) *domain.VerificationCode {
	return &domain.VerificationCode{
		UserID:      userID,
		Code:        c.generateCode(),
		ExpiresAt:   time.Now().Add(c.expiresAt),
		MaxAttempts: c.maxAttempts,
	}
}

// GenerateCode generates new code
func (c *CodeGenerator) generateCode() string {
	digit := rand.Intn(c.maxN)

	return fmt.Sprintf("%0*d", c.countDigits, digit)
}

func NewNoopCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		expiresAt:   4 * time.Minute,
		countDigits: 6,
		maxN:        int(math.Pow10(6)),
		maxAttempts: 5,
	}
}
