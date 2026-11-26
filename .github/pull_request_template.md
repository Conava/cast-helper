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

<!-- "Closes #123" closes the issue on merge, "Refs #123" only links it. Leave empty if there is none. -->

## Screenshots / recordings

<!-- Required if the change is visible in a client (Cast plugin, Dux Player) or on the cast device. -->

## Testing

<!-- Distribution, GPU, cast devices and the client you used, for example the Cast plugin with `./dev reload Cast` in dux-plugins. -->

## Checklist

- [ ] `gofmt`, `go vet ./...` and `go test ./...` pass
- [ ] Tests added or updated for new logic
- [ ] Tested end to end: cast to a real device from a client, not only the tests
- [ ] Protocol changes are additive, documented in `PROTOCOL.md` and raise `protocolVersion` in `main.go`
- [ ] New dependencies are in `deps.go` with package names for every distribution family, and in the README
- [ ] README updated if setup, dependencies or features changed
- [ ] Screenshots attached for visible changes
- [ ] Changes that also need the Cast plugin link the matching PR in [dux-plugins](https://github.com/Conava/dux-plugins)
- [ ] If AI wrote a meaningful part, it is mentioned here, and I have read, tested and can explain every line
