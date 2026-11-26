# cast-helper

A Go service that casts from Linux to Chromecast and AirPlay devices.
Clients (the Cast plugin in Dux Shell, Dux Player) start `cast-helper attach` and exchange JSON lines with one shared helper per user.
[PROTOCOL.md](PROTOCOL.md) is the contract with those clients.
The README is for users, so keep developer details out of it.

## Code map

Everything is one `main` package without subpackages.

| File | Description |
|---|---|
| `main.go` | Subcommands, the command switch, `version` and `protocolVersion`, the `state` event |
| `serve.go` | The shared helper: socket, lock, `attach`, sending events to all clients |
| `glue.go` | Discovery loop, merging the Chromecast and AirPlay entries of a TV, debug logging |
| `discovery.go` | Avahi over D-Bus, TXT records to devices |
| `castv2.go` | Cast v2 control channel (TLS, protobuf framing, JSON payloads) |
| `session.go`, `sender.go`, `source.go`, `vui.go` | Chromecast mirroring: portal, GStreamer pipelines, Cast Streaming RTP |
| `portal.go`, `audio.go` | ScreenCast portal, PipeWire sound routing |
| `airplay.go` | AirPlay mirroring through the control socket of doubletake |
| `queue.go`, `media.go`, `music.go` | Files, links and queues on the media player of the receiver |
| `convert.go`, `remux.go`, `thumbs.go` | GPU conversion, MKV remuxing, thumbnails |
| `youtube.go` | YouTube links through the YouTube app of the TV |
| `mpris.go` | The cast as an MPRIS player |
| `gst.go` | GStreamer loaded at runtime through purego |
| `deps.go` | The dependency table, distro detection, `check` and `install-missing` |

## Rules

- No cgo: releases are static `CGO_ENABLED=0` builds, and GStreamer is loaded at runtime with purego. New native calls follow the pattern in `gst.go`.
- Protocol: only add commands, events and fields. Never rename, retype or remove one unless it was agreed in an issue first. Each addition raises `protocolVersion` in `main.go` and gets a row in the Versions table of PROTOCOL.md and a `// N:` marker next to it.
- Dependencies: nothing is required to start. Each feature checks its packages when it is used and fails with `missingDependency`. A new system package goes into `dependencies` in `deps.go`, with its name for arch, debian, fedora and opensuse checked against their package indexes, and into the install commands in the README.
- Network and files: the README section of that name lists every host contacted and every path written. A change to either updates it.
- Errors: each carries a stable `code`, which clients translate, and an English `detail` for logs. A new code goes into the list in PROTOCOL.md.
- Comments: say why, in a line or two. `sender.go`, `youtube.go` and `airplay.go` open with a comment explaining their mechanism. Keep those accurate when the code changes.

## Verifying

```bash
gofmt -l . ; go vet ./... ; go test ./...
CGO_ENABLED=0 go build -o cast-helper .
echo '{"cmd": "check"}' | timeout 5 ./cast-helper
```

Tests can't replace a real cast.
For anything that reaches a device, build it through the dev script of the Cast plugin (`./dev reload Cast` in a dux-plugins checkout next to this one) and hand over the steps for a manual test on real hardware: what to cast, to which kind of device, and what the TV should show.
`-ldflags "-X main.fakeMissing=a,b"` makes a build report those dependency ids as missing, to test how a client handles it.

## Shipping

- Commits: `cast-helper: short lowercase summary`, or `ci:` and `repo:` for everything else.
- Releases: a `vX.Y.Z` tag that matches `const version` in `main.go`. CI then builds `cast-helper-x86_64`, `cast-helper-aarch64` and `SHA256SUMS` and publishes the release. The asset names follow `uname -m`, because the install commands and the updater of the Cast plugin download them by that name.
- AUR: after a release, `cast-helper` and `cast-helper-bin` are bumped to it. `cast-helper-git` builds `master` and needs nothing.
