package ape

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/avislash/nftstamper/ape/image"
	"github.com/avislash/nftstamper/ape/metadata"
	"github.com/avislash/nftstamper/config"
	"github.com/avislash/nftstamper/lib/ipfs"
	"github.com/avislash/nftstamper/root"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ipfsClient      *ipfs.Client
	stamper         *image.Processor
	metadataFetcher *metadata.SentinelMetadataFetcher
	configFile      string
)

var cmd = &cobra.Command{
	Use:     "apebot",
	Short:   "Instantiate NFT Stamper for supported Applied Primate Engineering Collections",
	Long:    "Instantiate NFT Stamper for supported Applied Primate Engineering Collections",
	PreRunE: botInit,
	RunE:    apeBot,
}

func init() {
	cmd.PersistentFlags().StringVar(&configFile, "config", "./ape/config.yaml", "Path to config file")
	root.Cmd.AddCommand(cmd)
}

func botInit(_ *cobra.Command, _ []string) error {
	var configParams config.Config
	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("Failed to read in config.yaml: %w", err)
	}

	err = yaml.Unmarshal(configFile, &configParams)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config.yaml: %w", err)
	}

	stamper, err = image.NewProcessor(configParams.ImageProcessorConfig)
	if err != nil {
		return fmt.Errorf("Error initializing Image Processor: %w", err)
	}

	ipfsClient, err = ipfs.NewClient(configParams.IPFSEndpoint, ipfs.WithPNGDecoder())
	if err != nil {
		return fmt.Errorf("Error creating IPFS Client: %w", err)
	}

	metadataFetcher = metadata.NewSentinelMetadataFetcher(configParams.MetadataEndpoint)
	return nil
}

func apeBot(cmd *cobra.Command, _ []string) error {
	token := os.Getenv("APE_DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	dg.AddHandler(gmInteraction)
	if err := dg.Open(); err != nil {
		return err
	}
	defer dg.Close()
	botID := dg.State.User.ID

	minID := float64(0)
	maxID := float64(10000)
	sentinelID := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "sentinel",
		Description: "Sentinel ID #",
		Required:    true,
		MinValue:    &minID,
		MaxValue:    maxID,
	}
	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "gm",
		Description: "Responds with a GM",
		Options:     []*discordgo.ApplicationCommandOption{sentinelID},
	})

	if err != nil {
		return err
	}

	log.Println("Bot started")
	<-cmd.Context().Done()
	return nil
}

func gmInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var mention string
	var name string
	if nil == interaction.Member {
		mention = interaction.User.Mention()
		name = interaction.User.Username
	} else {
		mention = interaction.Member.Mention()
		if nil != interaction.Member.User {
			name = interaction.Member.User.Username
		}
	}
	if discordgo.InteractionApplicationCommand == interaction.Type {
		cmdData := interaction.ApplicationCommandData()
		if cmdData.Name == "gm" {
			go func() {
				sentinelID := cmdData.Options[0].UintValue()

				metadata, err := metadataFetcher.Fetch(sentinelID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for Sentienl #%d: %w", sentinelID, err)
					log.Println("Error: ", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				sentinel, err := ipfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve Sentinel #%d image from IPFS: %w", sentinelID, err)
					log.Println("Error: ", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				buff, err := stamper.OverlayMug(sentinel, metadata.BaseArmor)
				if err != nil {
					err := fmt.Errorf("Failed to create GM image for Sentinel %d: %w ", sentinelID, err)
					log.Println("Error: ", err)
					sendErrorResponse(err, session, interaction)
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_gm_sentinel_%d.png", name, sentinelID),
					ContentType: "image/png",
					Reader:      buff,
				}
				response := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "GM " + mention,
						Files:   []*discordgo.File{file},
					},
				}

				if err := session.InteractionRespond(interaction.Interaction, response); err != nil {
					log.Println("Error sending message: ", err)
				}
			}()
		}

	}
}

func sendErrorResponse(err error, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: err.Error(),
		},
	}

	if err := session.InteractionRespond(interaction.Interaction, response); err != nil {
		log.Println("Error sending message: ", err)
	}

}
