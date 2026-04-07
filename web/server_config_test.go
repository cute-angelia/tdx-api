package main

import "testing"

func TestResolveServerAddrUsesDefaults(t *testing.T) {
	t.Setenv("ENV_TDX_API_HOST", "")
	t.Setenv("ENV_TDX_API_PORT", "")

	addr := resolveServerAddr()
	if addr != "localhost:8080" {
		t.Fatalf("unexpected addr: got %q want %q", addr, "localhost:8080")
	}
}

func TestResolveServerAddrUsesConfiguredHostAndPort(t *testing.T) {
	t.Setenv("ENV_TDX_API_HOST", "127.0.0.1")
	t.Setenv("ENV_TDX_API_PORT", "9090")

	addr := resolveServerAddr()
	if addr != "127.0.0.1:9090" {
		t.Fatalf("unexpected addr: got %q want %q", addr, "127.0.0.1:9090")
	}
}

func TestResolveServerAddrFallsBackWhenPortInvalid(t *testing.T) {
	t.Setenv("ENV_TDX_API_HOST", "0.0.0.0")
	t.Setenv("ENV_TDX_API_PORT", "abc")

	addr := resolveServerAddr()
	if addr != "0.0.0.0:8080" {
		t.Fatalf("unexpected addr: got %q want %q", addr, "0.0.0.0:8080")
	}
}
