package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/avislash/sentamper/image"
	"github.com/avislash/sentamper/ipfs"
	"github.com/avislash/sentamper/metadata"
	"github.com/bwmarrin/discordgo"
)

var ipfsClient *ipfs.Client
var stamper *image.Processor

func main() {

	var err error
	ipfsClient, err = ipfs.NewClient()
	if err != nil {
		panic("Error creating IPFS Client: " + err.Error())
	}

	stamper = image.NewProcessor()

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
