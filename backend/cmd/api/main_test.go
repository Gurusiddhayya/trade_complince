package main

import "testing"

func TestAllProducts(t *testing.T) {
	if len(products) != 13 {
		t.Fatalf("expected 13 products, got %d", len(products))
	}
	s := newStore()
	if len(s.Transactions) != 13 {
		t.Fatalf("expected 13 sample transactions, got %d", len(s.Transactions))
	}
	if s.Transactions["1"].RealizedValue != 60000 {
		t.Fatalf("unexpected export realization")
	}
}

func TestSecurityTokenParsingAndMode(t *testing.T) {
	t.Setenv("TCC_SECURITY_MODE", "production")
	t.Setenv("TCC_AUTH_TOKENS", "secret|u1|c1|customer")
	ps := configuredPrincipals()
	if len(ps) != 1 || ps[0].principal.CompanyID != "c1" || ps[0].principal.Role != "customer" {
		t.Fatalf("unexpected security principal configuration: %#v", ps)
	}
	if !authRequired() {
		t.Fatal("production mode must require authentication")
	}
}

func TestRateLimiter(t *testing.T) {
	l := newRateLimiter(2)
	if !l.allow("x") || !l.allow("x") || l.allow("x") {
		t.Fatal("rate limiter did not enforce configured limit")
	}
}
