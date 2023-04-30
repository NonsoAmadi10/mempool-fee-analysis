package mempoolfee

import (
	"log"
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
