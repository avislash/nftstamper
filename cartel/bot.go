package cartel

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/avislash/nftstamper/cartel/config"
	"github.com/avislash/nftstamper/cartel/image"
	"github.com/avislash/nftstamper/cartel/metadata"
	"github.com/avislash/nftstamper/lib/ipfs"
	"github.com/avislash/nftstamper/lib/log"
	"github.com/avislash/nftstamper/root"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var (
	ipfsClient           ipfs.Client
	maycIpfsClient       ipfs.Client
	stamper              *image.Processor
	houndMetadataFetcher *metadata.HoundMetadataFetcher
	maycMetadataFetcher  *metadata.MAYCMetadataFetcher
	logger               *log.SugaredLogger
	configFile           string
	botToken             string
	env                  string
)

type collectionOpt int

const (
	houndsOpt collectionOpt = iota
	maycOpt
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
	cmd.PersistentFlags().StringVar(&env, "env", "CARTEL", "Configuration Environment")
	root.Cmd.AddCommand(cmd)
}

func botInit(_ *cobra.Command, _ []string) error {
	cfg, err := config.LoadCfg(env, configFile)
	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}

	logger, err = log.NewSugaredLogger() //log.WithLogLevel(log.DEBUG))
	if err != nil {
		return fmt.Errorf("Unable to instantiate logger")
	}

	stamper, err = image.NewProcessor(cfg.ImageProcessorConfig)
	if err != nil {
		return fmt.Errorf("Error initializing Image Processor: %w", err)
	}

	ipfsClient, err = ipfs.NewClient(cfg.IPFSEndpoint, ipfs.WithJPEGDecoder())
	if err != nil {
		return fmt.Errorf("Error creating IPFS Client: %w", err)
	}

	maycIpfsClient, err = ipfs.NewClient(cfg.IPFSEndpoint, ipfs.WithPNGDecoder())
	if err != nil {
		return fmt.Errorf("Error creating IPFS Client: %w", err)
	}
	houndMetadataFetcher = metadata.NewHoundMetadataFetcher(cfg.HoundsMetadataEndpoint)
	maycMetadataFetcher = metadata.NewMAYCMetadataFetcher(cfg.MAYCMetadataEndpoint)
	botToken = "Bot " + cfg.BotToken
	return nil
}

func cartelBot(cmd *cobra.Command, _ []string) error {
	dg, err := discordgo.New(botToken)
	if err != nil {
		return err
	}

	dg.AddHandler(gmInteraction)
	dg.AddHandler(nfdInteraction)
	dg.AddHandler(suitInteraction)
	dg.AddHandler(pledgeInteraction)
	dg.AddHandler(apeBagInteraction)

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

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "nfd",
		Description: "Add NFD Campaign Merch",
		Options:     []*discordgo.ApplicationCommandOption{houndID},
	})
	if err != nil {
		return err
	}

	maxID = float64(30006)
	maycID := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "mayc",
		Description: "MAYC ID #",
		Required:    true,
		MinValue:    &minID,
		MaxValue:    maxID,
	}

	nfdSuitSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "nfd",
		Description: "NFD Suit",
		Options:     []*discordgo.ApplicationCommandOption{maycID},
	}

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "suit",
		Description: "Add Suit to MAYC",
		Options:     []*discordgo.ApplicationCommandOption{nfdSuitSubCmd},
	})
	if err != nil {
		return err
	}

	id := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "id",
		Description: "ID #",
		Required:    true,
		MinValue:    &minID,
		MaxValue:    maxID,
	}
	choices := []*discordgo.ApplicationCommandOptionChoice{&discordgo.ApplicationCommandOptionChoice{Name: "hound", Value: houndsOpt},
		&discordgo.ApplicationCommandOptionChoice{Name: "mayc", Value: maycOpt},
	}
	collectionChoice := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "collection",
		Description: "which collection",
		Required:    true,
		Choices:     choices,
	}
	orangeHandSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "orange",
		Description: "Orange Hand Stamp",
		Options:     []*discordgo.ApplicationCommandOption{collectionChoice, id},
	}
	blackHandSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "black",
		Description: "Black Hand Stamp",
		Options:     []*discordgo.ApplicationCommandOption{collectionChoice, id},
	}
	whiteHandSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "white",
		Description: "White Hand Stamp",
		Options:     []*discordgo.ApplicationCommandOption{collectionChoice, id},
	}
	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "pledge",
		Description: "Pledge MAYC with the Mutant Cartel",
		Options:     []*discordgo.ApplicationCommandOption{orangeHandSubCmd, blackHandSubCmd, whiteHandSubCmd},
	})
	if err != nil {
		return err
	}

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "apebag",
		Description: "Give MAYC a bag of $APE",
		Options:     []*discordgo.ApplicationCommandOption{maycID},
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
				houndID := cmdData.Options[0].UintValue()
				metadata, err := houndMetadataFetcher.Fetch(houndID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for Hound #%d: %w", houndID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve Hound #%d image from IPFS: %w", houndID, err)
					logger.Errorf("Error: %w", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				buff, err := stamper.OverlayBowl(hound, metadata.Background)
				if err != nil {
					err := fmt.Errorf("Failed to create GM image for Hound %d: %w ", houndID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_gm_hound_%d.png", name, houndID),
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

func nfdInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var name string
	if nil == interaction.Member {
		name = interaction.User.Username
	} else {
		if nil != interaction.Member.User {
			name = interaction.Member.User.Username
		}
	}
	if discordgo.InteractionApplicationCommand == interaction.Type {
		cmdData := interaction.ApplicationCommandData()
		if cmdData.Name == "nfd" {
			//Send ACK To meet the 3s turnaround and allow for more time to upload the image
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{}})
			go func() {
				houndID := cmdData.Options[0].UintValue()
				logger.Debugf("Getting metadata for hound #%d", houndID)
				metadata, err := houndMetadataFetcher.Fetch(houndID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for Hound #%d: %w", houndID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}
				logger.Debugf("Metadata: %+v", metadata)

				hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve Hound #%d image from IPFS: %w", houndID, err)
					logger.Errorf("Error: %w", err)
					sendErrorResponse(err, session, interaction)
					return
				}
				logger.Debugf("Got image from IPFS")

				logger.Debugf("Overlaying Image")
				buff, err := stamper.OverlayNFDMerch(hound, metadata)
				if err != nil {
					err := fmt.Errorf("Failed to create NFD Merch image for Hound %d: %w ", houndID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_nfd_merch_hound_%d.png", name, houndID),
					ContentType: "image/png",
					Reader:      buff,
				}

				content := "In NFD we trust"
				response := &discordgo.WebhookEdit{
					Content: &content,
					Files:   []*discordgo.File{file},
				}
				logger.Debugf("Uploading image")
				if _, err := session.InteractionResponseEdit(interaction.Interaction, response); err != nil {
					logger.Errorf("Error sending message: %s", err)
				}
			}()
		}

	}

}

func suitInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var name string
	if nil == interaction.Member {
		name = interaction.User.Username
	} else {
		if nil != interaction.Member.User {
			name = interaction.Member.User.Username
		}
	}
	if discordgo.InteractionApplicationCommand == interaction.Type {
		cmdData := interaction.ApplicationCommandData()
		if cmdData.Name == "suit" {
			//Send ACK To meet the 3s turnaround and allow for more time to upload the image
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{}})
			go func() {
				maycID := cmdData.Options[0].Options[0].UintValue()
				logger.Debugf("Getting metadata for MAYC #%d", maycID)
				metadata, err := maycMetadataFetcher.Fetch(maycID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for MAYC #%d: %w", maycID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}
				logger.Debugf("Metadata: %+v", metadata)

				mayc, err := maycIpfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve MAYC #%d image from IPFS: %w", maycID, err)
					logger.Errorf("Error: %w", err)
					sendErrorResponse(err, session, interaction)
					return
				}
				logger.Debugf("Got image from IPFS")

				logger.Debugf("Overlaying Image")
				buff, err := stamper.OverlayNFDSuit(mayc)
				if err != nil {
					err := fmt.Errorf("Failed to overlay NFD Suit to MAYC  %d: %w ", maycID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_nfd_suit_mayc_%d.png", name, maycID),
					ContentType: "image/png",
					Reader:      buff,
				}

				content := "In NFD we trust"
				response := &discordgo.WebhookEdit{
					Content: &content,
					Files:   []*discordgo.File{file},
				}
				logger.Debugf("Uploading image")
				if _, err := session.InteractionResponseEdit(interaction.Interaction, response); err != nil {
					logger.Errorf("Error sending message: %s", err)
				}
			}()
		}

	}

}

func pledgeInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var name string
	if nil == interaction.Member {
		name = interaction.User.Username
	} else {
		if nil != interaction.Member.User {
			name = interaction.Member.User.Username
		}
	}
	if discordgo.InteractionApplicationCommand == interaction.Type {
		cmdData := interaction.ApplicationCommandData()
		if cmdData.Name == "pledge" {
			//Send ACK To meet the 3s turnaround and allow for more time to upload the image
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{}})
			go func() {
				var (
					filename string
					image    *bytes.Buffer
				)

				options := cmdData.Options[0]
				color := options.Name

				choices := cmdData.Options[0]

				collection := choices.Options[0].UintValue()

				switch collectionOpt(collection) {
				case maycOpt:
					maycID := choices.Options[1].UintValue()
					filename = fmt.Sprintf("%s_pledge_mayc_%d.png", name, maycID)
					logger.Debugf("Getting metadata for MAYC #%d", maycID)
					metadata, err := maycMetadataFetcher.Fetch(maycID)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve metadata for MAYC #%d: %w", maycID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(err, session, interaction)
						return
					}
					logger.Debugf("Metadata: %+v", metadata)

					mayc, err := maycIpfsClient.GetImageFromIPFS(metadata.Image)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve MAYC #%d image from IPFS: %w", maycID, err)
						logger.Errorf("Error: %w", err)
						sendErrorResponse(err, session, interaction)
						return
					}
					logger.Debugf("Got image from IPFS")

					logger.Debugf("Overlaying Image")
					image, err = stamper.OverlayHandMAYC(mayc, metadata, color)
					if err != nil {
						err := fmt.Errorf("Failed to overlay Hand Stamp to MAYC  %d: %w ", maycID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(err, session, interaction)
						return
					}
				case houndsOpt:
					houndID := choices.Options[1].UintValue()
					filename = fmt.Sprintf("%s_pledge_hound_%d.png", name, houndID)
					logger.Debugf("Getting metadata for Hound #%d", houndID)
					metadata, err := houndMetadataFetcher.Fetch(houndID)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve metadata for Hound #%d: %w", houndID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(err, session, interaction)
						return
					}
					logger.Debugf("Metadata: %+v", metadata)

					hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve Hound #%d image from IPFS: %w", houndID, err)
						logger.Errorf("Error: %w", err)
						sendErrorResponse(err, session, interaction)
						return
					}
					logger.Debugf("Got image from IPFS")

					logger.Debugf("Overlaying Image")
					image, err = stamper.OverlayHandHound(hound, metadata, color)
					if err != nil {
						err := fmt.Errorf("Failed to overlay Hand Stamp on Hound  %d: %w ", houndID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(err, session, interaction)
						return
					}

				default:
					sendErrorResponse(fmt.Errorf("Unrecognized Colection: %d", collection), session, interaction)
					return

				}
				file := &discordgo.File{
					Name:        filename,
					ContentType: "image/png",
					Reader:      image,
				}

				content := "I swear by the Apes of old and by all that is sacred to Mutants that I stand with the Mutant Cartel"
				response := &discordgo.WebhookEdit{
					Content: &content,
					Files:   []*discordgo.File{file},
				}
				logger.Debugf("Uploading image")
				if _, err := session.InteractionResponseEdit(interaction.Interaction, response); err != nil {
					logger.Errorf("Error sending message: %s", err)
				}
			}()
		}

	}

}

func apeBagInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var name string
	if nil == interaction.Member {
		name = interaction.User.Username
	} else {
		if nil != interaction.Member.User {
			name = interaction.Member.User.Username
		}
	}
	if discordgo.InteractionApplicationCommand == interaction.Type {
		cmdData := interaction.ApplicationCommandData()
		if cmdData.Name == "apebag" {
			//Send ACK To meet the 3s turnaround and allow for more time to upload the image
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{}})
			go func() {
				maycID := cmdData.Options[0].UintValue()
				logger.Debugf("Getting metadata for MAYC #%d", maycID)
				metadata, err := maycMetadataFetcher.Fetch(maycID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for MAYC #%d: %w", maycID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}
				logger.Debugf("Metadata: %+v", metadata)

				mayc, err := maycIpfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve MAYC #%d image from IPFS: %w", maycID, err)
					logger.Errorf("Error: %w", err)
					sendErrorResponse(err, session, interaction)
					return
				}
				logger.Debugf("Got image from IPFS")

				logger.Debugf("Overlaying Image")
				buff, err := stamper.OverlayApeBag(mayc, metadata)
				if err != nil {
					err := fmt.Errorf("Failed to overlay Ape Bag to MAYC  %d: %w ", maycID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(err, session, interaction)
					return
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_apebag_mayc_%d.png", name, maycID),
					ContentType: "image/png",
					Reader:      buff,
				}

				response := &discordgo.WebhookEdit{
					Files: []*discordgo.File{file},
				}
				logger.Debugf("Uploading image")
				if _, err := session.InteractionResponseEdit(interaction.Interaction, response); err != nil {
					logger.Errorf("Error sending message: %s", err)
				}
			}()
		}

	}

}

func sendErrorResponse(err error, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	errMsg := err.Error()

	if strings.Contains(errMsg, "invalid character 'T' looking for beginning of value") {
		options := interaction.ApplicationCommandData().Options[0]
		if len(options.Options) != 0 {
			options = options.Options[0]
		}
		houndID := options.UintValue()
		errMsg = fmt.Sprintf("Error: Hound #%d has not yet been revealed", houndID)
	}

	if strings.Contains(errMsg, "invalid JPEG format") {
		options := interaction.ApplicationCommandData().Options[0]
		if len(options.Options) != 0 {
			options = options.Options[0]
		}
		houndID := options.UintValue()
		errMsg = fmt.Sprintf("Error: Is Hound #%d a Mega? Megas are not currently supported", houndID)
	}

	response := &discordgo.WebhookParams{
		Content: errMsg,
		Flags:   discordgo.MessageFlagsEphemeral,
	}

	if err := session.InteractionResponseDelete(interaction.Interaction); err != nil {
		logger.Errorf("Failed to delete interaction: %s", err)
	}

	if _, err := session.FollowupMessageCreate(interaction.Interaction, true, response); err != nil {
		logger.Errorf("Error sending message: %s", err)
	}

}
