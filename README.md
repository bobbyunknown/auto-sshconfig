# Auto SSH Config

Setup SSH key dan config otomatis dengan Go.

## Build

```bash
go build -o auto-sshconfig main.go
```

## Commands

```bash
auto-sshconfig              # Setup SSH baru (interactive)
auto-sshconfig -h           # Help
auto-sshconfig -r vpstest   # Remove SSH config
auto-sshconfig -c vpstest root@192.168.1.100  # Copy key ke server
auto-sshconfig -s           # Show SSH configs
```

## Setup SSH Baru

```bash
auto-sshconfig
```

Input yang diminta:
- Nama host (contoh: vpstest)
- IP address
- Username (default: root)

Yang dilakukan:
1. Hapus known_hosts yang bermasalah
2. Generate SSH key baru
3. Copy public key ke server
4. Update SSH config
5. Test koneksi

Setelah selesai bisa SSH dengan: `ssh namahost`

## Show SSH Configs

```bash
auto-sshconfig -s
```

Menampilkan daftar SSH configs yang ada.

## Copy Key ke Server Lain

```bash
auto-sshconfig -c vpstest root@192.168.1.100
```

Copy key yang sudah ada (vpstest) ke server lain.

## Remove SSH Config

```bash
auto-sshconfig -r vpstest
```

Menghapus key files dan entry dari SSH config.

## File yang dibuat

- `~/.ssh/id_rsa_namahost` - Private key
- `~/.ssh/id_rsa_namahost.pub` - Public key
- `~/.ssh/config` - SSH config (diupdate)
