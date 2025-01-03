# Neo Inu

## Prerequisite

- bash
- docker

## Installation

```bash
git clone https://github.com/TranThangBin/neo-inu.git
```

<p>Build the binary if you need nothing else</p>

```bash
make build_bin
```

## Usage

```
Usage of ./bin/neo-inu:
  -guild string
        Test guild ID default: "" (mean global)
  -rmcmd
        Remove all command after shutdown default: true (default true)
  -token string
        Your discord bot token look for TOKEN variable if not provide (default $TOKEN)
```

## Development

1. Give main.sh execute permission

```bash
chmod +x ./script/main.sh
```

2. Build the container

```bash
make
```

3. Run the container

```bash
make run
```
