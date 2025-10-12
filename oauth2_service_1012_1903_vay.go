// 代码生成时间: 2025-10-12 19:03:58
package main

import (
    "fmt"
    "net/http"
    "golang.org/x/oauth2"
    "github.com/gorilla/sessions"
    "log"
)

// OAuth2Config 配置OAuth2认证
type OAuth2Config struct {
    RedirectURL string
    ClientID    string
    ClientSecret string
    Scopes     []string
    AuthURL    string
    TokenURL   string
}

// OAuth2Service 提供OAuth2认证服务
type OAuth2Service struct {
    config *OAuth2Config
    oauth2Config *oauth2.Config
    store *sessions.CookieStore
}

// NewOAuth2Service 创建一个新的OAuth2Service实例
func NewOAuth2Service(config *OAuth2Config) *OAuth2Service {
    oAuth2Config := &oauth2.Config{
        ClientID:     config.ClientID,
        ClientSecret: config.ClientSecret,
        RedirectURL:  config.RedirectURL,
        Scopes:       config.Scopes,
        Endpoint: oauth2.Endpoint{
            AuthURL:  config.AuthURL,
            TokenURL: config.TokenURL,
        },
    }
    store := sessions.NewCookieStore([]byte("your_secret_key_here"))
    return &OAuth2Service{
        config:        config,
        oauth2Config: oAuth2Config,
        store:        store,
    }
}

// AuthHandler 处理OAuth2认证请求
func (s *OAuth2Service) AuthHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        state := generateRandomStateString()
        s.store.Save(r, w, sessions.NewSession(r, state))
        authURL := s.oauth2Config.AuthCodeURL(state)
        http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
    } else {
        http.Error(w, "Invalid request", http.StatusBadRequest)
    }
}

// CallbackHandler 处理OAuth2回调请求
func (s *OAuth2Service) CallbackHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        state := r.URL.Query().Get("state")
        session, _ := s.store.Get(r, state)
        if session == nil {
            http.Error(w, "Invalid state", http.StatusBadRequest)
            return
        }
        code := r.URL.Query().Get("code")
        if code == "" {
            http.Error(w, "No code provided", http.StatusBadRequest)
            return
        }
        token, err := s.oauth2Config.Exchange(r.Context(), code)
        if err != nil {
            log.Printf("Failed to exchange code for token: %v", err)
            http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "Token received: %+v", token)
    } else {
        http.Error(w, "Invalid request", http.StatusBadRequest)
    }
}

// generateRandomStateString 生成随机状态字符串
func generateRandomStateString() string {
    // Implement your own state string generation logic here
    return "random_state_string"
}

func main() {
    config := &OAuth2Config{
        RedirectURL:  "http://localhost:8080/callback",
        ClientID:    "your_client_id",
        ClientSecret: "your_client_secret",
        Scopes:      []string{"email", "profile"},
        AuthURL:     "https://provider.com/oauth/authorize",
        TokenURL:    "https://provider.com/oauth/token",
    }
    service := NewOAuth2Service(config)

    http.HandleFunc("/auth", service.AuthHandler)
    http.HandleFunc("/callback", service.CallbackHandler)

    log.Println("OAuth2 service listening on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start OAuth2 service: %v", err)
    }
}
