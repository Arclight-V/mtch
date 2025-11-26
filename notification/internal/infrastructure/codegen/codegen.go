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
	expiresAt    time.Duration
	countDigits  int
	maxN         int
	maxAttempts  int
	generateFunc func(int, int) string
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

	return &CodeGenerator{
		expiresAt:    cfg.ExpiresAt,
		countDigits:  count,
		maxN:         maxN,
		maxAttempts:  cfg.MaxAttempts,
		generateFunc: generateCode,
	}
}

// NewVerificationCode returns new domain.VerificationCode
func (c *CodeGenerator) NewVerificationCode(userID string) *domain.VerificationCode {
	return &domain.VerificationCode{
		UserID:      userID,
		Code:        c.generateFunc(c.maxN, c.countDigits),
		ExpiresAt:   time.Now().Add(c.expiresAt),
		MaxAttempts: c.maxAttempts,
	}
}

// GenerateCode generates new code
func generateCode(maxN, width int) string {
	digit := rand.Intn(maxN)
	return fmt.Sprintf("%0*d", width, digit)
}

// NewNoopCodeGenerator returns Noop CodeGenerator
func NewNoopCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		expiresAt:   4 * time.Minute,
		countDigits: 6,
		maxN:        int(math.Pow10(6)),
		maxAttempts: 5,
		generateFunc: func(int, int) string {
			return fmt.Sprintf("%0*d", 6, int(math.Pow10(6)))
		},
	}
}
