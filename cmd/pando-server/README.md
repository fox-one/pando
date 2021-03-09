# Pando Server

## Build

```shell
# build binary
make build-server

# build docker image
make pando/server
```

## API

### Login

```http request
POST /api/login
```

**Params:**

```json5
{
  "code": "mixin oauth code"
}
```

**Response:**

```json5
{
   "ts": 1614858782118,
   "data": {
      "avatar": "https://mixin-images.zeromesh.net/Fh-jsEMf7KYyjyhtUoEpVjMUhIT2cZPIGqfDxtHNxNoG-2ruJYFtAJoeqexkKBn8AlptnUSZW-eKTWF6KRbo9K7J=s256",
      "id": "8017d200-7870-4b82-b53f-74bae1d2dad7",
      "name": "yiplee@fox",
      "scope": "PROFILE:READ ASSETS:READ SNAPSHOTS:READ",
      "token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiI5YzMzNjhkNy02NjAzLTQ0ODQtYWQ5ZC1jMmUyNWFhYWZkMjIiLCJleHAiOjE2NDYzOTQ3ODEsImlhdCI6MTYxNDg1ODc4MSwiaXNzIjoiNjcwZTFmYWEtMjk3NS00OGQ5LWE4MWYtY2QwOTA1YWU4NDdlIiwic2NwIjoiUFJPRklMRTpSRUFEIEFTU0VUUzpSRUFEIFNOQVBTSE9UUzpSRUFEIn0.R15j1h9zPxL-isgIxqaBARHB5GX3caEwTMllidW6DyT8bdtCFK211_RMfGQ_fp8eFHTGCiTRTBeEhObhdQarN9rTv0qinL1Piv4cugWaEvygJofnEkE8q5Gu_5vAPLjbs7j1ghqVAfz2VHKLOE9GDnyqnF0rlGqI3hCSFzyq2tc"
   }
}
```

### List Assets

> List All Assets

```http request
GET /api/assets
```

**Response:**

```json5
{
  "ts": 1615273615698,
  "data": {
    "assets": [
      {
        "id": "43d61dcd-e413-450d-80b8-101d5e903357",
        "name": "Ether",
        "symbol": "ETH",
        "logo": "https://mixin-images.zeromesh.net/zVDjOxNTQvVsA8h2B4ZVxuHoCF3DJszufYKWpd9duXUSbSapoZadC7_13cnWBqg0EmwmRcKGbJaUpA8wFfpgZA=s128",
        "chain_id": "43d61dcd-e413-450d-80b8-101d5e903357",
        "chain": {
          "id": "43d61dcd-e413-450d-80b8-101d5e903357",
          "name": "Ether",
          "symbol": "ETH",
          "logo": "https://mixin-images.zeromesh.net/zVDjOxNTQvVsA8h2B4ZVxuHoCF3DJszufYKWpd9duXUSbSapoZadC7_13cnWBqg0EmwmRcKGbJaUpA8wFfpgZA=s128",
          "chain_id": "43d61dcd-e413-450d-80b8-101d5e903357",
          "chain": null,
          "price": "1827.79"
        },
        "price": "1827.79"
      }
    ]
  }
}
```

### Find Asset

```http request
GET /api/assets/asset_id
```

**Response:**

```json5
{
  "ts": 1615273696911,
  "data": {
    "id": "43d61dcd-e413-450d-80b8-101d5e903357",
    "name": "Ether",
    "symbol": "ETH",
    "logo": "https://mixin-images.zeromesh.net/zVDjOxNTQvVsA8h2B4ZVxuHoCF3DJszufYKWpd9duXUSbSapoZadC7_13cnWBqg0EmwmRcKGbJaUpA8wFfpgZA=s128",
    "chain_id": "43d61dcd-e413-450d-80b8-101d5e903357",
    "chain": {
      "id": "43d61dcd-e413-450d-80b8-101d5e903357",
      "name": "Ether",
      "symbol": "ETH",
      "logo": "https://mixin-images.zeromesh.net/zVDjOxNTQvVsA8h2B4ZVxuHoCF3DJszufYKWpd9duXUSbSapoZadC7_13cnWBqg0EmwmRcKGbJaUpA8wFfpgZA=s128",
      "chain_id": "43d61dcd-e413-450d-80b8-101d5e903357",
      "chain": null,
      "price": "1827.79"
    },
    "price": "1827.79"
  }
}
```

### List Oracles

```http request
GET /api/oracles
```

#### [Oracle]

| key      | description                       |
| -------- | --------------------------------- |
| asset_id | mixin asset id                    |
| hop      | price delay (seconds)             |
| current  | current price                     |
| next     | next price, peek at peek_at + hop |
| peek_at  | last update of price              |

**Response:**

```json5
{
    "data": {
        "oracles": [
            {
                "asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
                "current": "1",
                "hop": 3600,
                "next": "1",
                "peek_at": "2021-03-09T06:00:00Z"
            },
            {
                "asset_id": "6770a1e5-6086-44d5-b60f-545f9d9e8ffd",
                "current": "5",
                "hop": 3600,
                "next": "5",
                "peek_at": "2021-03-09T06:00:00Z"
            }
        ]
    },
    "ts": 1615272883769
}
```

### Find Oracles

```http request
GET /api/oracles/asset_id
```

**Response:**

```json5
{
    "data": {
        "asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
        "current": "1",
        "hop": 3600,
        "next": "1",
        "peek_at": "2021-03-09T06:00:00Z"
    },
    "ts": 1615273811937
}
```

### List Collaterals

> List All Collaterals

```http request
GET /api/cats
```

#### [Collateral]

| key  | description                      |
| ---- | -------------------------------- |
| name | name                             |
| gem  | collateral asset id              |
| dai  | debt asset id                    |
| ink  | gem total deposited amount       |
| rate | total normalised debt            |
| rho  | last update time of rate         |
| debt | total debt                       |
| line | max debt in total                |
| dust | minimum debt per vault           |
| mat  | liquidation rate                 |
| chop | liquidation fee                  |
| dunk | liquidation limit per flip       |
| beg  | [flip] minimum bid increase      |
| ttl  | [flip] bid duration in seconds   |
| tau  | [flip] auction length in seconds |
| live | collateral state                 |


**Response:**

```json5
{
   "ts": 1614857763109,
   "data": {
      "collaterals": [
         {
            "id": "0439b3e4-61a8-3ff4-9d3c-fe223ff55244",
            "created_at": "2021-03-03T08:29:35Z",
            "name": "XIN/CNB",
            "gem": "c94ac88f-4671-3976-b60a-09064f1811e8", // Collateral Asset ID
            "dai": "965e5c6e-434c-3fa9-b780-c50f43cd955c", // Debt Asset ID
            "ink": "1",                                    // Total Deposited
            "art": "120",                                  // Total Normalised Debt
            "rate": "1.0000731854582843",                  // Accumulated Rates
            "rho": "2021-03-04T06:10:53Z",                 // Rate Update Date
            "debt": "120",                                 // Total Debt
            "line": "10001",                               // Max Debt
            "dust": "100",                                 // Minimum Debt Per Vault
            "price": "150",                                // Current Price
            "duty": "1.03",                                // Stability Fee
            "mat": "1.1",                                  // Liquidation Rate
            "chop": "1.13",                                // Liquidation Fee
            "dunk": "5000",                                // Max Liquidation Debt
            "beg": "0.03",                                 // minimum bid increase
            "ttl": 900,                                    // bid duration in seconds
            "tau": 10800,                                  // auction length in seconds
            "live": true                                   // Collateral State
         }
      ]
   }
}
```

### List My Vaults

```http request
GET /api/vats
```

**Response:**

```json5
{
   "ts": 1614858173421,
   "data": {
      "vaults": [
         {
            "id": "e40060ae-fb63-3b6b-8c17-72550ffa5a5d",
            "created_at": "2021-03-03T08:56:34Z",
            "collateral_id": "0439b3e4-61a8-3ff4-9d3c-fe223ff55244",
            "ink": "1",                         // Total Deposited
            "art": "120"                        // Total Normalised Debt, debt = art * rate
         }
      ]
   }
}
```

### Get Tx

```http request
GET /transactions/{follow_id}
```

**Response:**

```json5
{
  "data": {
    "action": "CatEdit",
    "amount": "1",
    "asset_id": "965e5c6e-434c-3fa9-b780-c50f43cd955c",
    "created_at": "2021-03-09T06:29:21Z",
    "id": "3f973c40-1e6c-5572-b5d5-e55a7117a512",
    "msg": "",
    "parameters": "[\"f7c1ba17-67b9-34f3-adec-18bf4d63931b\",\"price\",\"5\",\"live\",\"1\"]",
    "status": "OK" // OK/Abort
  },
  "ts": 1615272801003
}
```

### Actions

```http request
POST /api/actions
```

**Params:**

```json5
{
  "asset_id": "", // payment asset id
  "amount": "123", // payment amount
  "actions": ["uuid","xxx"]
}
```

**Response:**

```json5
{
    "data": {
      "memo": "xxx",
      "code": "xxxx",
      "core_url": "https://mixin.one/codes/xxxx"
    },
    "ts": 1614858409361
}
```
