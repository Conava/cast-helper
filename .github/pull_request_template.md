## Description

<!-- What does this PR change, and why? -->

## Type of change

- [ ] Bug fix
- [ ] New feature
- [ ] Protocol change (new commands, events or fields)
- [ ] New dependency or distribution support
- [ ] Documentation
- [ ] Other

## Related issues

<!-- "Closes #123" closes the issue when this is merged; "Refs #123" only links it. Leave empty if there's none. -->

## Screenshots / recordings

<!-- Required when the change is visible in a client (the Cast plugin, Dux Player) or on the TV. -->

## How I tested it

<!-- Distribution, GPU, cast devices, and the client you used (e.g. the Cast plugin with `./dev reload Cast` in dux-plugins). -->

## Checklist

- [ ] `gofmt`, `go vet ./...` and `go test ./...` pass
- [ ] Added or updated tests for new logic
- [ ] Tested end to end: cast to a real device from a client, not just the tests
- [ ] Protocol changes are additive, documented in `PROTOCOL.md`, and bump `protocolVersion` in `main.go`
- [ ] New dependencies are in `deps.go` with package names for every distro family, and in the README
- [ ] Updated the README if setup, dependencies or features changed
- [ ] Screenshots attached for visible changes
- [ ] Changes that also need the Cast plugin link the matching PR in [dux-plugins](https://github.com/Conava/dux-plugins)
- [ ] If AI wrote a meaningful part, I've said so here, and I've read, tested and can explain every line
