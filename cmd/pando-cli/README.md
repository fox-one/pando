# Pando CLI

## install

```bash
# install task
go install github.com/go-task/task/v3/cmd/task@latest

# install cli, will install at $GOBIN/pd
task install-cli
```

## Config pd

Setup api

```bash
# for example, use pando production api
pd use https://leaf-api.pando.im
```

## Auth

Open the oauth page, scan with Mixin Messenger to get an auth code

```bash
pd auth login
```

Login with the auth code

```bash
pd auth login the_auth_code
```

Run ```pd help``` to get more
