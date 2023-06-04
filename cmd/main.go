package main

import (
	"log"
	"time"

	mempoolfee "github.com/NonsoAmadi10/mempool-fee/mempool-fee"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func startLoader() chan struct{} {
	loader := make(chan struct{})
	emojis := []string{"🔄", "⏳"}
	go func() {
		for {
			select {
			case <-loader:
				return
			default:
				for _, emoji := range emojis {
					log.Printf("%s Loading...", emoji)
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
	return loader
}

func closeLoader(loader chan struct{}) {
	close(loader)
}

// GetBestFee
var getBestFeeCmd = &cobra.Command{
	Use:   "getbestfee",
	Short: "Get the best fee rate",
	Run: func(cmd *cobra.Command, args []string) {
		loader := startLoader()
		time.Sleep(2 * time.Second)
		closeLoader(loader)
		bestFee := mempoolfee.GetBestFee()
		color.Cyan("Estimated fee rate: %f sat/vB\n", bestFee)
	},
}

// Get Priority Fees

var getPriorityFeeCmd = &cobra.Command{
	Use:   "getpriorityfees",
	Short: "Get fees based on priority",
	Run: func(cmd *cobra.Command, args []string) {
		loader := startLoader()
		time.Sleep(2 * time.Second)
		closeLoader(loader)
		high, normal, low, err := mempoolfee.GetPriorityFees()

		if err != nil {
			log.Fatal(err)
		}
		color.Green("High fee rate: %f sat/vB\n", high)
		color.Blue("Normal fee rate: %f sat/vB\n", normal)
		color.Magenta("Low fee rate: %f sat/vB\n", low)
	},
}

// GetHalfHourFee Rate
var getHalfRate = &cobra.Command{
	Use:   "gethalfhour",
	Short: "Get the fee rate for half-hour blocks",
	Run: func(cmd *cobra.Command, args []string) {
		loader := startLoader()
		time.Sleep(2 * time.Second)
		closeLoader(loader)
		bestFee := mempoolfee.GetHalfHourFee()
		color.HiYellow("Estimated fee rate for half-hour: %f sat/vB\n", bestFee)
	},
}

func main() {
	rootCmd := &cobra.Command{
		Version: "0.0.1",
		Use:     "ily",
		Short:   "mempool fee analyser",
		Long:    "estimates the fee for a Bitcoin transaction by analyzing the current mempool",
	}

	rootCmd.AddCommand(getBestFeeCmd)
	rootCmd.AddCommand(getPriorityFeeCmd)
	rootCmd.AddCommand(getHalfRate)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
