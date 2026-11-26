# Contributing

Whether cast-helper works depends on the GPU and its driver, the cast device and the distribution.
I can only test a small part of that myself, with an AMD system and the few distributions I can reasonably test by hand.
Testers with other hardware or software are more than welcome: an Intel or NVIDIA GPU, another TV brand, an Apple TV, Debian, Fedora, NixOS or alike.
Please report anything that doesn't work on your machine, and try the fix branches created when you report hardware-specific issues.

## Before you start

- Bigger changes: open an issue before a new feature, a new dependency or a refactoring, so we can agree on the approach before you spend time on it. Small features or fixes can go straight to a pull request.
- Plans: the [plans and development goals](README.md#plans-and-development-goals) in the README are where help is most useful.
- AI: using AI tools to write code is fine, but you are responsible for the result. You have to understand every change, have tested it and be able to explain and discuss it. Mention in the PR if a meaningful part was generated. AI-generated issues, PR descriptions or comments are strictly forbidden and will be closed without discussion. If writing it yourself isn't worth your time, reading it isn't worth mine.
- [AGENTS.md](AGENTS.md) is the guide for AI agents. Make sure that any AI agent you use follows it, and that you understand what it does. This exists to help improve the quality of generated contributions, but you remain solely responsible for the result.

## Setup

You need Go 1.25 or newer and the dependencies from the [README](README.md#installation).

```bash
git clone https://github.com/Conava/cast-helper.git && cd cast-helper
go test ./...
CGO_ENABLED=0 go build -o cast-helper .
```

Without arguments, a build reads commands from stdin. That is the quickest way to try it:

```bash
echo '{"cmd": "check"}' | timeout 5 ./cast-helper
```

To test with a real client, check out [dux-plugins](https://github.com/Conava/dux-plugins) next to this repository and run `./dev reload Cast` there.
It builds this checkout with detailed logging, and the Cast plugin uses that build instead of the installed one.
The log is in `~/.cache/dms/cast/serve.log`.

## Making a change

- Scope: one fix or feature per PR.
- Checks: run `gofmt`, `go vet ./...` and `go test ./...` before you push. CI fails if `go vet` or the tests fail.
- Testing: write tests for new logic. Tests can't connect to an actual device, so cast to a real device before you open the PR, and mention the device, GPU and distribution you used.
- No cgo: the binary is built without cgo and loads GStreamer at runtime. Don't add anything that needs cgo.
- Protocol: changes only add commands, events or fields. Nothing gets changed or removed unless discussed before. Document them in [PROTOCOL.md](PROTOCOL.md) and raise `protocolVersion` in `main.go`.
- Dependencies: a new system package goes into the table in `deps.go`, with its name for every distribution family, and into the install commands in the README.
- Docs: when installation, dependencies or features change, update the README in the same PR.
- Commits: `cast-helper: use VA-API low-power encoders`, lowercase after the prefix. CI and repository changes use `ci:` or `repo:`.
- Cast plugin: if the plugin needs a matching change, open that PR in dux-plugins and link the two PRs.

## Reporting bugs

Use the [bug report form](https://github.com/Conava/cast-helper/issues/new/choose).

- With the Cast plugin: turn on *Detailed logging* under Troubleshooting in its settings, reproduce the problem and paste what *Copy report* gives you.
- Without the plugin: attach `~/.cache/dms/cast/serve.log` and the output of `echo '{"cmd": "check"}' | timeout 5 cast-helper`.

Logs can contain device names, local addresses and file names, so please check them before you post.
