# I'm Pando

Pando is a decentralized financial network built with the MTG technology, and its underlying financial algorithm is inspired by [MakerDao](https://makerdao.com) and Synthetix.

Pando takes multiple over-collateralized assets and grows multiple Pando assets, which(the grown Pando assets), e.g. pUSD, can also be used as a pledge.

## Mixin Network & MTG

### Mixin Network

Mixin Network is a public blockchain driven by TEE (Trusted Execution Environment) based on the DAG with aBFT. Unlike other projects which have great theories but hardly any actual implementations of blockchain transaction solution, Mixin Network provides a more secure, private, 0 fees, developer-friendly and user-friendly transaction solution with lightning speed.

### Mixin Trusted Group

[Mixin Trusted Group](https://developers.mixin.one/document/mainnet/mtg) is an alternative to smart contacts on Mixin Network.

Basically, MTG is a Multi-signature custodian consensus solution. Several teams will be selected and arranged as the “Trusted Group” in Pando, becoming the “Nodes”.
Concensus has to be reached among the nodes to perform certain administrative actions. As a result, stable services and asset safety are guaranteed.   

For example, let’s say there is a M/N multi-sig group where M represents the number of nodes, and the group manages some assets in the multi-sig address. When one of the nodes needs to transfer some assets out, it needs to collect at least N signatures from others to perform the action.

MTG is the framework. Pando is an application designed using the framework on Mixin Network. 

## An intro to Pando

### pUSD

pUSD is a stable coin launched on Dec 25, 2020. Since the launch, the backing ratio of pUSD has always been higher than 200%. 

### Pledge

Nodes have the ability to add any asset as a pledge by voting. Once nodes vote one asset to be a pledge, anybody can borrow pUSD by providing enough specified asset.

Nodes also have the ability to adjust various parameters of pledges by voting.

All assets supported by Mixin Network, including BTC, ETH, etc, are potential pledges for Pando.

### Vaults

All approved pledges can be deposited in the Pando multi-signature address to create a vault to generate pUSD for any Pando user.

As long as the price of the pledge is higher than the minimum requirement, the creators have the complete control of their vaults.

### Interact with Pando

Both users and node administrators must use multi-signature transactions to interact with Pando.

The methods and parameters of the transactions are all written in the memo which contains extra information of each transfer.

Currently, in order to protect user privacy, all information in the memo must be encrypted by the following algorithm:

Pando will generate a pair of ed25519 public and private keys at first, and publish the public key.

For each transaction, the user generates a pair of temporary ed25519 public and private keys, and generates a string of 32-bit bytes.

The first 16 bits of these bytes will be the key for AES encryption, and the last 16 bits will be the `iv` of AES encryption.

These bytes must encrypt the original memo and generate a result we call it encrypted bytes. The client should put encrypted bytes and the user’s public key ​​together and encode in base64 as the final memo for the transfer.

The nodes of Pando synchronize all transfers from the Mixin Network. When a node recognizes a valid transfer, it performs a reverse operation to restore the key and iv encrypted by AES, and then decrypts the parameters of the action.

The nodes need to ensure that they process the user interaction in the same order to ensure that the data stored in the database is completely consistent with other nodes; It should follow the same order when transferring money to ensure that all nodes choose when completing the transfer in the same utxo.

### Liquidate high-risk vaults

In order to ensure that there is always enough pledge in Pando to endorse the loaned pUSD, all vaults will require an excess mortgage such as 150%.

When the value of the mortgaged assets is insufficient due to market price fluctuations, the vault will enter a high-risk state and be open for liquidation to redeem the pUSD.The liquidation is carried out as auction:

- If the pUSD got by the auction is enough to pay off the debt in the vault and the liquidation penalty, the auction will sell just the enough amount of pledge to cover the debt and the remaining pledge will be returned to the original owner.
- If the pUSD got from the auction of the pledge is not enough to pay off the debt and the liquidation penalty, the loss will become Pando’s liability and be shared by all nodes.

### Price Oracle

Pando needs to synchronize the prices of the pledged assets to update the collateral rate of the vaults and liquidate the high-risk vaults.

The price data of Pando relies on an external decentralized price service. Pando will not use the price data directly.

Pando adds an one-hour delay to all price data. During this period, if someone attacks a price service, the nodes can vote to change the state of the pledge to be frozen urgently. Nodes can also vote for new prices.

### Summary and future works

Pando has achieved the goal that decentralized the consensus among trusted nodes, bringing the stable coin minting service to all users of the Mixin Network.

Pando also reserves the ability to extend the lending to non-stable assets. For example, it can issue a 1:1 token pTesla against Tesla stock on the Mixin network, and then plege the Bitcoin on Pando to generate pTesla, which will build connections between the assets in the Mixin Network and the external assets of outside world.
