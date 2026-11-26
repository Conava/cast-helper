# cast-helper

Cast your Linux desktop, files, music and links to Chromecast, Google TV and AirPlay devices.

cast-helper is developed to serve the backend of [Cast plugin](https://github.com/Conava/dux-plugins/tree/master/Cast)
in Dux Shell and the casting in [Dux Player](https://github.com/Conava/dux-player), but it has an open protocol and is compatible with any other app.
You install it once, and every app that uses it shares the same cast, so you can control the music cast from dux-player in your Cast plugin widget.

## Features

- Screen mirroring incl. sound, encoded on the GPU in real time. Chromecast gets up to 1080p at 60 fps. AirPlay receivers such as an Apple TV work too.
- Sound on the cast device, on your computer, or on both.
- Videos, photos and music on Chromecast devices, with title, artist and cover on the TV. Incompatible files like 8K video, AV1 or HEIC photos are converted on the fly.
- Queues of files and links, photo shows on a timer, and `.m3u` and `.pls` playlists.
- Support for direct media links on the web and YouTube (Controls the cast devices youtube app)
- MPRIS support, so a cast shows up in media widgets and answers the media keys like any other player.

## Supported devices

| Device | Mirroring | Files, music, links |
|---|---|---|
| Chromecast, Google TV, TVs with Chromecast built-in | Yes | Yes |
| Apple TV and TVs with AirPlay | Yes, through doubletake | Not yet |
| HomePod and other AirPlay speakers | No | Not yet |

A TV that speaks both Chromecast and AirPlay allows to select the connection protocol with Chromecast as the default. AirPlay has fewer features, but may offer better quality and lower latency. 

## Installation

cast-helper runs on Linux, on x86_64 and aarch64. 
It finds devices through Avahi and uses GStreamer to encode and convert, and mirroring requires a GPU driver that can encode H.264 or HEVC through VA-API.

### Arch Linux

cast-helper is packaged on the AUR.

| Package | Description |
|---|---|
| [`cast-helper`](https://aur.archlinux.org/packages/cast-helper) | The latest release, built on your machine |
| [`cast-helper-bin`](https://aur.archlinux.org/packages/cast-helper-bin) | Prebuilt binary of the latest release |
| [`cast-helper-git`](https://aur.archlinux.org/packages/cast-helper-git) | The git `master` branch to get all changes immediately. Potentially unstable |

```bash
yay -S cast-helper
sudo systemctl enable --now avahi-daemon
```

Optional dependencies for more formats, sound on the TV only, HEIC photos and AirPlay:
```bash
yay -S --needed --asdeps gst-plugins-good gst-libav libheif pipewire libpulse doubletake-bin
```

### Debian, Ubuntu, Fedora, openSUSE

Install the dependencies first. The first line of each command are required dependencies, the second contains the optional dependencies for full feature support and can be omitted.

Debian / Ubuntu:

```bash
sudo apt install avahi-daemon gstreamer1.0-plugins-base gstreamer1.0-plugins-bad gstreamer1.0-pipewire \
  gstreamer1.0-plugins-good gstreamer1.0-libav pipewire-bin libheif-examples
```

Fedora:

```bash
sudo dnf install avahi gstreamer1-plugins-base gstreamer1-plugins-bad-free pipewire-gstreamer \
  gstreamer1-plugins-good gstreamer1-plugin-libav pipewire-utils libheif-tools
```

openSUSE:

```bash
sudo zypper install avahi gstreamer-plugins-base gstreamer-plugins-bad gstreamer-plugin-pipewire \
  gstreamer-plugins-good gstreamer-plugins-libav pipewire-tools
```

Then start Avahi:

```bash
sudo systemctl enable --now avahi-daemon
```

Fedora's and openSUSE's Mesa can't encode H.264 or HEVC. For mirroring, install the full drivers from [RPM Fusion](https://rpmfusion.org/Howto/Multimedia) or [Packman](https://en.opensuse.org/Additional_package_repositories#Packman).

Now install cast-helper:

Apps look for it on `PATH` and in `~/.local/share/dms/cast/bin`, so put the release build there:

```bash
mkdir -p ~/.local/share/dms/cast/bin && cd ~/.local/share/dms/cast/bin \
  && curl -fLO https://github.com/Conava/cast-helper/releases/latest/download/cast-helper-$(uname -m) \
  && curl -fsL https://github.com/Conava/cast-helper/releases/latest/download/SHA256SUMS \
    | grep " cast-helper-$(uname -m)$" | sha256sum -c \
  && mv cast-helper-$(uname -m) cast-helper && chmod +x cast-helper
```

The Cast plugin can also download it for you from its settings in the background and will handle updates itself.

### Building from source / other distributions

Install Avahi and these GStreamer packages under your distribution's names: `gstreamer`, `gst-plugins-base`, `gst-plugins-bad`, `gst-plugin-pipewire` and `gst-plugin-va`.
Optionally install `gst-plugins-good`, `gst-libav` and `libheif`. 
Then build with Go 1.25 or newer:

```bash
git clone https://github.com/Conava/cast-helper.git && cd cast-helper && CGO_ENABLED=0 go build -o ~/.local/share/dms/cast/bin/cast-helper .
```

### AirPlay

cast-helper mirrors to AirPlay through [doubletake](https://github.com/omarroth/doubletake), which is not affiliated with this project. Install and trust it at your own discretion.
cast-helper is tested with doubletake 0.5.0 and looks for it on `PATH`, in `~/.local/share/dms/cast/bin` or in `~/go/bin`.

On Arch, you can install `doubletake-bin` from the AUR. Elsewhere, download the release (on aarch64 the file is `doubletake-arm64`):

```bash
curl -fL -o ~/.local/share/dms/cast/bin/doubletake \
  https://github.com/omarroth/doubletake/releases/download/v0.5.0/doubletake \
  && chmod +x ~/.local/share/dms/cast/bin/doubletake
```

doubletake runs `gst-launch-1.0` and `pactl` too. They're in `gstreamer1.0-tools` and `pulseaudio-utils` on Debian and Ubuntu, `pulseaudio-utils` on Fedora, `gstreamer-utils` and `pulseaudio-utils` on openSUSE, and `libpulse` on Arch.

The first AirPlay cast asks for a PIN that is shown on the TV. 
With AirPlay, sound on the TV is currently not supported, so it stays on your computer.

### Checking the setup

```bash
cast-helper version
cast-helper install-missing
```

`install-missing` asks to install whatever's missing through your package manager. 

## Network and files

cast-helper talks to the cast devices on your network. 
Beyond that it contacts any server that hosts a file you want to cast when a link is provided, and youtube.com to control the YouTube app on the cast device when a YouTube link is provided.
It never sends your files to the internet and opens no other connections.
Your files are served to the cast device and nothing else.

| Path | Description |
|---|---|
| `~/.cache/dms/cast` | Converted files while they are playing, thumbnails, the log |
| `~/.local/share/dms/cast` | AirPlay pairings, a downloaded cast-helper or doubletake |
| `$XDG_RUNTIME_DIR` | `cast-helper.sock` and `cast-helper.lock`, and `dms-cast/` during an AirPlay cast |

## Issues and questions

If the casting fails, [open an issue here](https://github.com/Conava/cast-helper/issues/new/choose):
For example: no image on the cast device, a device that won't connect, a file that won't play or a package reported missing that's installed. 
Buttons, menus and settings in Dux Shell are part of the [Cast plugin](https://github.com/Conava/dux-plugins/issues). If in doubt, please open an issue in the Plugin repo since that is monitored more closely.

Feature requests are welcome too. 

Often your support is required to debug issues and to successfully test bug fixes or new requested features, because this project depends a lot on the hardware and software of your computer, your network and your cast device. Please be patient and provide as much information as possible.

Want to use this in your app? [PROTOCOL.md](PROTOCOL.md) explains how to talk to cast-helper. 
To help with the code or contribute features yourself, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Plans and development goals

- Native AirPlay: A sender built into cast-helper to replace doubletake, support sound on the TV at the right volume, make files and music possible over AirPlay and to support AirPlay speakers like a HomePod.
- Fewer dependencies: Implement as many features as possible in cast-helper itself, so that it can be installed with fewer dependencies and run on more distributions.

## License

MIT, see [LICENSE](LICENSE).
