package sav

import "testing"

func TestPlayerSaveNameRejectsCompanionFiles(t *testing.T) {
	for _, tc := range []struct {
		name string
		want bool
	}{
		{"3C22D58D000000000000000000000000.sav", true},
		{"3c22d58d000000000000000000000000.SAV", true},
		{"3C22D58D000000000000000000000000_dps.sav", false},
		{"Level.sav", false},
		{"3C22D58D00000000000000000000000Z.sav", false},
	} {
		if got := isPlayerSaveName(tc.name); got != tc.want {
			t.Errorf("isPlayerSaveName(%q) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestPlayerSaveUIDPrefersEmbeddedGUID(t *testing.T) {
	data := propertyMap{"PlayerUId": {Type: "StructProperty", Value: structData{Type: "Guid", Value: "embedded-guid"}}}
	if got := playerSaveUID(data, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.sav"); got != "embedded-guid" {
		t.Fatalf("playerSaveUID = %q, want embedded GUID", got)
	}
	if got := playerSaveUID(propertyMap{}, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.sav"); got != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" {
		t.Fatalf("fallback playerSaveUID = %q", got)
	}
}

func TestPlayerSaveLocationReadsLastTransformTranslation(t *testing.T) {
	want := Vector{X: 12.5, Y: -8, Z: 300}
	data := propertyMap{"LastTransform": {
		Type: "StructProperty",
		Value: structData{Type: "Transform", Value: propertyMap{"Translation": {
			Type: "StructProperty", Value: structData{Type: "Vector", Value: want},
		}}},
	}}
	got, ok := playerSaveLocation(data)
	if !ok || got != want {
		t.Fatalf("playerSaveLocation = %+v, %v; want %+v, true", got, ok, want)
	}
}
