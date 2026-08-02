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

func TestTreasureMapPointGUIDKeyHint(t *testing.T) {
	value := &gw{}
	value.bytes(intProp("Rarity", 3))
	value.fstr("None")
	entries := &gw{}
	entries.guid(testGUID(0x71))
	entries.bytes(value.b)

	record := &gw{}
	record.bytes(mapProp("FoundTreasureMapPointMap", "StructProperty", "StructProperty", 1, entries.b))
	record.fstr("None")
	save := &gw{}
	save.bytes(structProp("RecordData", "PalPlayerRecordData", record.b))
	save.fstr("None")
	top := &gw{}
	top.bytes(structProp("SaveData", "PalPlayerSaveData", save.b))
	top.fstr("None")

	stats := newStats()
	props, err := readProperties(newReaderWithStats(top.b, &stats), "", &stats)
	if err != nil {
		t.Fatal(err)
	}
	saveProps, ok := propertyProperties(props, "SaveData")
	if !ok {
		t.Fatal("SaveData was not decoded")
	}
	recordProps, ok := propertyProperties(saveProps, "RecordData")
	if !ok {
		t.Fatal("RecordData was not decoded")
	}
	treasureMap := recordProps["FoundTreasureMapPointMap"]
	if treasureMap == nil {
		t.Fatal("FoundTreasureMapPointMap was not decoded")
	}
	points, ok := treasureMap.Value.([]mapEntry)
	if !ok || len(points) != 1 {
		t.Fatalf("treasure-map points = %#v, want one entry", treasureMap.Value)
	}
	point, ok := asProperties(points[0].Value)
	if !ok || firstInt(point, "Rarity") != 3 {
		t.Fatalf("treasure-map point value = %#v", points[0].Value)
	}
	if stats.SkippedProperties != 0 || stats.SkippedStructs != 0 {
		t.Fatalf("unexpected skips: %+v", stats)
	}
}
