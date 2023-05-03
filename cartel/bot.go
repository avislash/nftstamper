package cartel

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/avislash/nftstamper/cartel/image"
	"github.com/avislash/nftstamper/cartel/metadata"
	"github.com/avislash/nftstamper/config"
	libImg "github.com/avislash/nftstamper/lib/image"
	"github.com/avislash/nftstamper/lib/ipfs"
	"github.com/avislash/nftstamper/root"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ipfsClient      *ipfs.Client
	stamper         *image.Processor
	metadataFetcher *metadata.HoundMetadataFetcher
	configFile      string
)

var cmd = &cobra.Command{
	Use:     "cartelbot",
	Short:   "Instantiate NFT Stamper for supported Mutant Cartel Collections",
	Long:    "Instantiate NFT Stamper for supported Mutant Cartel Collections",
	PreRunE: botInit,
	RunE:    cartelBot,
}

func init() {
	cmd.PersistentFlags().StringVar(&configFile, "config", "./cartel/config.yaml", "Path to config file")
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

	ipfsClient, err = ipfs.NewClient(&libImg.JPEGDecoder{})
	if err != nil {
		return fmt.Errorf("Error creating IPFS Client: %w", err)
	}

	metadataFetcher = metadata.NewHoundMetadataFetcher(configParams.MetadataEndpoint)
	return nil
}

func cartelBot(cmd *cobra.Command, _ []string) error {
	token := os.Getenv("CARTEL_DISCORD_BOT_TOKEN")
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
	houndID := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "hound",
		Description: "Hound ID #",
		Required:    true,
		MinValue:    &minID,
		MaxValue:    maxID,
	}
	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "gm",
		Description: "Responds with a GM",
		Options:     []*discordgo.ApplicationCommandOption{houndID},
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
				houndID := cmdData.Options[0].UintValue()
				metadata, err := metadataFetcher.Fetch(houndID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for Hound #%d: %w", houndID, err)
					log.Println("Error: ", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve Hound #%d image from IPFS: %w", houndID, err)
					log.Println("Error: ", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				buff, err := stamper.OverlayBowl(hound, metadata.Background)
				if err != nil {
					err := fmt.Errorf("Failed to create GM image for Hound %d: %w ", houndID, err)
					log.Println("Error: ", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_gm_hound_%d.png", name, houndID),
					ContentType: "image/png",
					Reader:      buff,
				}

				//Send ACK To meet the 3s turnaround and allow for more time to upload the image
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{}})

				content := "GM " + mention
				response := &discordgo.WebhookEdit{
					Content: &content,
					Files:   []*discordgo.File{file},
				}
				if _, err := session.InteractionResponseEdit(interaction.Interaction, response); err != nil {
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
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	if err := session.InteractionRespond(interaction.Interaction, response); err != nil {
		log.Println("Error sending message: ", err)
	}

}
