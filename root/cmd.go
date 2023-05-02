package root

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "nftstamper",
	Short: "Instantiate the NFT Stamper Discord bot for a supported collection",
	Long:  "Instantiate the NFT Stamper Discord bot for a supported collection",
}
