package signature

import (
	"fmt"
	"strconv"
	"time"
)

// Verifier API签名验证器
type Verifier struct {
	signer     *Signer
	maxAge     int64 // 最大允许的时间差（秒）
	requiredTS bool  // 是否必须提供timestamp
}

// VerifierConfig 验证器配置
type VerifierConfig struct {
	MaxAge     int64  // 最大时间差，默认300秒（5分钟）
	RequiredTS bool   // 是否强制要求时间戳
	SecretKey  string // 密钥
}

// DefaultVerifierConfig 默认配置
var DefaultVerifierConfig = VerifierConfig{
	MaxAge:     300,
	RequiredTS: true,
	SecretKey:  "",
}

// NewVerifier 创建签名验证器
func NewVerifier(config VerifierConfig) *Verifier {
	if config.SecretKey == "" {
		config.SecretKey = DefaultVerifierConfig.SecretKey
	}
	if config.MaxAge == 0 {
		config.MaxAge = 300
	}
	return &Verifier{
		signer:     NewSigner(config.SecretKey),
		maxAge:     config.MaxAge,
		requiredTS: config.RequiredTS,
	}
}

// Verify 验证请求签名
// params: 请求中的参数字典（不包含sign字段）
// sign: 请求头或参数中的签名值
func (v *Verifier) Verify(params map[string]string, sign string) error {
	if sign == "" {
		return fmt.Errorf("缺少签名参数")
	}

	tsStr := params["timestamp"]
	if v.requiredTS && tsStr == "" {
		return fmt.Errorf("缺少timestamp参数")
	}

	if tsStr != "" {
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid timestamp: %w", err)
		}

		if time.Now().Unix()-ts > v.maxAge {
			return fmt.Errorf("请求已过期，时间戳差值超过%d秒", v.maxAge)
		}
	}

	expectedSign := v.signer.GenerateSign(params, getTimestamp(params))
	if sign != expectedSign {
		return fmt.Errorf("签名验证失败")
	}

	return nil
}

// VerifyFromHeaders 从请求头中获取签名并验证
func (v *Verifier) VerifyFromHeaders(header map[string]string, params map[string]string) error {
	sign := header["X-Sign"]
	if sign == "" {
		sign = header["Sign"]
	}
	return v.Verify(params, sign)
}

// VerifyWithApp 验证应用级签名（含app_id和app_secret）
func (v *Verifier) VerifyWithApp(appID string, appSecret string, params map[string]string, sign string) error {
	// 创建专用密钥
	appSigner := NewSigner(appSecret)

	tsStr := params["timestamp"]
	if v.requiredTS && tsStr == "" {
		return fmt.Errorf("缺少timestamp参数")
	}

	var ts int64
	if tsStr != "" {
		var err error
		ts, err = strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid timestamp: %w", err)
		}
		if time.Now().Unix()-ts > v.maxAge {
			return fmt.Errorf("请求已过期")
		}
	}

	expectedSign := appSigner.GenerateSign(params, ts)
	if sign != expectedSign {
		return fmt.Errorf("应用签名验证失败")
	}

	return nil
}

// GenerateRequest 生成完整请求参数（含签名和时间戳）
func (v *Verifier) GenerateRequest(params map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range params {
		result[k] = v
	}

	sign, ts := v.signer.GenerateTimestampedSign(result)
	result["timestamp"] = strconv.FormatInt(ts, 10)
	result["sign"] = sign

	return result
}

// getTimestamp 从参数字典中获取时间戳
func getTimestamp(params map[string]string) int64 {
	tsStr := params["timestamp"]
	if tsStr == "" {
		return time.Now().Unix()
	}
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return time.Now().Unix()
	}
	return ts
}
