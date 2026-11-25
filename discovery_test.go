package main

import "testing"

func TestParseTXTAndDedupe(t *testing.T) {
	cast := parseTXT("chromecast", "Chromecast-abc", [][]byte{[]byte("id=abc123"), []byte("md=Chromecast Ultra"), []byte("fn=Living Room"), []byte("rs=")})
	if cast.ID != "chromecast:abc123" || cast.Name != "Living Room" || cast.Model != "Chromecast Ultra" {
		t.Fatalf("cast: %+v", cast)
	}
	air := parseTXT("airplay", "Bedroom TV", [][]byte{[]byte("deviceid=AA:BB:CC:DD:EE:FF"), []byte("features=0x5A7FFFF7,0x1E"), []byte("model=AppleTV5,3")})
	if air.ID != "airplay:aa:bb:cc:dd:ee:ff" || air.Name != "Bedroom TV" || air.Model != "AppleTV5,3" || air.Features != "0x5A7FFFF7,0x1E" {
		t.Fatalf("airplay: %+v", air)
	}
	if noID := parseTXT("chromecast", "Foo", nil); noID.ID != "chromecast:Foo" {
		t.Fatalf("fallback id: %+v", noID)
	}

	r := newRegistry()
	eth := instance{iface: 2, name: "Chromecast-abc", stype: "_googlecast._tcp", domain: "local"}
	wifi := eth
	wifi.iface = 3
	r.pending(eth)
	r.pending(wifi)
	if !r.resolved(eth, cast) || r.resolved(wifi, cast) {
		t.Fatal("second interface must not add the device again")
	}
	if _, gone := r.removed(eth); gone {
		t.Fatal("device still on wifi")
	}
	if d, gone := r.removed(wifi); !gone || d.ID != cast.ID {
		t.Fatal("last instance gone must remove the device")
	}

	// removed while resolving: the late resolve result is dropped
	r.pending(eth)
	r.removed(eth)
	if r.resolved(eth, cast) {
		t.Fatal("stale resolve must not add")
	}
}

func TestCanMirror(t *testing.T) {
	for features, want := range map[string]bool{
		"0x5A7FFFF7,0x1E":       true,  // Apple TV
		"0x7F8AD0,0x38BCF46":    true,  // LG webOS TV
		"0x444f8a00":            false, // Yamaha receiver, audio only
		"0x4A7FCA00,0xBC354BD0": false, // speaker without bit 7
		"":                      true,  // unreadable: keep it
	} {
		if got := canMirror(features); got != want {
			t.Errorf("canMirror(%q) = %v, want %v", features, got, want)
		}
	}
}
