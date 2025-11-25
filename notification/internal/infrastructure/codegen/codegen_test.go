package codegen

import (
	"math"
	"testing"
	"time"

	"github.com/Arclight-V/mtch/pkg/platform/config"
)

func TestNewCodeGenerator(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.CodeGeneratorCfg
		want *CodeGenerator
	}{
		{
			name: "success",
			cfg: &config.CodeGeneratorCfg{
				ExpiresAt:   3 * time.Minute,
				CountDigits: 6,
				MaxAttempts: 10,
			},
			want: &CodeGenerator{
				expiresAt:   3 * time.Minute,
				countDigits: 6,
				maxAttempts: 10,
				maxN:        int(math.Pow10(6)),
			},
		},
		{
			name: "max count digits",
			cfg: &config.CodeGeneratorCfg{
				ExpiresAt:   3 * time.Minute,
				CountDigits: 11,
				MaxAttempts: 10,
			},
			want: &CodeGenerator{
				expiresAt:   3 * time.Minute,
				countDigits: maxCountDigits,
				maxAttempts: 10,
				maxN:        int(math.Pow10(maxCountDigits)),
			},
		},

		{
			name: "min count digits",
			cfg: &config.CodeGeneratorCfg{
				ExpiresAt:   3 * time.Minute,
				CountDigits: -1,
				MaxAttempts: 10,
			},
			want: &CodeGenerator{
				expiresAt:   3 * time.Minute,
				countDigits: minCountDigits,
				maxAttempts: 10,
				maxN:        int(math.Pow10(minCountDigits)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newGen := NewCodeGenerator(tt.cfg)
			if newGen.expiresAt != tt.want.expiresAt {
				t.Errorf("Error for set expiresAt want: %v, got: %v", tt.want.expiresAt, newGen.expiresAt)
			}
			if newGen.countDigits != tt.want.countDigits {
				t.Errorf("Error for set countDigits want: %v, got: %v", tt.want.countDigits, newGen.countDigits)
			}
			if newGen.maxAttempts != tt.want.maxAttempts {
				t.Errorf("Error for set maxAttempts want: %v, got: %v", tt.want.maxAttempts, newGen.maxAttempts)
			}
			if newGen.maxN != tt.want.maxN {
				t.Errorf("Error for set maxN want: %v, got: %v", tt.want.maxN, newGen.maxN)
			}

		})
	}
}

func TestCodeGenerator_Generate(t *testing.T) {
	cfg := &config.CodeGeneratorCfg{
		ExpiresAt:   3 * time.Minute,
		CountDigits: 6,
		MaxAttempts: 10,
	}

	newGen := NewCodeGenerator(cfg)
	code := newGen.generateCode()
	if len(code) != newGen.countDigits {
		t.Errorf("Error len code want:%v, got: %v", newGen.countDigits, len(code))
	}

}
