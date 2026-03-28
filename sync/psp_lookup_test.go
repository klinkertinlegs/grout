package sync

import "testing"

func TestLookupPSPTitle_WithSuffix(t *testing.T) {
	title, ok := LookupPSPTitle("UCUS98662_GameData0")
	if !ok {
		t.Fatal("expected to find LocoRoco")
	}
	if title != "LocoRoco" {
		t.Errorf("got %q, want 'LocoRoco'", title)
	}
}

func TestLookupPSPTitle_PlainID(t *testing.T) {
	title, ok := LookupPSPTitle("UCUS98662")
	if !ok {
		t.Fatal("expected to find LocoRoco")
	}
	if title != "LocoRoco" {
		t.Errorf("got %q, want 'LocoRoco'", title)
	}
}

func TestLookupPSPTitle_WithDashes(t *testing.T) {
	title, ok := LookupPSPTitle("UCUS-98662_GameData0")
	if !ok {
		t.Fatal("expected to find LocoRoco with dashes")
	}
	if title != "LocoRoco" {
		t.Errorf("got %q, want 'LocoRoco'", title)
	}
}

func TestLookupPSPTitle_NotFound(t *testing.T) {
	_, ok := LookupPSPTitle("ZZZZ99999_GameData0")
	if ok {
		t.Error("expected not found for fake game ID")
	}
}

func TestLookupPSPTitle_EmptyString(t *testing.T) {
	_, ok := LookupPSPTitle("")
	if ok {
		t.Error("expected not found for empty string")
	}
}
