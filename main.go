package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/avislash/sentamper/config"
	"github.com/avislash/sentamper/image"
	"github.com/avislash/sentamper/ipfs"
	"github.com/avislash/sentamper/metadata"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var ipfsClient *ipfs.Client
var stamper *image.Processor
var metadataFetcher *metadata.SentinelMetadataFetcher

func init() {
	var configParams config.Config
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to read in config.yaml: %s", err))
	}

	err = yaml.Unmarshal(configFile, &configParams)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config.yaml: %s", err))
	}

	stamper, err = image.NewProcessor(configParams.ImageProcessorConfig)
	if err != nil {
		panic("Error initializing Image Processor: " + err.Error())
	}

	ipfsClient, err = ipfs.NewClient()
	if err != nil {
		panic("Error creating IPFS Client: " + err.Error())
	}
}

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(gmInteraction)
	if err := dg.Open(); err != nil {
		panic(err)
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
		panic(err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	log.Println("Bot started")
	<-ctx.Done()
	log.Println("Exit")
	_ = sentinelID
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
			sentinelID := cmdData.Options[0].UintValue()

			fetcher := metadata.NewSentinelMetadataFetcher("https://api.appliedprimate.dev/sentinels/metadata")
			metadata, err := fetcher.FetchMetdata(sentinelID)
			if err != nil {
				err := fmt.Errorf("Failed to retrieve metadata for Sentienl #%d: %w", sentinelID, err)
				log.Println("Error: ", err)
				sendErrorResponse(err, session, interaction)
				return
			}

			sentinel, err := ipfsClient.GetSentinelFromIPFS(metadata.Image)
			if err != nil {
				err := fmt.Errorf("Failed to retrieve Sentinel #%d image from IPFS: %w", sentinelID, err)
				log.Println("Error: ", err)
				sendErrorResponse(err, session, interaction)
				return
			}

			buff, err := stamper.OverlayMug(sentinel, metadata.BaseArmor) //combineImages(metadata)
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
