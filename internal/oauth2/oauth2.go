package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Provider 第三方登录提供商
type Provider int

const (
	ProviderWechat Provider = iota
	ProviderQQ
	ProviderGitHub
)

func (p Provider) String() string {
	return [...]string{"wechat", "qq", "github"}[p]
}

// OAuth2Config OAuth2配置
type OAuth2Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

// OAuth2Client OAuth2客户端
type OAuth2Client struct {
	Provider Provider
	Config   OAuth2Config
	httpCli  *http.Client
}

// NewOAuth2Client 创建OAuth2客户端
func NewOAuth2Client(provider Provider, config OAuth2Config) *OAuth2Client {
	return &OAuth2Client{
		Provider: provider,
		Config:   config,
		httpCli:  &http.Client{Timeout: 10 * time.Second},
	}
}

// AuthURL 获取授权URL
func (c *OAuth2Client) AuthURL(state string) string {
	switch c.Provider {
	case ProviderWechat:
		return fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
			c.Config.ClientID, c.Config.RedirectURL, state)
	case ProviderQQ:
		return fmt.Sprintf("https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
			c.Config.ClientID, c.Config.RedirectURL, state)
	case ProviderGitHub:
		return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s",
			c.Config.ClientID, c.Config.RedirectURL, state)
	default:
		return ""
	}
}

// ExchangeCode 用授权码交换Token
func (c *OAuth2Client) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	var tokenResp TokenResponse
	var reqErr error

	switch c.Provider {
	case ProviderWechat:
		tokenResp, reqErr = c.exchangeWechat(ctx, code)
	case ProviderQQ:
		tokenResp, reqErr = c.exchangeQQ(ctx, code)
	case ProviderGitHub:
		tokenResp, reqErr = c.exchangeGitHub(ctx, code)
	default:
		return nil, fmt.Errorf("不支持的提供商: %s", c.Provider)
	}

	if reqErr != nil {
		return nil, fmt.Errorf("%s授权码交换失败: %w", c.Provider, reqErr)
	}

	return &tokenResp, nil
}

// exchangeWechat 微信授权码交换
func (c *OAuth2Client) exchangeWechat(ctx context.Context, code string) (TokenResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		c.Config.ClientID, c.Config.ClientSecret, code)

	resp, err := c.httpCli.Get(url)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return TokenResponse{}, err
	}

	// 检查是否有错误码
	if _, ok := result["errcode"]; ok {
		msg, _ := result["errmsg"].(string)
		return TokenResponse{}, fmt.Errorf("微信错误: %s", msg)
	}

	return TokenResponse{
		AccessToken: result["access_token"].(string),
		OpenID:      result["openid"].(string),
		ExpiresIn:   int64(result["expires_in"].(float64)),
	}, nil
}

// exchangeQQ QQ授权码交换
func (c *OAuth2Client) exchangeQQ(ctx context.Context, code string) (TokenResponse, error) {
	url := fmt.Sprintf("https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		c.Config.ClientID, c.Config.ClientSecret, code, c.Config.RedirectURL)

	resp, err := http.Get(url)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, err
	}

	// QQ返回的是key=value&key=value格式
	values := make(map[string]string)
	parts := splitString(string(body), "&")
	for _, part := range parts {
		kv := splitString(part, "=")
		if len(kv) == 2 {
			values[kv[0]] = kv[1]
		}
	}

	return TokenResponse{
		AccessToken: values["access_token"],
		ExpiresIn:   parseQQInt(values["expires"]),
	}, nil
}

// exchangeGitHub GitHub授权码交换
func (c *OAuth2Client) exchangeGitHub(ctx context.Context, code string) (TokenResponse, error) {
	data := map[string]string{
		"client_id":     c.Config.ClientID,
		"client_secret": c.Config.ClientSecret,
		"code":          code,
		"redirect_uri":  c.Config.RedirectURL,
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return TokenResponse{}, err
	}

	bodyBytes, _ := json.Marshal(data)
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return TokenResponse{}, err
	}

	if _, ok := result["error"]; ok {
		msg, _ := result["error_description"].(string)
		return TokenResponse{}, fmt.Errorf("GitHub错误: %s", msg)
	}

	return TokenResponse{
		AccessToken: result["access_token"].(string),
		Scope:       result["scope"].(string),
		TokenType:   "bearer",
	}, nil
}

// GetUserInfo 获取用户信息
func (c *OAuth2Client) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	switch c.Provider {
	case ProviderWechat:
		return c.getUserInfoWechat(ctx, token)
	case ProviderQQ:
		return c.getUserInfoQQ(ctx, token)
	case ProviderGitHub:
		return c.getUserInfoGitHub(ctx, token)
	default:
		return nil, fmt.Errorf("不支持的提供商: %s", c.Provider)
	}
}

// getUserInfoWechat 微信获取用户信息
func (c *OAuth2Client) getUserInfoWechat(ctx context.Context, accessToken string) (*UserInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, c.Config.ClientID)

	resp, err := c.httpCli.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		Provider: c.Provider.String(),
		Nickname: fmt.Sprintf("%v", result["nickname"]),
	}
	if avatarurl, ok := result["headimgurl"].(string); ok {
		userInfo.Avatar = avatarurl
	}
	if openid, ok := result["openid"].(string); ok {
		userInfo.OpenID = openid
	}

	return userInfo, nil
}

// getUserInfoQQ QQ获取用户信息
func (c *OAuth2Client) getUserInfoQQ(ctx context.Context, accessToken string) (*UserInfo, error) {
	url := fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&openid=%s", accessToken, c.Config.ClientID)

	resp, err := c.httpCli.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		Provider: c.Provider.String(),
		Nickname: fmt.Sprintf("%v", result["nickname"]),
		Avatar:   fmt.Sprintf("%v", result["figureurl_qq"]),
	}

	return userInfo, nil
}

// getUserInfoGitHub GitHub获取用户信息
func (c *OAuth2Client) getUserInfoGitHub(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		Provider: c.Provider.String(),
		Nickname: fmt.Sprintf("%v", result["login"]),
		Name:     fmt.Sprintf("%v", result["name"]),
		Email:    fmt.Sprintf("%v", result["email"]),
		Avatar:   fmt.Sprintf("%v", result["avatar_url"]),
	}

	return userInfo, nil
}

// TokenResponse OAuth2令牌响应
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	OpenID      string `json:"openid,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

// UserInfo OAuth2用户信息
type UserInfo struct {
	Provider string `json:"provider"` // wechat, qq, github
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// splitString 分割字符串
func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

// parseQQInt 解析QQ返回的整数字符串
func parseQQInt(s string) int64 {
	var n int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		}
	}
	return n
}
