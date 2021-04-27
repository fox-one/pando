# Pando/worker

## config.yaml

```yaml
db:
  dialect: mysql
  host: db
  port: 3306
  user: root
  password: pando
  database: pando_1

# node_1
dapp:
  pin: ''
  client_id:
  session_id:
  pin_token:
  private_key:

group:
  # shared ed25519 private key among all members
  private_key: wp+18tCGzWix0Qp3BHJBni/tREgLdF77nRFEz1ZrNSAM/15Fhnywrv9tgX3xgbaL/cgx1fI9VCajSv/1bI0Ddg==
  admins:
    - 8017d200-7870-4b82-b53f-74bae1d2dad7
  members:
    - 670e1faa-2975-48d9-a81f-cd0905ae847e
    - 229fc7ac-9d09-4a6a-af5a-78f7439dce76
    - 8017d200-7870-4b82-b53f-74bae1d2dad7
    - 170e40f0-627f-4af2-acf5-0f25c009e523
    - dfa655ef-55db-4e18-bdd7-29a7c576a223
  threshold: 3
```

## Build

```bash
# install task
go install github.com/go-task/task/v3/cmd/task@latest

# build docker image
task pando/worker

# build binary
task build-worker
```

## Deploy

### Run Binary

```bash
./pando-worker --config config.yaml
```

### Docker Compose

> **pando/worker** https://github.com/fox-one/pando/packages/752854

```yaml
version: "3.9"

services:
  worker:
    # cat github_token | docker login -u username --password-stdin docker.pkg.github.com
    # generate github token with package:read scope on https://github.com/settings/tokens
    image: docker.pkg.github.com/fox-one/pando/worker:1.4.8
    restart: always
    volumes:
      - ./config.yaml:/app/config.yaml
    ports:
      - "7777:7777"

  db:
    image: mysql:5.7
    restart: always
    volumes:
      - data:/var/lib/mysql
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=pando
      - MYSQL_USER=pando
      - MYSQL_PASSWORD=pando

volumes:
  data:
```
