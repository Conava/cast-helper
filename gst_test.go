package main

import (
	"slices"
	"testing"
	"time"
)

// TestEncoderControl needs GStreamer and a VA encoder; it checks that a key
// frame can be forced and the bitrate changed while the pipeline runs.
func TestEncoderControl(t *testing.T) {
	p, err := startGstPipeline("videotestsrc is-live=true pattern=ball ! video/x-raw,width=640,height=360,framerate=30/1 ! " +
		"vapostproc ! vah264enc name=enc aud=true b-frames=0 key-int-max=300 rate-control=cbr bitrate=4000 ! " +
		"h264parse config-interval=-1 ! video/x-h264,stream-format=byte-stream,alignment=au ! appsink name=out sync=false")
	if err != nil {
		t.Skip(err)
	}
	defer p.Stop()

	type au struct {
		at   time.Time
		size int
		key  bool
	}
	aus := make(chan au, 1000)
	go readAccessUnits(p, nalCodec{}, func(b []byte, key bool) { aus <- au{time.Now(), len(b), key} })
	collect := func(d time.Duration) (keys, bytes int) {
		end := time.After(d)
		for {
			select {
			case a := <-aus:
				if a.key {
					keys++
				}
				bytes += a.size
			case <-end:
				return
			}
		}
	}

	collect(time.Second) // skip the first key frame and the encoder warming up
	asked := time.Now()
	p.KeyFrame()
	for a := range aus {
		if a.key {
			t.Logf("forced key frame after %v", a.at.Sub(asked).Round(time.Millisecond))
			break
		}
		if a.at.Sub(asked) > time.Second {
			t.Fatal("no key frame within 1 s of asking (GOP is 10 s)")
		}
	}

	_, high := collect(2 * time.Second)
	p.SetBitrate(500)
	collect(500 * time.Millisecond)
	_, low := collect(2 * time.Second)
	t.Logf("2 s at 4000 kbit/s: %d bytes, at 500 kbit/s: %d bytes", high, low)
	if low*3 > high {
		t.Fatal("output did not shrink after lowering the bitrate")
	}
}

// TestPickVAElements covers the element names GStreamer gives VA devices: the
// first device plain, the others by render node, and low-power encoders.
func TestPickVAElements(t *testing.T) {
	nodes := []string{"/dev/dri/renderD128", "/dev/dri/renderD129"}
	cases := []struct {
		name, codec string
		elements    []string
		pp, enc     string
		ok          bool
	}{
		{"one GPU", "h264", []string{"vapostproc", "vah264enc"}, "vapostproc", "vah264enc", true},
		{"low power only", "hevc", []string{"vapostproc", "vah265lpenc"}, "vapostproc", "vah265lpenc", true},
		{"normal over low power", "h264", []string{"vapostproc", "vah264enc", "vah264lpenc"}, "vapostproc", "vah264enc", true},
		// An NVIDIA GPU on nouveau decodes but can't encode; the Intel one is second.
		{"second GPU encodes", "h264", []string{"vapostproc", "vah264dec", "varenderD129postproc", "varenderD129h264lpenc"},
			"varenderD129postproc", "varenderD129h264lpenc", true},
		{"encoder without its postproc", "h264", []string{"varenderD129h264enc"}, "vapostproc", "vah264enc", false},
		{"nothing", "hevc", nil, "vapostproc", "vah265enc", false},
	}
	for _, c := range cases {
		has := func(n string) bool { return slices.Contains(c.elements, n) }
		pp, enc, ok := pickVAElements(c.codec, nodes, has)
		if pp != c.pp || enc != c.enc || ok != c.ok {
			t.Errorf("%s: got %s %s %v, want %s %s %v", c.name, pp, enc, ok, c.pp, c.enc, c.ok)
		}
	}
}
