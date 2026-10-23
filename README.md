# GoDeb Tunnel (Arch / Linux)

Native GoDeb Tunnel client for Arch Linux (and other Linux distros with `wg-quick`).

Matches the Android / Windows GoDeb clients for:

- **Trusted endpoint auto-disconnect** (per tunnel via conf comment)
- **Excluded AllowedIPs** (CIDR subtract helper)

Uses the system `wireguard-tools` (`wg`, `wg-quick`) — no extra tunnel stack.

## Install (Arch)

From this repo (local package):

```bash
cd packaging/arch
# point source at a local tree or release tarball, then:
makepkg -si
```

Or build the binary directly:

```bash
go build -o godeb-wg ./cmd/godeb-wg
sudo install -m755 godeb-wg /usr/bin/godeb-wg
sudo install -m644 systemd/godeb-wg-monitor.service /usr/lib/systemd/user/
```

Enable the user monitor:

```bash
systemctl --user enable --now godeb-wg-monitor.service
```

## Usage

```bash
# Import a GoDeb / WireGuard client config
godeb-wg import home ~/Downloads/home.conf

godeb-wg list
sudo godeb-wg up home
sudo godeb-wg down home

# Compute AllowedIPs = 0.0.0.0/0 minus LAN
godeb-wg exclude '0.0.0.0/0, ::/0' '172.16.1.0/24'

# Foreground trusted monitor (or use the systemd user unit)
godeb-wg monitor
```

### Conf comments (same as Android / Windows)

```
# GoDeb-Trusted-Endpoint: http://172.16.1.20:8080/api/trusted
# GoDeb-AllowedIPs-Base: 0.0.0.0/0, ::/0
# GoDeb-ExcludedIPs: 172.16.1.0/24
```

Configs are stored in `~/.config/godeb-wg/tunnels/`.

## Related

- [godeb-tunnel-android](https://github.com/ZeroDos1988/godeb-tunnel-android)
- Windows client: not public yet
- [GoDeb](https://github.com/ZeroDos1988/GoDeb)

## License

GPL-3.0-only, see [LICENSE](LICENSE). Copyright (C) 2026 ZeroDos1988.

## Trademark

GoDeb Tunnel is an independent project and is not affiliated with or endorsed by the WireGuard project. It uses the system's `wg-quick` / `wireguard-tools`. "WireGuard" is a registered trademark of Jason A. Donenfeld.
