# OSU-API-Server Self-Hosted Guide

## With Docker

This part explains how to install oas in your server with Docker.

### Docker Compose

It's easy to set up and maintain oas with docker-compose. Below is
an easy setup `docker-compose.yml` example.

```yaml
version: '3'

services:
  server:
    image: avimitin/osu-api-server:v0.1
    restart: always
    environment:
      - OSU_API_KEY=YOURKEY
      - OSU_DB_TYPE=redis
      - OSU_DB_HOST=database
    volumes:
      - ./data:/data
    ports:
      - "11451:80"
    depends_on:
      - database
  
  database:
    image: redis:latest
    restart: always
```

Run the command `docker-compose up -d` to start the server

## Build

You can also choose to build the binary yourself.

### STEP 1

Install the [Go dependency](https://golang.org/dl/)
and [Redis](https://redis.io/download)

### STEP 2

Run command to build the server binary.

```bash
go build -o /usr/local/bin/osu-api-server cmd/cmd.go
```

### STEP 3

Put your config under `$HOME/.config/osuapi/config.json`

### STEP 4

Write a systemd service profile to guard the program.

---
osu-api-server.service

```unit file (systemd)
[Unit]
Description=Osu API SERVER
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/osu-api-server
Restart=always

[Install]
WantedBy=multi-user.target
```

---

Start the service with command: `systemctl start osu-api-server`