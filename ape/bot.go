package ape

import (
	"fmt"

	"github.com/avislash/nftstamper/ape/config"
	"github.com/avislash/nftstamper/ape/image"
	"github.com/avislash/nftstamper/ape/metadata"
	"github.com/avislash/nftstamper/lib/ipfs"
	"github.com/avislash/nftstamper/lib/log"
	"github.com/avislash/nftstamper/root"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var (
	ipfsClient      ipfs.Client
	stamper         *image.Processor
	metadataFetcher *metadata.SentinelMetadataFetcher
	logger          *log.SugaredLogger
	configFile      string
	botToken        string
	env             string
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
	cmd.PersistentFlags().StringVar(&env, "env", "APE", "Configuration Environment")
	root.Cmd.AddCommand(cmd)
}

func botInit(_ *cobra.Command, _ []string) error {
	cfg, err := config.LoadCfg(env, configFile)
	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}

	logger, err = log.NewSugaredLogger(log.WithLogLevel(log.Level(cfg.LogLevel)))
	if err != nil {
		return fmt.Errorf("Unable to instantiate logger")
	}

	stamper, err = image.NewProcessor(cfg.ImageProcessorConfig, logger)
	if err != nil {
		return fmt.Errorf("Error initializing Image Processor: %w", err)
	}

	ipfsClient, err = ipfs.NewClient(cfg.IPFSEndpoint, ipfs.WithPNGDecoder())
	if err != nil {
		return fmt.Errorf("Error creating IPFS Client: %w", err)
	}

	metadataFetcher = metadata.NewSentinelMetadataFetcher(cfg.MetadataEndpoint)
	botToken = "Bot " + cfg.BotToken
	return nil
}

func apeBot(cmd *cobra.Command, _ []string) error {
	dg, err := discordgo.New(botToken)
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
	maxID := float64(9999)
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

	logger.Info("Bot started")
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
			//Send ACK To meet the 3s turnaround and allow for more time to upload the image
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{}})
			go func() {
				sentinelID := cmdData.Options[0].UintValue()

				metadata, err := metadataFetcher.Fetch(sentinelID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for Sentinel #%d: %w", sentinelID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				sentinel, err := ipfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve Sentinel #%d image from IPFS: %w", sentinelID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				buff, err := stamper.OverlayMug(sentinel, metadata)
				if err != nil {
					err := fmt.Errorf("Failed to create GM image for Sentinel %d: %w ", sentinelID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_gm_sentinel_%d.png", name, sentinelID),
					ContentType: "image/png",
					Reader:      buff,
				}

				content := "GM " + mention
				response := &discordgo.WebhookEdit{
					Content: &content,
					Files:   []*discordgo.File{file},
				}
				if _, err := session.InteractionResponseEdit(interaction.Interaction, response); err != nil {
					logger.Errorf("Error sending message: %s", err)
				}
			}()
		}

	}
}

func sendErrorResponse(err error, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	response := &discordgo.WebhookParams{
		Content: err.Error(),
		Flags:   discordgo.MessageFlagsEphemeral,
	}

	if err := session.InteractionResponseDelete(interaction.Interaction); err != nil {
		logger.Errorf("Failed to delete interaction: %s", err)
	}

	if _, err := session.FollowupMessageCreate(interaction.Interaction, true, response); err != nil {
		logger.Errorf("Error sending message: %s", err)
	}
}
