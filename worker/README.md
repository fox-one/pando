# Pando Workers

| Name      | Description                                                            |
| --------- | ---------------------------------------------------------------------- |
| cashier   | handle pending transfers ordered by transfer id                        |
| events    | tracking new transactions and notify tx's operator                     |
| messenger | send messages to mixin messenger user and clean sended                 |
| payee     | pull new multisig utxo and process actions according to the utxo order |
| pricesync | sync asset price with mixin wallet api                                 |
| spentsync | sync transfer status on mixin mainnet and notify the receiver if done  |
| txsender  | send raw transactions to mixin mainnet                                 |
