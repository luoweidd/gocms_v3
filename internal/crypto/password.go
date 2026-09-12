package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt对密码进行哈希加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword 验证密码是否匹配哈希值
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPasswordWithCost 使用指定成本因子进行哈希加密
func HashPasswordWithCost(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}
	return string(bytes), nil
}

// IsWeakPassword 检测弱密码
func IsWeakPassword(password string) bool {
	// 检查密码强度：长度至少8位
	if len(password) < 8 {
		return true
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	specialChars := "!@#$%^&*()_+-=[]{}|;:',.<>?/"

	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			for _, s := range specialChars {
				if r == s {
					hasSpecial = true
					break
				}
			}
		}
	}

	// 密码需要至少包含3种字符类型
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasDigit {
		count++
	}
	if hasSpecial {
		count++
	}

	return count < 3
}

// PasswordPolicy 密码策略配置
type PasswordPolicy struct {
	MinLength      int  // 最小长度，默认8
	RequireUpper   bool // 需要大写字母
	RequireLower   bool // 需要小写字母
	RequireDigit   bool // 需要数字
	RequireSpecial bool // 需要特殊字符
}

// DefaultPasswordPolicy 默认密码策略
var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:      8,
	RequireUpper:   true,
	RequireLower:   true,
	RequireDigit:   true,
	RequireSpecial: false,
}

// ValidatePassword 验证密码是否符合策略
func (p *PasswordPolicy) ValidatePassword(password string) error {
	if len(password) < p.MinLength {
		return fmt.Errorf("密码长度不能少于%d位", p.MinLength)
	}

	if p.RequireUpper && !hasUpperCase(password) {
		return fmt.Errorf("密码必须包含大写字母")
	}

	if p.RequireLower && !hasLowerCase(password) {
		return fmt.Errorf("密码必须包含小写字母")
	}

	if p.RequireDigit && !hasNumber(password) {
		return fmt.Errorf("密码必须包含数字")
	}

	if p.RequireSpecial && !hasSpecialChar(password) {
		return fmt.Errorf("密码必须包含特殊字符")
	}

	return nil
}

func hasUpperCase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func hasLowerCase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:',.<>?/"
	for _, r := range s {
		for _, c := range specialChars {
			if r == c {
				return true
			}
		}
	}
	return false
}
