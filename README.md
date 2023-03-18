# mempool-fee-analysis
An API and cli that helps users decide what amount of fee to pay for bitcoin transaction based on mempool analysis


## Alogrithm at a High Level 
To estimate the fee, we look at the current mempool, which is a list of all unconfirmed transactions that have been broadcast to the Bitcoin network. We then calculate the fee rate for each transaction in the mempool. The fee rate is the amount of Bitcoin (in satoshis) that the transaction pays per byte of data it includes.

To calculate the fee rate for each transaction, we first need to calculate the total input and output values of the transaction. The input values are the amounts of Bitcoin that the transaction is spending, and the output values are the amounts of Bitcoin that the transaction is sending to its recipients. The difference between the input and output values is the transaction fee.

We then divide the transaction fee by the size of the transaction in bytes to get the fee rate. The fee rate is expressed in satoshis per byte. For example, if a transaction has a fee of 100 satoshis and a size of 200 bytes, the fee rate is 0.5 satoshis per byte.

Once we have calculated the fee rate for each transaction in the mempool, we sort the transactions by fee rate, from highest to lowest. We then select the nth highest fee rate as our estimated fee rate.

In my current implementation, I used the 10th highest fee rate as an example. This means that I selected the fee rate of the 10th transaction in the sorted list of transactions as the estimated fee rate. The idea is that by selecting a fee rate that is higher than most of the other transactions in the mempool, we can ensure that our transaction gets confirmed relatively quickly. However, selecting a fee rate that is too high can result in unnecessarily high transaction fees, so it's important to strike a balance between transaction confirmation speed and cost.

