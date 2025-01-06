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

```
Usage of ./neo-inu
  -guild string
        Test guild ID. $GUILD is prioritized
  -rmcmd
        Remove all command after shutdown. $RMCMD is prioritized (default true)
  -token string
        Your discord bot token. $TOKEN is prioritized.
```

## Development

### Start the container

<p>Remember to create your .env file</p>

```bash
touch .env
```

```bash
TOKEN="<CHECK -token ON USAGE>"
RMCMD="<CHECK -rmcmd ON USAGE>"
GUILD="<CHECK -guild ON USAGE>"
```

```bash
docker compose up -d
```

<p>Modify files in ./cmd/ ./pkg/ ./internal/ and the app will live reload</p>

## TODO

- More details Yu-gi-oh card response
- Play music
