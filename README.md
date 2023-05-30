# Mempool Fee Analyzer

This is a project written in Golang that estimates the fee for a Bitcoin transaction by analyzing the current mempool.

## Prerequisites

- Golang v1.16 or higher
- Bitcoin node with RPC enabled
- Access to the Bitcoin RPC credentials

## Installation

1. Clone the repository to your local machine.
2. Install the necessary dependencies by running `go mod download`.

## Usage

1. Start your Bitcoin node with RPC enabled.
2. In the terminal, navigate to the project directory.
3. Run `go run main.go` to start the fee estimator.
4. Make an API request to `http://localhost:4000/best-fee`
5. The estimated fee rate will be displayed in satoshis per byte.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Reference

To estimate the fee, we look at the current mempool, which is a list of all unconfirmed transactions that have been broadcast to the Bitcoin network. We then calculate the fee rate for each transaction in the mempool. The fee rate is the amount of Bitcoin (in satoshis) that the transaction pays per byte of data it includes.

To calculate the fee rate for each transaction, we first need to calculate the total input and output values of the transaction. The input values are the amounts of Bitcoin that the transaction is spending, and the output values are the amounts of Bitcoin that the transaction is sending to its recipients. The difference between the input and output values is the transaction fee.

We then divide the transaction fee by the size of the transaction in bytes to get the fee rate. The fee rate is expressed in satoshis per byte. For example, if a transaction has a fee of 100 satoshis and a size of 200 bytes, the fee rate is 0.5 satoshis per byte.

Once we have calculated the fee rate for each transaction in the mempool, we sort the transactions by fee rate, from highest to lowest. We then select the highest fee rate in the nth-percentile as our estimated fee rate.

In this implementation, I used the highest fee rate in the 80th-percentile as an example. The idea is that by selecting a fee rate that is higher than most of the other transactions in the mempool, we can ensure that our transaction gets confirmed relatively quickly. However, selecting a fee rate that is too high can result in unnecessarily high transaction fees, so it's important to strike a balance between transaction confirmation speed and cost.

### Half-hour Fee Estimation
The half-hour algorithm is a way to estimate the transaction fee based on the recent transaction activity in the Bitcoin network. The idea is to calculate the median fee rate of the transactions included in the previous half hour, and use that as an estimate for the current transaction fee.

The reason why this algorithm is useful is that the transaction fees in the Bitcoin network can be highly variable, depending on the amount of traffic and the urgency of the transaction. By using the recent transaction activity as a guide, the half-hour algorithm can provide a more accurate estimate of the appropriate transaction fee.

Here's how the getHalfhourFee function implements the half-hour algorithm:

1. First, it gets the current block height from the Bitcoin client using the GetBlockCount function.

2. Then, it calculates the block height of the block that was mined half an hour ago by subtracting 6 blocks (since blocks are mined roughly every 10 minutes) from the current block height.

3. Next, it retrieves the block hash of the block that was mined half an hour ago using the GetBlockHash function.

4. It then retrieves the block data for the block that was mined half an hour ago using the GetBlockVerbose function.

5. From the block data, it extracts the list of transactions that were included in the block.

6. For each transaction in the list, it calculates the fee per byte by dividing the transaction fee by the transaction size in bytes.

7. It then adds the fee per byte for each transaction to a slice of fee rates.

8. Once it has collected all the fee rates, it sorts the slice in ascending order.

9. It calculates the median fee rate by taking the middle value of the sorted slice.

10. It returns the median fee rate as the estimated fee per byte for the current transaction.

The code we just cooked implements the half-hour algorithm in Go using the btcd package, which provides an interface to a Bitcoin client. It uses a combination of Bitcoin client functions and built-in Go functions to retrieve the necessary data and perform the necessary calculations. Overall, the getHalfhourFee function provides a convenient way to estimate the appropriate transaction fee based on recent transaction activity in the Bitcoin network.