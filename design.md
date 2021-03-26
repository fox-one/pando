# Pando 设计文档

## 与 Pando 交互

Pando 里面的所有角色都是通过给 Pando 多签收款地址转账来完成，把要操作的方法和参数写在 memo 里面；节点 Syncer（worker/syncer）会把这些外部转账即 Mixin Multisig Output 同步过来，然后交给 Payee（worker/payee）去按顺序处理。

### Mixin Multisig Output

**Output:**

| field     | description      |
| --------- | ---------------- |
| Sender    | user mixin id    |
| CreatedAt | payment time     |
| AssetID   | payment asset id |
| Amount    | payment amount   |
| Mmeo      | extra message    |

**Output Memo:**

如果 memo 是 aes 加密的，那么前 32 位是用于计算 aes key/iv 的 ed25519 公钥，后面为加密后的附加信息。

memo 解密之后得到 **TransactionAction** (core/action.go#TransactionAction)

### TransactionAction

| field    | description                                | type  |
| -------- | ------------------------------------------ | ----- |
| FollowID | user defined trace id for this transaction | uuid  |
| Body     | action type & parameters                   | bytes |

## Actions

all actions supported by Pando with groups cat,flip,oracle,proposal,sys and vat. see [core/action](core/action.go) for details.

### Sys - system operations

#### #1 Withdraw

withdraw any assets from the multisig wallet, proposal required.

[pkg/maker/sys/withdraw](pkg/maker/sys/withdraw.go)

**Parameters:**

| name     | type | description         |
| -------- | ---- | ------------------- |
| asset    | uuid | withdraw asset id   |
| amount   | uuid | withdraw amount     |
| opponent | uuid | receiver's mixin id |

### Proposal - governance system

#### #11 Make 

create a new proposal

[pkg/maker/proposal/make](pkg/maker/proposal/make.go)

**Parameters:**

| name | type  | description                                         |
| ---- | ----- | --------------------------------------------------- |
| data | bytes | action type & parameters will be executed if passed |

#### #12 Shout

request node administrator to vote for this proposal

[pkg/maker/proposal/shout](pkg/maker/proposal/shout.go)

**Parameters:**

| name | type | description       |
| ---- | ---- | ----------------- |
| id   | uuid | proposal trace id |

#### #13 Vote

vote for a proposal, node only. If enough votes collected, the attached action will be executed on all nodes automatically.

[pkg/maker/proposal/vote](pkg/maker/proposal/vote.go)

**Parameters:**

| name | type | description       |
| ---- | ---- | ----------------- |
| id   | uuid | proposal trace id |

### Cat - manager collaterals

#### #21 Create

create a new collateral type, proposal required.

[pkg/maker/cat/create](pkg/maker/cat/create.go)

**Parameters:**

| name | type   | description          |
| ---- | ------ | -------------------- |
| gem  | uuid   | collateral asset id  |
| dai  | uuid   | debt asset id        |
| name | string | collateral type name |

#### #22 Supply

supply dai token to increase the total debt ceiling for this collateral type.
Payment asset id must be equal with the debt asset id.

[pkg/maker/cat/supply](pkg/maker/cat/supply.go)

**Parameters:**

| name | type | description         |
| ---- | ---- | ------------------- |
| id   | uuid | collateral trace id |

#### #23 Edit

modify collateral's one or more attributes, proposal required.

[pkg/maker/cat/edit](pkg/maker/cat/edit.go)

**Parameters:**

| name  | type   | description         |
| ----- | ------ | ------------------- |
| id    | uuid   | collateral trace id |
| key   | string | attribute name      |
| value | string | attributes value    |

#### #24 Fold

modify the debt multiplier(rate), creating / destroying corresponding debt.

[pkg/maker/cat/fold](pkg/maker/cat/fold.go)

**Parameters:**

| name | type | description         |
| ---- | ---- | ------------------- |
| id   | uuid | collateral trace id |

### Vat - manager vaults

#### #31 Open 

open a new vault with the special collateral type

[pkg/maker/vat/open](pkg/maker/vat/open.go)

**Parameters:**

| name | type    | description         |
| ---- | ------- | ------------------- |
| id   | uuid    | collateral trace id |
| debt | decimal | initial debt        |

#### #32 Deposit

transfer collateral into a Vault.

[pkg/maker/vat/deposit](pkg/maker/vat/deposit.go)

**Parameters:**

| name | type | description    |
| ---- | ---- | -------------- |
| id   | uuid | vault trace id |

#### #33 Withdraw

withdraw collateral from a Vault, vault owner only.

[pkg/maker/vat/withdraw](pkg/maker/vat/withdraw.go)

**Parameters:**

| name | type    | description          |
| ---- | ------- | -------------------- |
| id   | uuid    | vault trace id       |
| dink | decimal | change in collateral |

#### #34 Payback

 decrease Vault debt.

 [pkg/maker/vat/payback](pkg/maker/vat/payback.go)

**Parameters:**

| name | type | description    |
| ---- | ---- | -------------- |
| id   | uuid | vault trace id |

#### #35 Generate

increase Vault debt, vault owner only.

 [pkg/maker/vat/generate](pkg/maker/vat/generate.go)

**Parameters:**

| name | type    | description    |
| ---- | ------- | -------------- |
| id   | uuid    | vault trace id |
| debt | decimal | change in debt |

### Flip - manager auctions

#### #41 Kick

put collateral up for auction from an unsafe vault.

[pkg/maker/flip/kick](pkg/maker/flip/kick.go)

**Parameters:**

| name | type | description    |
| ---- | ---- | -------------- |
| id   | uuid | vault trace id |

#### #42 Bid

bid for the collateral.

> Starting in the tend-phase, bidders compete for a fixed lot amount of Gem with increasing bid amounts of Dai. Once tab amount of Dai has been raised, the auction moves to the dent-phase. The point of the tend phase is to raise Dai to cover the system's debt.
> During the dent-phase bidders compete for decreasing lot amounts of Gem for the fixed tab amount of Dai. Forfeited Gem is returned to the liquidated Urn for the owner to retrieve. The point of the dent phase is to return as much collateral to the Vault holder as the market will allow.

[pkg/maker/flip/bid](pkg/maker/flip/bid.go)

**Parameters:**

| name | type    | description       |
| ---- | ------- | ----------------- |
| id   | uuid    | flip trace id     |
| lot  | decimal | collateral amount |

#### #43 Deal

claim a winning bid / settles a completed auction

[pkg/maker/flip/deal](pkg/maker/flip/deal.go)

**Parameters:**

| name | type | description   |
| ---- | ---- | ------------- |
| id   | uuid | flip trace id |

### Oracle - manager price oracle

#### #51 Create

register a new oracle for asset, proposal required.

[pkg/maker/oracle/create](pkg/maker/oracle/create.go)

**Parameters:**

| name      | type      | description                              |
| --------- | --------- | ---------------------------------------- |
| id        | uuid      | asset id                                 |
| price     | decimal   | initial price                            |
| hop       | int64     | time delay in seconds between poke calls |
| threshold | int64     | number of governors required when poke   |
| ts        | timestamp | request timestamp                        |

#### #52 Edit

modify a oracle's next price, hop & threshold, proposal required.

[pkg/maker/oracle/edit](pkg/maker/oracle/edit.go)

**Parameters:**

| name  | type   | description      |
| ----- | ------ | ---------------- |
| id    | uuid   | asset id         |
| key   | string | attribute name   |
| value | string | attributes value |

#### #53 Poke

updates the current feed value and queue up the next one.

[pkg/maker/oracle/poke](pkg/maker/oracle/poke.go)

**Parameters:**

| name  | type      | description       |
| ----- | --------- | ----------------- |
| id    | uuid      | asset id          |
| price | decimal   | new next price    |
| ts    | timestamp | request timestamp |

#### #54 Rely

add a new price feed to the whitelist, proposal required

[pkg/maker/oracle/rely](pkg/maker/oracle/rely.go)

| name | type  | description   |
| ---- | ----- | ------------- |
| id   | uuid  | feed mixin id |
| key  | bytes | public key    |

#### #55 Deny

remove a price feed from the whitelist, proposal required

[pkg/maker/oracle/deny](pkg/maker/oracle/deny.go)

| name | type | description   |
| ---- | ---- | ------------- |
| id   | uuid | feed mixin id |
