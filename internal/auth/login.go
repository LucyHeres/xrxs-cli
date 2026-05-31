package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/LucyHeres/xrxs-cli/pkg/encrypt"
)

type loginResponse struct {
	Code    json.RawMessage `json:"code"`
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (lr loginResponse) codeIsZero() bool {
	var s string
	if json.Unmarshal(lr.Code, &s) == nil {
		return s == "0"
	}
	var n int
	if json.Unmarshal(lr.Code, &n) == nil {
		return n == 0
	}
	return false
}

func (lr loginResponse) codeStr() string {
	var s string
	if json.Unmarshal(lr.Code, &s) == nil {
		return s
	}
	return "unknown"
}

type loginInfoData struct {
	PasswordKey string `json:"passwordKey"`
}

type loginResultData struct {
	SSOToken string `json:"ssotoken"`
	Redirect string `json:"redirect"`
}

type predataData struct {
	CSRFToken string `json:"csrfToken"`
}

// Login performs password login and returns a Session with cookies and CSRF token.
func Login(baseURL, username, password string) (*Session, error) {
	baseURL = strings.TrimRight(baseURL, "/")

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	var lr loginResponse

	// Step 1: Get passwordKey (sets zcp cookie)
	infoReq, _ := http.NewRequest("POST", baseURL+"/account-center/service/sso/ajax-get-login-info",
		strings.NewReader("fromUrl=&fromType=&appId=app-admin"))
	infoReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	infoResp, err := client.Do(infoReq)
	if err != nil {
		return nil, fmt.Errorf("获取登录信息失败: %w", err)
	}
	if err := json.NewDecoder(infoResp.Body).Decode(&lr); err != nil {
		infoResp.Body.Close()
		return nil, fmt.Errorf("解析登录信息失败: %w", err)
	}
	infoResp.Body.Close()
	if !lr.codeIsZero() {
		return nil, fmt.Errorf("获取登录信息失败: %s (code=%s)", lr.Message, lr.codeStr())
	}

	var info loginInfoData
	if err := json.Unmarshal(lr.Data, &info); err != nil || info.PasswordKey == "" {
		info.PasswordKey = encrypt.DefaultEncryptKey
	}

	// Step 2: Login with encrypted password
	encPwd, err := encrypt.RC4Encrypt(password, info.PasswordKey)
	if err != nil {
		return nil, fmt.Errorf("加密密码失败: %w", err)
	}

	body := url.Values{}
	body.Set("verifyMode", "0")
	body.Set("accountName", username)
	body.Set("password", encPwd)
	body.Set("passwordKey", info.PasswordKey)

	req, _ := http.NewRequest("POST", baseURL+"/account-center/service/sso/ajax-password-login", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("登录请求失败: %w", err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err := json.Unmarshal(respBody, &lr); err != nil {
		return nil, fmt.Errorf("解析登录响应: %w", err)
	}
	if !lr.codeIsZero() {
		return nil, fmt.Errorf("登录失败: %s (code=%s)", lr.Message, lr.codeStr())
	}

	var loginData loginResultData
	if err := json.Unmarshal(lr.Data, &loginData); err != nil {
		return nil, fmt.Errorf("解析登录结果: %w", err)
	}

	// Step 3: Exchange ssotoken for real session cookies
	predataURL := fmt.Sprintf("%s/support/service/storm/ajax-get-predata-v2?ssotoken=%s",
		baseURL, url.QueryEscape(loginData.SSOToken))

	predataResp, err := client.Get(predataURL)
	if err != nil {
		return nil, fmt.Errorf("获取 session 失败: %w", err)
	}
	defer predataResp.Body.Close()

	var predataLR loginResponse
	if err := json.NewDecoder(predataResp.Body).Decode(&predataLR); err != nil {
		return nil, fmt.Errorf("解析 session 响应: %w", err)
	}
	if !predataLR.codeIsZero() {
		return nil, fmt.Errorf("获取 session 失败: %s (code=%s)", predataLR.Message, predataLR.codeStr())
	}

	// Collect all cookies including session cookies from step 3
	u, _ := url.Parse(baseURL)
	allCookies := jar.Cookies(u)

	// Extract CSRF token from predata response
	csrfToken := ""
	var pd predataData
	if json.Unmarshal(predataLR.Data, &pd) == nil && pd.CSRFToken != "" {
		csrfToken = pd.CSRFToken
	}

	return &Session{
		Cookies:   allCookies,
		CSRFToken: csrfToken,
		BaseURL:   baseURL,
		CreatedAt: time.Now(),
	}, nil
}

// Company represents a company/subsidiary from the switch-company list.
type Company struct {
	ID   string
	Name string
}

type companyItem struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	IsActive int           `json:"isActive"`
	Level    int           `json:"level"`
	Switch   bool          `json:"switch"`
	ItemList []companyItem `json:"itemList"`
}

type companyListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    struct {
		CurrentID string        `json:"currentId"`
		ItemList  []companyItem `json:"itemList"`
	} `json:"data"`
}

type switchCompanyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

// FetchCompanyList gets the list of switchable companies for the current session.
func FetchCompanyList(baseURL string, cookies []*http.Cookie, csrfToken string) ([]Company, error) {
	baseURL = strings.TrimRight(baseURL, "/")

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(baseURL)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}

	req, _ := http.NewRequest("POST", baseURL+"/support/service/storm/ajax-get-switch-companyList", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if csrfToken != "" {
		req.Header.Set("X-CSRF-TOKEN", csrfToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取公司列表失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var clr companyListResponse
	if err := json.Unmarshal(bodyBytes, &clr); err != nil {
		return nil, fmt.Errorf("解析公司列表响应: %w", err)
	}
	if clr.Code != 0 {
		return nil, fmt.Errorf("获取公司列表失败: %s (code=%d)", clr.Message, clr.Code)
	}

	var companies []Company
	for _, ho := range clr.Data.ItemList {
		for _, co := range ho.ItemList {
			if co.Switch {
				companies = append(companies, Company{ID: co.ID, Name: co.Name})
			}
		}
	}
	return companies, nil
}

// SwitchCompany switches to the specified company and returns updated session cookies.
func SwitchCompany(baseURL string, cookies []*http.Cookie, csrfToken string, targetID string) ([]*http.Cookie, error) {
	baseURL = strings.TrimRight(baseURL, "/")

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(baseURL)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}

	body := url.Values{}
	body.Set("targetId", targetID)

	req, _ := http.NewRequest("POST", baseURL+"/account-center/service/sso/ajax-change-login", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if csrfToken != "" {
		req.Header.Set("X-CSRF-TOKEN", csrfToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("切换公司失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var scr switchCompanyResponse
	if err := json.Unmarshal(bodyBytes, &scr); err != nil {
		return nil, fmt.Errorf("解析切换公司响应: %w", err)
	}
	if scr.Code != 0 {
		return nil, fmt.Errorf("切换公司失败: %s (code=%d)", scr.Message, scr.Code)
	}

	return jar.Cookies(u), nil
}
