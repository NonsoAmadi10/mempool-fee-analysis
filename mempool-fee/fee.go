package mempoolfee

import (
	"log"
	"math"
	"sort"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

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

	// Get the current mempool state.
	txs, err := client.GetRawMempool()
	if err != nil {
		log.Fatal(err)
	}

	// Calculate fee rates for each transaction in the mempool.
	feeRates := make([]float64, len(txs))
	for i, txHash := range txs {
		txVerbose, err := client.GetRawTransactionVerbose(txHash)
		if err != nil {
			log.Fatal(err)
		}

		var totalIn, totalOut float64
		for _, vin := range txVerbose.Vin {
			txHashFromStr, _ := chainhash.NewHashFromStr(vin.Txid)
			txInVerbose, err := client.GetRawTransactionVerbose(txHashFromStr)
			if err != nil {
				log.Fatal(err)
			}
			totalIn += txInVerbose.Vout[vin.Vout].Value
		}

		for _, vout := range txVerbose.Vout {
			totalOut += vout.Value
		}

		fee := totalIn - totalOut
		size := float64(txVerbose.Size)
		feeRates[i] = float64(fee) / size
	}

	// Sort transactions in descending order based on fee rates.
	sort.SliceStable(feeRates, func(i, j int) bool {
		return feeRates[i] > feeRates[j]
	})

	// Calculate median fee rate of top 50 transactions.
	var totalFeeRate float64
	n := len(feeRates)
	if n < 50 {
		log.Printf("Warning: Mempool contains only %d transactions, expected at least 50", n)
		n = 50
	}
	for i := 0; i < n && i < 50; i++ {
		totalFeeRate += feeRates[i]
	}
	halfhourFeeRate := totalFeeRate / 50.0

	return halfhourFeeRate * 1e8

}
