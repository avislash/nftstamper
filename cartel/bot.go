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
	dg.AddHandler(jerseyInteraction)

	if err := dg.Open(); err != nil {
		return err
	}

	defer dg.Close()
	botID := dg.State.User.ID

	minID := float64(0)
	maxHoundID := float64(10000)
	maxMAYCID := float64(30006)

	houndID := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "hound",
		Description: "Hound ID #",
		Required:    true,
		MinValue:    &minID,
		MaxValue:    maxHoundID,
	}
	maycID := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "mayc",
		Description: "MAYC ID #",
		Required:    true,
		MinValue:    &minID,
		MaxValue:    maxMAYCID,
	}

	//TODO make cfg a global and just grab the values programtically from the map
	liquidChoices := []*discordgo.ApplicationCommandOptionChoice{
		&discordgo.ApplicationCommandOptionChoice{Name: "coffee", Value: "coffee"},
		&discordgo.ApplicationCommandOptionChoice{Name: "serum", Value: "serum"},
	}

	maycGMCmdLiquidChoices := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "liquid",
		Description: "Coffe Mug Liquid",
		Required:    true,
		Choices:     liquidChoices,
	}

	logoChoices := []*discordgo.ApplicationCommandOptionChoice{
		&discordgo.ApplicationCommandOptionChoice{Name: "albino asylum", Value: "albino asylum"},
		&discordgo.ApplicationCommandOptionChoice{Name: "armoured guards", Value: "armoured guards"},
		&discordgo.ApplicationCommandOptionChoice{Name: "bionic army", Value: "bionic army"},
		&discordgo.ApplicationCommandOptionChoice{Name: "blood hounds", Value: "blood hounds"},
		&discordgo.ApplicationCommandOptionChoice{Name: "cartel", Value: "cartel"},
		&discordgo.ApplicationCommandOptionChoice{Name: "death pack", Value: "death pack"},
		&discordgo.ApplicationCommandOptionChoice{Name: "deathbot army", Value: "deathbot army"},
		&discordgo.ApplicationCommandOptionChoice{Name: "demon council", Value: "demon council"},
		&discordgo.ApplicationCommandOptionChoice{Name: "dmt cartel", Value: "dmt cartel"},
		&discordgo.ApplicationCommandOptionChoice{Name: "flesh eaters", Value: "flesh eaters"},
		&discordgo.ApplicationCommandOptionChoice{Name: "golem gang", Value: "golem gang"},
		&discordgo.ApplicationCommandOptionChoice{Name: "haunted howlers", Value: "haunted howlers"},
		&discordgo.ApplicationCommandOptionChoice{Name: "laughing legion", Value: "laughing legion"},
		&discordgo.ApplicationCommandOptionChoice{Name: "metal militia", Value: "metal militia"},
		&discordgo.ApplicationCommandOptionChoice{Name: "midnight marauders", Value: "midnight marauders"},
		&discordgo.ApplicationCommandOptionChoice{Name: "noisy syndicate", Value: "noisy syndicate"},
		&discordgo.ApplicationCommandOptionChoice{Name: "royal hounds", Value: "royal hounds"},
		&discordgo.ApplicationCommandOptionChoice{Name: "skull legion", Value: "skull legion"},
		&discordgo.ApplicationCommandOptionChoice{Name: "trippy brigade", Value: "trippy brigade"},
		&discordgo.ApplicationCommandOptionChoice{Name: "wolf pack", Value: "wolf pack"},
		&discordgo.ApplicationCommandOptionChoice{Name: "zombie horde", Value: "zombie horde"},
	}

	maycGMCmdLogoChoices := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "logo",
		Description: "Coffe Mug Logo",
		Required:    true,
		Choices:     logoChoices,
	}

	maycGmSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "mayc",
		Description: "Responds with a MAYC GM",
		Options:     []*discordgo.ApplicationCommandOption{maycID, maycGMCmdLiquidChoices, maycGMCmdLogoChoices},
	}

	houndGmSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "hound",
		Description: "Responds with a Mutant Hound GM",
		Options:     []*discordgo.ApplicationCommandOption{houndID},
	}

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "gm",
		Description: "Responds with a GM",
		Options:     []*discordgo.ApplicationCommandOption{houndGmSubCmd, maycGmSubCmd},
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

	id := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "id",
		Description: "ID #",
		Required:    true,
		MinValue:    &minID,
		MaxValue:    maxMAYCID,
	}
	choices := []*discordgo.ApplicationCommandOptionChoice{&discordgo.ApplicationCommandOptionChoice{Name: "hound", Value: houndsOpt},
		&discordgo.ApplicationCommandOptionChoice{Name: "mayc", Value: maycOpt},
	}
	collectionChoices := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "collection",
		Description: "which collection",
		Required:    true,
		Choices:     choices,
	}

	pledgeColorChoices := []*discordgo.ApplicationCommandOptionChoice{
		&discordgo.ApplicationCommandOptionChoice{Name: "orange", Value: "orange"},
		&discordgo.ApplicationCommandOptionChoice{Name: "black", Value: "black"},
		&discordgo.ApplicationCommandOptionChoice{Name: "white", Value: "white"},
	}
	pledgeCmdColorChoices := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "color",
		Description: "Pledge Hand Color",
		Required:    true,
		Choices:     pledgeColorChoices,
	}

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "pledge",
		Description: "Pledge with the Mutant Cartel",
		Options:     []*discordgo.ApplicationCommandOption{pledgeCmdColorChoices, collectionChoices, id},
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

	jerseyCollectionChoices := []*discordgo.ApplicationCommandOptionChoice{&discordgo.ApplicationCommandOptionChoice{Name: "hound", Value: houndsOpt}}
	jerseyCmdCollectionChoices := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "collection",
		Description: "which collection",
		Required:    true,
		Choices:     jerseyCollectionChoices,
	}

	nflAFCTeams := []*discordgo.ApplicationCommandOptionChoice{
		&discordgo.ApplicationCommandOptionChoice{Name: "cartel", Value: "cartel"},
		&discordgo.ApplicationCommandOptionChoice{Name: "bengals", Value: "bengals"},
		&discordgo.ApplicationCommandOptionChoice{Name: "bills", Value: "bills"},
		&discordgo.ApplicationCommandOptionChoice{Name: "bronocs", Value: "broncos"},
		&discordgo.ApplicationCommandOptionChoice{Name: "browns", Value: "browns"},
		&discordgo.ApplicationCommandOptionChoice{Name: "chiefs", Value: "chiefs"},
		&discordgo.ApplicationCommandOptionChoice{Name: "chargers", Value: "chargers"},
		&discordgo.ApplicationCommandOptionChoice{Name: "colts", Value: "colts"},
		&discordgo.ApplicationCommandOptionChoice{Name: "dolphins", Value: "dolphins"},
		&discordgo.ApplicationCommandOptionChoice{Name: "jaguars", Value: "jaguars"},
		&discordgo.ApplicationCommandOptionChoice{Name: "jets", Value: "jets"},
		&discordgo.ApplicationCommandOptionChoice{Name: "patriots", Value: "patriots"},
		&discordgo.ApplicationCommandOptionChoice{Name: "raiders", Value: "raiders"},
		&discordgo.ApplicationCommandOptionChoice{Name: "ravens", Value: "ravens"},
		&discordgo.ApplicationCommandOptionChoice{Name: "steelers", Value: "steelers"},
		&discordgo.ApplicationCommandOptionChoice{Name: "texans", Value: "texans"},
		&discordgo.ApplicationCommandOptionChoice{Name: "titans", Value: "titans"},
	}

	nflNFCTeams := []*discordgo.ApplicationCommandOptionChoice{
		&discordgo.ApplicationCommandOptionChoice{Name: "cartel", Value: "cartel"},
		&discordgo.ApplicationCommandOptionChoice{Name: "49ers", Value: "49ers"},
		&discordgo.ApplicationCommandOptionChoice{Name: "bears", Value: "bears"},
		&discordgo.ApplicationCommandOptionChoice{Name: "buccaneers", Value: "buccaneers"},
		&discordgo.ApplicationCommandOptionChoice{Name: "cardinals", Value: "cardinals"},
		&discordgo.ApplicationCommandOptionChoice{Name: "commanders", Value: "commanders"},
		&discordgo.ApplicationCommandOptionChoice{Name: "cowboys", Value: "cowboys"},
		&discordgo.ApplicationCommandOptionChoice{Name: "eagles", Value: "eagles"},
		&discordgo.ApplicationCommandOptionChoice{Name: "falcons", Value: "falcons"},
		&discordgo.ApplicationCommandOptionChoice{Name: "giants", Value: "giants"},
		&discordgo.ApplicationCommandOptionChoice{Name: "lions", Value: "lions"},
		&discordgo.ApplicationCommandOptionChoice{Name: "packers", Value: "packers"},
		&discordgo.ApplicationCommandOptionChoice{Name: "panthers", Value: "panthers"},
		&discordgo.ApplicationCommandOptionChoice{Name: "rams", Value: "rams"},
		&discordgo.ApplicationCommandOptionChoice{Name: "saints", Value: "saints"},
		&discordgo.ApplicationCommandOptionChoice{Name: "seahawks", Value: "seahawks"},
		&discordgo.ApplicationCommandOptionChoice{Name: "vikings", Value: "vikings"},
	}

	jerseyCmdAFCTeamChoices := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "team",
		Description: "which team",
		Required:    true,
		Choices:     nflAFCTeams,
	}

	jerseyCmdNFCTeamChoices := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "team",
		Description: "which team",
		Required:    true,
		Choices:     nflNFCTeams,
	}

	nflJerseyAFCSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "nfl-afc",
		Description: "Overlay an NFL AFC Jersey",
		Options:     []*discordgo.ApplicationCommandOption{jerseyCmdCollectionChoices, jerseyCmdAFCTeamChoices, houndID},
	}

	nflJerseyNFCSubCmd := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "nfl-nfc",
		Description: "Overlay an NFL NFC Jersey",
		Options:     []*discordgo.ApplicationCommandOption{jerseyCmdCollectionChoices, jerseyCmdNFCTeamChoices, houndID},
	}

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "jersey",
		Description: "Don a jersey representign your favorite sports team",
		Options:     []*discordgo.ApplicationCommandOption{nflJerseyAFCSubCmd, nflJerseyNFCSubCmd},
	})

	if err != nil {
		return err
	}

	suitChoices := []*discordgo.ApplicationCommandOptionChoice{
		&discordgo.ApplicationCommandOptionChoice{Name: "ape", Value: "ape"},
		&discordgo.ApplicationCommandOptionChoice{Name: "brown", Value: "brown"},
		&discordgo.ApplicationCommandOptionChoice{Name: "cartel", Value: "cartel"},
		&discordgo.ApplicationCommandOptionChoice{Name: "cartel comic", Value: "cartel comic"},
		&discordgo.ApplicationCommandOptionChoice{Name: "cheetah", Value: "cheetah"},
		&discordgo.ApplicationCommandOptionChoice{Name: "demon", Value: "demon"},
		&discordgo.ApplicationCommandOptionChoice{Name: "kodamara", Value: "kodamara"},
		&discordgo.ApplicationCommandOptionChoice{Name: "luke", Value: "luke"},
		&discordgo.ApplicationCommandOptionChoice{Name: "mayc", Value: "mayc"},
		&discordgo.ApplicationCommandOptionChoice{Name: "nfd", Value: "nfd"},
		&discordgo.ApplicationCommandOptionChoice{Name: "red hat", Value: "red hat"},
		&discordgo.ApplicationCommandOptionChoice{Name: "roc", Value: "roc"},
		&discordgo.ApplicationCommandOptionChoice{Name: "trippy", Value: "trippy"},
		&discordgo.ApplicationCommandOptionChoice{Name: "tux", Value: "tux"},
	}

	suitOption := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "style",
		Description: "Suit Style",
		Required:    true,
		Choices:     suitChoices,
	}

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "suit",
		Description: "Outfit MAYC with a stylish suit",
		Options:     []*discordgo.ApplicationCommandOption{suitOption, maycID},
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
				subCmd := cmdData.Options[0]
				var gmFile *discordgo.File
				switch subCmd.Name {
				case "hound":
					houndID := subCmd.Options[0].UintValue()
					metadata, err := houndMetadataFetcher.Fetch(houndID)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve metadata for Hound #%d: %w", houndID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(houndID, err, session, interaction)
						return
					}

					hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve Hound #%d image from IPFS: %w", houndID, err)
						logger.Errorf("Error: %w", err)
						sendErrorResponse(houndID, err, session, interaction)
						return
					}

					buff, err := stamper.OverlayBowl(hound, metadata.Background)
					if err != nil {
						err := fmt.Errorf("Failed to create GM image for Hound %d: %w ", houndID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(houndID, err, session, interaction)
						return
					}

					gmFile = &discordgo.File{
						Name:        fmt.Sprintf("%s_gm_hound_%d.png", name, houndID),
						ContentType: "image/png",
						Reader:      buff,
					}

				case "mayc":
					maycID := subCmd.Options[0].UintValue()
					liquid := subCmd.Options[1].StringValue()
					logo := subCmd.Options[2].StringValue()
					metadata, err := maycMetadataFetcher.Fetch(maycID)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve metadata for MAYC #%d: %w", maycID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(maycID, err, session, interaction)
						return
					}

					mayc, err := maycIpfsClient.GetImageFromIPFS(metadata.Image)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve MAYC #%d image from IPFS: %w", maycID, err)
						logger.Errorf("Error: %w", err)
						sendErrorResponse(maycID, err, session, interaction)
						return
					}

					buff, err := stamper.OverlayCoffeeMug(mayc, metadata, liquid, logo)
					if err != nil {
						err := fmt.Errorf("Failed to create GM image for MAYC %d: %w", maycID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(maycID, err, session, interaction)
					}

					gmFile = &discordgo.File{
						Name:        fmt.Sprintf("%s_gm_mayc%d_with_%s_%s_mug.png", name, maycID, logo, liquid),
						ContentType: "image/png",
						Reader:      buff,
					}
				default:
					logger.Errorf("GM Interaction called with unrecognized sub command: %s", subCmd.Name)
					sendErrorResponse(0, fmt.Errorf("Unrecognized GM sub command option %s", subCmd.Name), session, interaction)
				}

				content := "GM " + mention
				response := &discordgo.WebhookEdit{
					Content: &content,
					Files:   []*discordgo.File{gmFile},
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
					sendErrorResponse(houndID, err, session, interaction)
					return
				}
				logger.Debugf("Metadata: %+v", metadata)

				hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve Hound #%d image from IPFS: %w", houndID, err)
					logger.Errorf("Error: %w", err)
					sendErrorResponse(houndID, err, session, interaction)
					return
				}
				logger.Debugf("Got image from IPFS")

				logger.Debugf("Overlaying Image")
				buff, err := stamper.OverlayNFDMerch(hound, metadata)
				if err != nil {
					err := fmt.Errorf("Failed to create NFD Merch image for Hound %d: %w ", houndID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(houndID, err, session, interaction)
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
				suit := cmdData.Options[0].StringValue()
				maycID := cmdData.Options[1].UintValue()
				logger.Debugf("Getting metadata for MAYC #%d", maycID)
				metadata, err := maycMetadataFetcher.Fetch(maycID)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve metadata for MAYC #%d: %w", maycID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(maycID, err, session, interaction)
					return
				}
				logger.Debugf("Metadata: %+v", metadata)

				mayc, err := maycIpfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve MAYC #%d image from IPFS: %w", maycID, err)
					logger.Errorf("Error: %w", err)
					sendErrorResponse(maycID, err, session, interaction)
					return
				}
				logger.Debugf("Got image from IPFS")

				logger.Debugf("Overlaying Image")
				buff, err := stamper.OverlaySuit(suit, mayc, metadata)
				if err != nil {
					err := fmt.Errorf("Failed to overlay NFD Suit to MAYC  %d: %w ", maycID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(maycID, err, session, interaction)
					return
				}

				file := &discordgo.File{
					Name:        fmt.Sprintf("%s_%s_suit_style_mayc_%d.png", name, suit, maycID),
					ContentType: "image/png",
					Reader:      buff,
				}

				var content string
				switch suit {
				case "brown":
					content = "Well :poop:"
				case "cartel", "cartel comic":
					content = "I swear by the Apes of old and by all that is sacred to Mutants that I stand with the Mutant Cartel"
				case "cheetah":
					content = "Fast AF boi"
				case "demon":
					content = "What a handsome devil :smiling_imp:"
				case "kodamara":
					content = "Wtf is a koda?"
				case "nfd":
					content = "In NFD we trust"
				case "roc":
					content = "F*ck It"
				}

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

				color := cmdData.Options[0].StringValue()

				collection := cmdData.Options[1].UintValue()

				switch collectionOpt(collection) {
				case maycOpt:
					maycID := cmdData.Options[2].UintValue()
					filename = fmt.Sprintf("%s_pledge_mayc_%d.png", name, maycID)
					logger.Debugf("Getting metadata for MAYC #%d", maycID)
					metadata, err := maycMetadataFetcher.Fetch(maycID)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve metadata for MAYC #%d: %w", maycID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(maycID, err, session, interaction)
						return
					}
					logger.Debugf("Metadata: %+v", metadata)

					mayc, err := maycIpfsClient.GetImageFromIPFS(metadata.Image)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve MAYC #%d image from IPFS: %w", maycID, err)
						logger.Errorf("Error: %w", err)
						sendErrorResponse(maycID, err, session, interaction)
						return
					}
					logger.Debugf("Got image from IPFS")

					logger.Debugf("Overlaying Image")
					image, err = stamper.OverlayHandMAYC(mayc, metadata, color)
					if err != nil {
						err := fmt.Errorf("Failed to overlay Hand Stamp to MAYC  %d: %w ", maycID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(maycID, err, session, interaction)
						return
					}
				case houndsOpt:
					//				houndID := choices.Options[1].UintValue()
					houndID := cmdData.Options[2].UintValue()
					filename = fmt.Sprintf("%s_pledge_hound_%d.png", name, houndID)
					logger.Debugf("Getting metadata for Hound #%d", houndID)
					metadata, err := houndMetadataFetcher.Fetch(houndID)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve metadata for Hound #%d: %w", houndID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(houndID, err, session, interaction)
						return
					}
					logger.Debugf("Metadata: %+v", metadata)

					hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
					if err != nil {
						err := fmt.Errorf("Failed to retrieve Hound #%d image from IPFS: %w", houndID, err)
						logger.Errorf("Error: %w", err)
						sendErrorResponse(houndID, err, session, interaction)
						return
					}
					logger.Debugf("Got image from IPFS")

					logger.Debugf("Overlaying Image")
					image, err = stamper.OverlayHandHound(hound, metadata, color)
					if err != nil {
						err := fmt.Errorf("Failed to overlay Hand Stamp on Hound  %d: %w ", houndID, err)
						logger.Errorf("Error: %s", err)
						sendErrorResponse(houndID, err, session, interaction)
						return
					}

				default:
					sendErrorResponse(0, fmt.Errorf("Unrecognized Colection: %d", collection), session, interaction)
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
					sendErrorResponse(maycID, err, session, interaction)
					return
				}
				logger.Debugf("Metadata: %+v", metadata)

				mayc, err := maycIpfsClient.GetImageFromIPFS(metadata.Image)
				if err != nil {
					err := fmt.Errorf("Failed to retrieve MAYC #%d image from IPFS: %w", maycID, err)
					logger.Errorf("Error: %w", err)
					sendErrorResponse(maycID, err, session, interaction)
					return
				}
				logger.Debugf("Got image from IPFS")

				logger.Debugf("Overlaying Image")
				buff, err := stamper.OverlayApeBag(mayc, metadata)
				if err != nil {
					err := fmt.Errorf("Failed to overlay Ape Bag to MAYC  %d: %w ", maycID, err)
					logger.Errorf("Error: %s", err)
					sendErrorResponse(maycID, err, session, interaction)
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

func jerseyInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var name string //TODO resolving the name can be its own method as well
	if nil == interaction.Member {
		name = interaction.User.Username
	} else {
		if nil != interaction.Member.User {
			name = interaction.Member.User.Username
		}
	}
	cmdData := interaction.ApplicationCommandData()
	if cmdData.Name == "jersey" {
		//Send ACK To meet the 3s turnaround and allow for more time to upload the image
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{}})
		go func() {
			//logger.Infof("Jersey Command called by: %s", name)
			//logger.Infof("Comamnd: %s", spew.Sdump(cmdData))
			subCmd := cmdData.Options[0]
			collection := collectionOpt(subCmd.Options[0].UintValue())
			team := subCmd.Options[1].StringValue()
			houndID := subCmd.Options[2].UintValue()

			if collection != houndsOpt {
				err := fmt.Errorf("Unsupported Collection: %s", cmdData.Options[0].Name)
				logger.Errorf("Error: %s", err)
				sendErrorResponse(houndID, err, session, interaction)
				return
			}

			//TODO: Refactor Getting Metadata and Image into its own method
			logger.Debugf("Getting metadata for Hound #%d", houndID)
			metadata, err := houndMetadataFetcher.Fetch(houndID)
			if err != nil {
				err := fmt.Errorf("Failed to fetch Metadata for Hound ID %d: %w", houndID, err)
				logger.Errorf("Error: %s", err)
				sendErrorResponse(houndID, err, session, interaction)
				return
			}

			hound, err := ipfsClient.GetImageFromIPFS(metadata.Image)
			if err != nil {
				err := fmt.Errorf("Failed to fetch image from IPFS for Hound ID %d: %w", houndID, err)
				logger.Errorf("Error: %s", err)
				sendErrorResponse(houndID, err, session, interaction)
				return
			}

			logger.Debugf("Overlaying Image")
			buff, err := stamper.OverlayHoundJersey(hound, metadata, team)
			if err != nil {
				err := fmt.Errorf("Failed to overlay %s %s jersey for Hound ID %d: %w", team, subCmd.Name, houndID, err)
				logger.Errorf("Error: %s", err)
				sendErrorResponse(houndID, err, session, interaction)
				return
			}

			file := &discordgo.File{
				Name:        fmt.Sprintf("%s_%s_%s_jersey_%d.png", name, team, strings.ReplaceAll(subCmd.Name, "-", "_"), houndID),
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

func sendErrorResponse(id uint64, err error, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	errMsg := err.Error()

	if strings.Contains(errMsg, "invalid character 'T' looking for beginning of value") {
		houndID := id
		errMsg = fmt.Sprintf("Error: Hound #%d has not yet been revealed", houndID)
	}

	if strings.Contains(errMsg, "invalid JPEG format") {
		houndID := id
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
