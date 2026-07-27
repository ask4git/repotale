package main

import "testing"

func TestSaveLoadConfigRoundTrip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	want := Config{Token: "ghp_test123", Repo: "acme/widgets"}
	if err := saveConfig(want); err != nil {
		t.Fatalf("saveConfig: %v", err)
	}

	got, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}

	if got != want {
		t.Fatalf("loadConfig = %+v, want %+v", got, want)
	}
}
