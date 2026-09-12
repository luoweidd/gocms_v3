package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

// Signer HMAC-SHA256签名生成器
type Signer struct {
	secretKey []byte
}

// NewSigner 创建签名器
func NewSigner(secretKey string) *Signer {
	return &Signer{
		secretKey: []byte(secretKey),
	}
}

// GenerateSign 生成API签名
// params: 待签名的参数字典（已按key排序）
// timestamp: 时间戳（Unix秒）
func (s *Signer) GenerateSign(params map[string]string, timestamp int64) string {
	// 拼接参数为 key1=value1&key2=value2 格式
	signString := fmt.Sprintf("timestamp=%d", timestamp)
	for k, v := range params {
		if k != "sign" && k != "timestamp" {
			signString += "&" + k + "=" + v
		}
	}

	// 追加密钥: key1=value1&key2=value2&secret=your_secret_key
	signString += "&secret=" + string(s.secretKey)

	// HMAC-SHA256签名
	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write([]byte(signString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// GenerateTimestampedSign 生成带时间戳的签名（防重放）
func (s *Signer) GenerateTimestampedSign(params map[string]string) (string, int64) {
	timestamp := time.Now().Unix()
	sign := s.GenerateSign(params, timestamp)
	return sign, timestamp
}

// VerifySign 验证API签名
func (s *Signer) VerifySign(params map[string]string, sign string, maxAge int64) error {
	// 获取请求中的时间戳
	tsStr := params["timestamp"]
	if tsStr == "" {
		return fmt.Errorf("缺少timestamp参数")
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}

	// 检查时间戳是否在允许范围内（默认5分钟）
	if maxAge == 0 {
		maxAge = 300 // 5分钟
	}
	if time.Now().Unix()-ts > maxAge {
		return fmt.Errorf("签名已过期")
	}

	// 重新计算签名
	expectedSign := s.GenerateSign(params, ts)
	if sign != expectedSign {
		return fmt.Errorf("签名验证失败")
	}

	return nil
}

// GenerateAppSign 为应用生成签名（含app_id）
func (s *Signer) GenerateAppSign(appID string, params map[string]string) (string, int64) {
	timestamp := time.Now().Unix()
	params["app_id"] = appID
	params["timestamp"] = strconv.FormatInt(timestamp, 10)

	// 按key排序确保签名一致性
	sortedParams := make(map[string]string)
	for k, v := range params {
		if k != "sign" {
			sortedParams[k] = v
		}
	}

	sign := s.GenerateSign(sortedParams, timestamp)
	return sign, timestamp
}
