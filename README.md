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
go build ./cmd/neo-inu/main.go
```

## Usage

```bash
GUILD=<STRING> RMCMD=<BOOLEAN> TOKEN=<STRING> ./neo-inu
```

## Development

```bash
docker compose up -d
```

<p>Modify files in ./cmd/ ./pkg/ ./internal/ and the app will live reload</p>

## TODO

- More details Yu-gi-oh card response
- Play music
