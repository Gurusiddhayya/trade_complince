package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type securityPrincipal struct {
	UserID    string
	CompanyID string
	Role      string
}

type securityContextKey string

const principalKey securityContextKey = "tcc-principal"

func principalFromContext(ctx context.Context) (securityPrincipal, bool) {
	p, ok := ctx.Value(principalKey).(securityPrincipal)
	return p, ok
}

// TCC_AUTH_TOKENS format: token|user_id|company_id|role,token2|user2|company2|role2
// This is a development/prototype authentication boundary. Production should replace
// static tokens with OIDC/OAuth/JWT validation without changing handler authorization checks.
func configuredPrincipals() []struct {
	token     string
	principal securityPrincipal
} {
	raw := strings.TrimSpace(os.Getenv("TCC_AUTH_TOKENS"))
	if raw == "" {
		return nil
	}
	var out []struct {
		token     string
		principal securityPrincipal
	}
	for _, item := range strings.Split(raw, ",") {
		p := strings.Split(item, "|")
		if len(p) != 4 || p[0] == "" {
			continue
		}
		out = append(out, struct {
			token     string
			principal securityPrincipal
		}{p[0], securityPrincipal{UserID: p[1], CompanyID: p[2], Role: strings.ToLower(p[3])}})
	}
	return out
}

func authRequired() bool {
	return strings.EqualFold(os.Getenv("TCC_SECURITY_MODE"), "production") || os.Getenv("TCC_AUTH_REQUIRED") == "1"
}

func isPublicPath(path string) bool {
	switch path {
	case "/api/v1/health", "/api/v1/products", "/api/v1/knowledge", "/api/v1/regulatory/updates":
		return true
	default:
		return false
	}
}

func authenticate(next http.Handler) http.Handler {
	principals := configuredPrincipals()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) && r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if !authRequired() && auth == "" {
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
			return
		}
		token := strings.TrimSpace(auth[len("Bearer "):])
		for _, item := range principals {
			if subtle.ConstantTimeCompare([]byte(token), []byte(item.token)) == 1 {
				ctx := context.WithValue(r.Context(), principalKey, item.principal)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid authentication token"})
	})
}

func requireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, role := range roles {
		allowed[strings.ToLower(role)] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := principalFromContext(r.Context())
			if !authRequired() && !ok {
				next.ServeHTTP(w, r)
				return
			}
			if !ok || !allowed[strings.ToLower(p.Role)] {
				writeJSON(w, http.StatusForbidden, map[string]string{"error": "insufficient role"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type fixedWindowLimiter struct {
	mu     sync.Mutex
	window time.Time
	counts map[string]int
	limit  int
}

func newRateLimiter(limit int) *fixedWindowLimiter {
	return &fixedWindowLimiter{counts: map[string]int{}, limit: limit}
}
func (l *fixedWindowLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if l.window.IsZero() || now.Sub(l.window) >= time.Minute {
		l.window = now
		l.counts = map[string]int{}
	}
	l.counts[key]++
	return l.counts[key] <= l.limit
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if r.TLS != nil || strings.EqualFold(os.Getenv("TCC_SECURITY_MODE"), "production") {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

func secureCORS(next http.Handler) http.Handler {
	origins := strings.TrimSpace(os.Getenv("TCC_ALLOWED_ORIGINS"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origins != "" {
			w.Header().Set("Access-Control-Allow-Origin", origins)
		} else if !authRequired() {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requestSecurity(rate *fixedWindowLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			w.Header().Set("X-Request-ID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
		}
		key := r.RemoteAddr
		if p, ok := principalFromContext(r.Context()); ok && p.UserID != "" {
			key = "user:" + p.UserID
		}
		if !rate.allow(key) {
			writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
			return
		}
		log.Printf("request method=%s path=%s remote=%s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func securityMode() string {
	if authRequired() {
		return "production"
	}
	return "development"
}

func securityConfigSummary() map[string]string {
	limit := os.Getenv("TCC_RATE_LIMIT_PER_MINUTE")
	if limit == "" {
		limit = "120"
	}
	if _, err := strconv.Atoi(limit); err != nil {
		limit = "120"
	}
	return map[string]string{"mode": securityMode(), "rate_limit_per_minute": limit, "authentication": "Bearer token boundary; replace with OIDC/OAuth/JWT for production", "tenant_scope": "principal carries company_id; handler-level scoping required", "cors": "configured via TCC_ALLOWED_ORIGINS in production"}
}
