package mempoolfee

import (
	"log"
	"math"
	"sort"

	"github.com/NonsoAmadi10/mempool-fee/utils"
	"github.com/btcsuite/btcd/btcutil"
)

func GetBestFee() float64 {

	client := utils.Bitcoind()
	defer client.Shutdown()

	// Get the current mempool
	mempool, err := client.GetRawMempool()
	if err != nil {
		log.Fatal(err)
	}

	// Calculate fee rates for each transaction in the mempool
	type txInfo struct {
		feeRate float64
		size    int
	}
	txInfos := make([]txInfo, 0, len(mempool))
	for _, txid := range mempool {
		tx, err := client.GetRawTransaction(txid)
		if err != nil {
			log.Fatal(err)
		}

		totalIn := 0.0
		for _, in := range tx.MsgTx().TxIn {
			prevTxHash := &in.PreviousOutPoint.Hash
			prevTx, err := client.GetRawTransaction(prevTxHash)
			if err != nil {
				log.Fatal(err)
			}
			prevOut := prevTx.MsgTx().TxOut[in.PreviousOutPoint.Index]
			totalIn += btcutil.Amount(prevOut.Value).ToBTC()
		}
		totalOut := 0.0
		for _, out := range tx.MsgTx().TxOut {
			totalOut += btcutil.Amount(out.Value).ToBTC()
		}
		fee := totalIn - totalOut
		feeRate := fee / float64(tx.MsgTx().SerializeSize())
		txInfos = append(txInfos, txInfo{feeRate, tx.MsgTx().SerializeSize()})
	}

	// Sort transactions by fee rate
	sort.Slice(txInfos, func(i, j int) bool {
		return txInfos[i].feeRate > txInfos[j].feeRate
	})

	// Estimate the fee rate as the n-th highest fee rate
	n := 12 // Use the 10th highest fee rate as an example
	if n > len(txInfos) {
		n = len(txInfos) + 1
	}
	feeRateEstimate := txInfos[n-2].feeRate

	// Print the estimated fee rate
	log.Printf("Estimated fee rate: %f sat/vB\n", feeRateEstimate*1e8)

	bestFee := feeRateEstimate * 1e8

	return bestFee
}

func GetPriorityFees() (highFeeRate, normalFeeRate, lowFeeRate float64, err error) {

	client := utils.Bitcoind()

	defer client.Shutdown()

	// Get the current mempool
	mempool, err := client.GetRawMempool()
	if err != nil {
		return 0, 0, 0, err
	}

	// Calculate fee rates for each transaction in the mempool
	type txInfo struct {
		feeRate float64
		size    int
	}
	txInfos := make([]txInfo, 0, len(mempool))
	for _, txid := range mempool {
		tx, err := client.GetRawTransaction(txid)
		if err != nil {
			return 0, 0, 0, err
		}
		totalIn := 0.0
		for _, in := range tx.MsgTx().TxIn {
			prevTxHash := &in.PreviousOutPoint.Hash
			prevTx, err := client.GetRawTransaction(prevTxHash)
			if err != nil {
				return 0, 0, 0, err
			}
			prevOut := prevTx.MsgTx().TxOut[in.PreviousOutPoint.Index]
			totalIn += btcutil.Amount(prevOut.Value).ToBTC()
		}
		totalOut := 0.0
		for _, out := range tx.MsgTx().TxOut {
			totalOut += btcutil.Amount(out.Value).ToBTC()
		}
		fee := totalIn - totalOut
		feeRate := fee / float64(tx.MsgTx().SerializeSize())
		txInfos = append(txInfos, txInfo{feeRate, tx.MsgTx().SerializeSize()})
	}

	// Sort transactions by fee rate
	sort.Slice(txInfos, func(i, j int) bool {
		return txInfos[i].feeRate > txInfos[j].feeRate
	})

	// Estimate the fee rate as the n-th highest fee rate
	n := 12 // Use the 10th highest fee rate as an example
	if n > len(txInfos) {
		n = len(txInfos)
	}
	highPriority := txInfos[0].feeRate * 1e8
	normalPriority := txInfos[n-2].feeRate * 1e8
	lowPriority := txInfos[len(txInfos)-1].feeRate * 1e8

	return highPriority, normalPriority, lowPriority, nil
}

func GetImprovedBestFee() float64 {
	client := utils.Bitcoind()
	defer client.Shutdown()

	// Get the current mempool
	mempool, err := client.GetRawMempool()
	if err != nil {
		log.Fatal(err)
	}

	// Calculate fee rates for each transaction in the mempool
	type txInfo struct {
		feeRate float64
		size    int
	}
	txInfos := make([]txInfo, 0, len(mempool))
	for _, txid := range mempool {
		tx, err := client.GetRawTransaction(txid)
		if err != nil {
			log.Fatal(err)
		}

		totalIn := 0.0
		for _, in := range tx.MsgTx().TxIn {
			prevTxHash := &in.PreviousOutPoint.Hash
			prevTx, err := client.GetRawTransaction(prevTxHash)
			if err != nil {
				log.Fatal(err)
			}
			prevOut := prevTx.MsgTx().TxOut[in.PreviousOutPoint.Index]
			totalIn += btcutil.Amount(prevOut.Value).ToBTC()
		}
		totalOut := 0.0
		for _, out := range tx.MsgTx().TxOut {
			totalOut += btcutil.Amount(out.Value).ToBTC()
		}
		fee := totalIn - totalOut
		feeRate := fee / float64(tx.MsgTx().SerializeSize())
		txInfos = append(txInfos, txInfo{feeRate, tx.MsgTx().SerializeSize()})
	}

	// Sort transactions by fee rate
	sort.Slice(txInfos, func(i, j int) bool {
		return txInfos[i].feeRate > txInfos[j].feeRate
	})

	// Calculate deviation and filter out outliers
	var sum, sumOfSquares float64
	for _, tx := range txInfos {
		sum += tx.feeRate
		sumOfSquares += tx.feeRate * tx.feeRate
	}
	mean := sum / float64(len(txInfos))
	stddev := math.Sqrt(sumOfSquares/float64(len(txInfos)) - mean*mean)
	lowerBound := mean - 2*stddev
	upperBound := mean + 2*stddev
	filteredTxInfos := make([]txInfo, 0, len(txInfos))
	for _, tx := range txInfos {
		if tx.feeRate >= lowerBound && tx.feeRate <= upperBound {
			filteredTxInfos = append(filteredTxInfos, tx)
		}
	}

	// Estimate the fee rate as the n-th highest fee rate
	n := int(float64(len(filteredTxInfos)) * 0.9) // Use the 80th percentile as an example
	if n >= len(filteredTxInfos) {
		n = len(filteredTxInfos) - 1
	}
	feeRateEstimate := filteredTxInfos[n].feeRate

	// Print the estimated fee rate
	log.Printf("Estimated fee rate: %f sat/vB\n", feeRateEstimate*1e8)

	bestFee := feeRateEstimate * 1e8

	return bestFee
}

func GetHalfHourFee() float64 {
	client := utils.Bitcoind()
	defer client.Shutdown()

	// Get the height of the block chain
	_, height, err := client.GetBestBlock()
	if err != nil {
		log.Fatal(err)
	}

	// Get the height of the block chain half an hour ago
	targetHeight := height - int32(6) // 6 blocks per hour
	targetBlockHash, err := client.GetBlockHash(int64(targetHeight))
	if err != nil {
		log.Fatal(err)
	}

	// Get the block half an hour ago
	targetBlock, err := client.GetBlock(targetBlockHash)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the total fees and transaction sizes in the block
	var totalFees btcutil.Amount
	var totalSize int
	for _, tx := range targetBlock.Transactions {
		// Calculate the total input value of the transaction
		var totalIn btcutil.Amount
		for _, in := range tx.TxIn {
			outpoint := in.PreviousOutPoint
			prevTx, err := client.GetRawTransaction(&outpoint.Hash)
			if err != nil {
				log.Fatal(err)
			}
			prevOut := prevTx.MsgTx().TxOut[outpoint.Index]
			totalIn += btcutil.Amount(prevOut.Value)
		}

		// Calculate the total output value of the transaction
		var totalOut btcutil.Amount
		for _, out := range tx.TxOut {
			totalOut += btcutil.Amount(out.Value)
		}

		// Calculate the transaction fee and size
		fee := totalIn - totalOut
		size := tx.SerializeSize()

		totalFees += fee
		totalSize += size
	}

	// Calculate the fee rate for the block
	feeRate := float64(totalFees) / float64(totalSize)

	return feeRate
}
