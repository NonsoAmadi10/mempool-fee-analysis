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

Once we have calculated the fee rate for each transaction in the mempool, we sort the transactions by fee rate, from highest to lowest. We then select the nth highest fee rate as our estimated fee rate.

In this implementation, I used the 10th highest fee rate as an example. This means that I selected the fee rate of the 10th transaction in the sorted list of transactions as the estimated fee rate. The idea is that by selecting a fee rate that is higher than most of the other transactions in the mempool, we can ensure that our transaction gets confirmed relatively quickly. However, selecting a fee rate that is too high can result in unnecessarily high transaction fees, so it's important to strike a balance between transaction confirmation speed and cost.