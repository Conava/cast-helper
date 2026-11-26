# cast-helper protocol

This describes how an app casts through cast-helper.
The Cast plugin and Dux Player both use it, and any app that can start a process and read its output line by line can do the same.

Each direction carries one JSON object per line. Clients send commands and the helper sends events.
The helper pushes the full state on every change, so a client never has to poll or keep its own state.

## Trying it

```console
$ cast-helper attach
{"devices":[],"event":"state","noCast":null,"session":null}
{"cmd": "discover", "on": true}
{"devices":[{"id":"6f3c…","name":"Living Room TV","model":"Chromecast","protocol":"chromecast","protocols":["chromecast"]}],"event":"state","noCast":null,"session":null}
{"cmd": "play", "device": "6f3c…", "path": "/home/me/Videos/holiday.mp4"}
```

Every following `state` event shows the session moving from `starting` to `casting`, including the position.
`{"cmd": "stop"}` ends the cast.

## Writing a client

1. Find the binary: `cast-helper` on `PATH`, otherwise `~/.local/share/dms/cast/bin/cast-helper`. If neither exists, point the user to the [installation guide](README.md#installation).
2. Run `cast-helper attach` and keep it running while your app is open.
3. Send `hello` with the id and name of your app, so other clients can show who started a cast.
4. Send `check`. Compare its `protocol` with the version your client needs (see [Versions](#versions)) and show `missing` to the user, because nothing works without those packages.
5. Build your UI from the latest `state` event only. It always contains everything.
6. Turn `discover` on while a device list is shown, and off again when it closes.
7. Show your own text for each error `code`. `detail` is meant for logs.

## Connecting

- `cast-helper attach`: what a client runs. It connects its stdin and stdout to the shared helper on `$XDG_RUNTIME_DIR/cast-helper.sock` and starts `cast-helper serve` first if none is running. It exits with 0 when its stdin closes and with 1 when the helper is gone. Clients can also connect to the socket directly.
- Shared helper: sends every client a `state` event right after it connects, then one on every change. A cast keeps running when the client that started it disconnects. The helper exits a few seconds after the last client is gone and nothing is casting.
- Plain `cast-helper`: speaks the same protocol on its own stdin and stdout for a single client, and stops its cast when stdin closes.
- Replies to `check`, `thumbnail` and `playlist`, and errors about the command itself (`busy`, `unsupportedLink`), only go to the client that sent the command. Everything else goes to all clients.

## Versions

`check` reports the `protocol` version.
Changes are additive: commands and fields get added, but never changed or removed.
A client works with any helper whose `protocol` is at least the version it was written for, and should suggest an update otherwise.
Helpers that don't report a version speak protocol 1.

| Protocol | Added |
|---|---|
| 1 | Everything below not marked otherwise |
| 2 | `serve` and `attach`, `hello`, `session.owner`, `title`/`artist`/`album`/`art` on items, `version` in `check` |
| 3 | `edit`, `repeat` and `times` for a running queue |
| 4 | `noCast` in `state`, `debug` |

## Other subcommands

```bash
cast-helper version
cast-helper install-missing [all|required]   # default all
```

`install-missing` installs everything `check` reports as missing through the distribution's package manager, or through an AUR helper for AUR packages.
It is meant to run in a terminal, so it can ask for the password and for confirmation.
A client that offers an install button should open a terminal with it.

## Commands

```jsonc
{"cmd": "hello", "client": "dplayer", "name": "Dux Player"}   // 2: names this client for session.owner
{"cmd": "check"}
{"cmd": "discover", "on": true}                                // while a device list is on screen
{"cmd": "start", "device": "<id>", "audio": "tv|both|local", "width": 1920, "height": 1080, "fps": 30,
 "kbps": 0 /* auto */, "codec": "auto|h264|hevc", "delay": 200, "restoreToken": "…",
 "protocol": "chromecast|airplay" /* where a device has both */, "airplayPorts": "47840-47849"}
{"cmd": "pair", "pin": "1234"}                                 // AirPlay PIN or password
{"cmd": "play", "device": "<id>", "path": "/abs/movie.mkv" | "url": "https://…", "port": 47830,
 "title"?, "artist"?, "album"?, "art"?}
{"cmd": "queue", "device": "<id>", "items": [{"path": "…"} | {"url": "…"}, …], "index": 0,
 "interval": 8, "loop": false, "shrinkPhotos": true, "convertVideos": true, "port": 47830,
 "repeat"?: "show|item", "times"?: 1 /* 3: 0 until stopped */}
{"cmd": "pause"} {"cmd": "resume"} {"cmd": "seek", "position": 30} {"cmd": "subtitles", "on": false}
{"cmd": "next"} {"cmd": "previous"} {"cmd": "jump", "index": 3} {"cmd": "slideshow", "interval": 0}
{"cmd": "edit", "items": [{"path": "…"} | {"url": "…"}, …], "index": 0}  // 3: replaces the running queue
{"cmd": "repeat", "repeat": "show|item|", "times": 0}           // 3: for the running queue
{"cmd": "move", "device": "<id>"}                              // the same cast on another device
{"cmd": "thumbnail", "path": "/abs/video.mp4", "size": 256}    // a song's thumbnail is its cover
{"cmd": "playlist", "path": "/abs/list.m3u8"}                  // .m3u, .m3u8, .pls
{"cmd": "stop"}
{"cmd": "debug", "on": true}                                   // 4: detailed logging in serve.log
{"cmd": "volume", "level": 0.4}                                // receiver volume, where supported
```

Queue items and `play` take optional metadata (protocol 2): `title`, `artist`, `album`, and `art`, which is an http(s) link or the absolute path of a cover image.
The device shows it, and the session reports it instead of the file name or link.
The tags of a song file fill in whatever isn't given.
A playlist file in `items` is expanded in place.

`edit` (protocol 3) replaces the queue of a running cast, for example when a show is edited while it plays.
Items that were already in the queue keep what was prepared for them.
If the item at `index` is the one playing, it keeps playing. Otherwise the device switches to the item at `index`.

`repeat` (protocol 3) sets what plays again when it ends on its own: the whole queue (`show`) or the current item (`item`), either `times` more times or until stopped with `times` 0.
An empty `repeat` turns it off.
`loop: true` is the same as `show` until stopped.
Next and previous always move on.

## Events

```jsonc
{"event": "state", "devices": [{"id", "name", "model", "protocol": "chromecast|airplay", "protocols": […]}],
 "noCast": ["Smart TV"],  // 4: Android TVs on the network with Cast turned off or missing
 "session": null | {"kind": "mirror|media", "device", "name",
   "phase": "starting|authorizing|pairing|converting|casting", "since", "volume", "canVolume",
   "owner"?: {"id": "dplayer", "name": "Dux Player"},   // 2: the client that started it, after its hello
   "protocol"?: "airplay", "credential"?: "pin|password" /* while pairing */,
   // media only:
   "title", "progress", "playerState", "position", "positionAt", "duration", "rate",
   "subtitles", "subtitlesOn", "index", "interval", "paused", "loop", "repeat", "times" /* 3 */,
   "queue": [{"source" /* 3: path or link */, "title", "artist", "album", "art"? /* links only */, "kind": "image|video|audio|link|youtube",
              "status": "waiting|preparing|ready|failed", "progress", "error"?, "feature"?, "missing"?}]}}
{"event": "check", "protocol": 4, "version": "0.2.0", "distro": "arch|debian|fedora|opensuse|other",
 "missing": [required ids], "deps": [{"id", "package", "required", "present", "features": […]}]}
{"event": "restoreToken", "token": "…"}                      // keep it for the next mirroring start
{"event": "thumbnail", "path": "…", "file": "~/.cache/dms/cast/thumbs/<hash>.jpg" | ""}
{"event": "playlist", "path": "…", "items": [{"path"?, "url"?, "title"?}]}
{"event": "error", "code": "…", "detail": "…", "feature"?: "mirroring", "missing"?: ["gst-plugin-va"]}
{"event": "warning", "code": "soundOnBoth", "feature": "tvOnlySound", "missing": ["pipewire"]}
{"event": "warning", "code": "airplaySound"}                 // AirPlay sound stays on this computer
```

`detail` is English text for logs. Clients show their own text for each `code`: `airplayFailed`, `airplayNoMedia`, `avahi`, `busy`, `deviceLost`, `deviceNotFound`, `encoderMissing`, `launchRefused`, `loadFailed`, `missingDependency`, `pipelineFailed`, `portalDenied`, `portInUse`, `receiverStopped`, `serverUnreachable`, `unsupportedFile`, `unsupportedLink`, `youtubeFailed`, `youtubeInQueue`.

## Behavior

- Dependencies: none are needed to start. Each feature checks its packages when it is used and fails with `missingDependency`, naming the missing ones. The features are `discovery`, `mirroring`, `sound`, `tvOnlySound`, `files`, `mkv`, `convert`, `photos`, `heic` and `thumbnails`.
- Queue: a single `play` is a queue with one item. The current item and the next video or next two photos are prepared ahead (served, remuxed, converted, shrunk). Photos advance on the interval, everything else when it ends. A YouTube link plays on its own, in the YouTube app of the device.
- Devices: a TV with Chromecast and AirPlay is one device and keeps its Chromecast id. Mirroring uses Chromecast unless `protocol` asks for AirPlay. Files, links and shows always use Chromecast.
- Position: the receiver only reports it when it changes. While playing, count on from `positionAt`.
- MPRIS: a cast shows up as `org.mpris.MediaPlayer2.cast`. File playback supports play, pause, next, previous and seek. For mirroring, pause and stop end the cast.
- Debugging: `CAST_DEBUG=1` logs every Cast message and HTTP request. The shared helper logs to `~/.cache/dms/cast/serve.log`.
