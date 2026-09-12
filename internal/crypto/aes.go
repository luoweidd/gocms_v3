package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// AESCipher AES-256-GCM加解密器
type AESCipher struct {
	key []byte // 32字节密钥
}

// NewAESCipher 创建AES加密器
func NewAESCipher(secret string) *AESCipher {
	// 使用SHA256将任意长度密钥转换为32字节
	hash := sha256.Sum256([]byte(secret))
	return &AESCipher{
		key: hash[:],
	}
}

// Encrypt 加密字符串，返回Base64编码的密文(含IV)
func (c *AESCipher) Encrypt(plaintext string) (string, error) {
	plainBytes := []byte(plaintext)

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("创建AES cipher失败: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	// 生成随机nonce (IV)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", fmt.Errorf("生成nonce失败: %w", err)
	}

	// 加密并附加认证标签
	ciphertext := aesGCM.Seal(nonce, nonce, plainBytes, nil)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密Base64编码的密文，返回原始字符串
func (c *AESCipher) Decrypt(ciphertextBase64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("创建AES cipher失败: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("密文数据不完整")
	}

	nonce, ciphertextData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密并验证认证标签
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}

	return string(plaintext), nil
}

// EncryptBytes 加密字节数组
func (c *AESCipher) EncryptBytes(data []byte) (string, error) {
	return c.Encrypt(string(data))
}

// DecryptBytes 解密并返回字节数组
func (c *AESCipher) DecryptBytes(ciphertextBase64 string) ([]byte, error) {
	plaintext, err := c.Decrypt(ciphertextBase64)
	if err != nil {
		return nil, err
	}
	return []byte(plaintext), nil
}
