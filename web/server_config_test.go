package main

import "testing"

func TestResolveServerAddrUsesContainerDefaults(t *testing.T) {
	t.Setenv("ENV_TDX_API_HOST", "")
	t.Setenv("ENV_TDX_API_PORT", "")

	originalDetector := containerDetector
	containerDetector = func() bool { return true }
	t.Cleanup(func() {
		containerDetector = originalDetector
	})

	addr := resolveServerAddr()
	if addr != "0.0.0.0:8080" {
		t.Fatalf("unexpected addr: got %q want %q", addr, "0.0.0.0:8080")
	}
}

func TestResolveServerAddrUsesDefaults(t *testing.T) {
	t.Setenv("ENV_TDX_API_HOST", "")
	t.Setenv("ENV_TDX_API_PORT", "")

	originalDetector := containerDetector
	containerDetector = func() bool { return false }
	t.Cleanup(func() {
		containerDetector = originalDetector
	})

	addr := resolveServerAddr()
	if addr != "localhost:8080" {
		t.Fatalf("unexpected addr: got %q want %q", addr, "localhost:8080")
	}
}

func TestResolveServerAddrUsesConfiguredHostAndPort(t *testing.T) {
	t.Setenv("ENV_TDX_API_HOST", "127.0.0.1")
	t.Setenv("ENV_TDX_API_PORT", "9090")

	originalDetector := containerDetector
	containerDetector = func() bool { return true }
	t.Cleanup(func() {
		containerDetector = originalDetector
	})

	addr := resolveServerAddr()
	if addr != "127.0.0.1:9090" {
		t.Fatalf("unexpected addr: got %q want %q", addr, "127.0.0.1:9090")
	}
}

func TestResolveServerAddrFallsBackWhenPortInvalid(t *testing.T) {
	t.Setenv("ENV_TDX_API_HOST", "0.0.0.0")
	t.Setenv("ENV_TDX_API_PORT", "abc")

	originalDetector := containerDetector
	containerDetector = func() bool { return false }
	t.Cleanup(func() {
		containerDetector = originalDetector
	})

	addr := resolveServerAddr()
	if addr != "0.0.0.0:8080" {
		t.Fatalf("unexpected addr: got %q want %q", addr, "0.0.0.0:8080")
	}
}
