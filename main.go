package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/ipfs/boxo/files"
	ipfsClient "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	//fmt.Println(token)
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

	_, err = dg.ApplicationCommandCreate(botID, "", &discordgo.ApplicationCommand{
		Name:        "gm",
		Description: "Responds with a GM",
	})

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	<-ctx.Done()
	log.Println("Exit")
}

func gmInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var mention string
	if nil == interaction.Member {
		mention = interaction.User.Mention()
	} else {
		mention = interaction.Member.Mention()
	}
	if discordgo.InteractionApplicationCommand == interaction.Type {
		if interaction.ApplicationCommandData().Name == "gm" {
			var response *discordgo.InteractionResponse
			buff, err := combineImages()

			if err == nil {
				file := &discordgo.File{
					Name:        "gm.png",
					ContentType: "image/png",
					Reader:      buff,
				}
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "GM " + mention,
						Files:   []*discordgo.File{file},
					},
				}
			} else {
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to create Image: " + err.Error(),
					},
				}
			}

			if err := session.InteractionRespond(interaction.Interaction, response); err != nil {
				log.Println("Error sending message: ", err)
			}
		}
	}
}

func getSentinelFromIPFS() (image.Image, error) {
	// Create a new IPFS client
	client, err := ipfsClient.NewLocalApi()
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %w", err)
	}

	// Sentinel CID
	//cid := path.New("QmVttt4xLfRkGXAgMDqvLXXgse4GWTczHXBWpYJrSgyZeu") //trippy sentinel
	cid := path.New("QmY7xvucdb7DqSRvWEpotPPM3yUKijMgN5BMTWtLeKyqFG")
	// Retrieve the file from IPFS
	node, err := client.Unixfs().Get(context.Background(), cid)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving centinel from IPFS Hash %s: %w", cid, err)
	}

	file := files.ToFile((node))
	defer file.Close()

	sentinel, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Error decoding IPFS File as PNG: %w", err)
	}

	return sentinel, nil
}

func combineImages() (*bytes.Buffer, error) {
	sentinel, err := getSentinelFromIPFS()
	if err != nil {
		return nil, err
	}

	// Open the second image file
	sentinelHandFile, err := os.Open("./mugs/trippyMug.png")
	if err != nil {
		panic(err)
	}
	defer sentinelHandFile.Close()

	sentinelHand, err := png.Decode(sentinelHandFile)
	if err != nil {
		panic(err)
	}

	// Create a new image with the size of the larger image
	combinedWidth := max(sentinel.Bounds().Max.X, sentinelHand.Bounds().Max.X)
	combinedHeight := max(sentinel.Bounds().Max.Y, sentinelHand.Bounds().Max.Y)
	combinedImg := image.NewRGBA(image.Rect(0, 0, combinedWidth, combinedHeight))

	// Draw the first image onto the combined image
	draw.Draw(combinedImg, sentinel.Bounds(), sentinel, image.ZP, draw.Src)

	// Draw the second image onto the combined image with an offset
	offset := image.Pt((combinedWidth-sentinelHand.Bounds().Dx())/2, (combinedHeight-sentinelHand.Bounds().Dy())/2)
	drawRect := sentinelHand.Bounds()
	drawRect = drawRect.Add(offset)
	drawRect = drawRect.Intersect(combinedImg.Bounds())
	drawRect = drawRect.Sub(offset)
	drawRect = drawRect.Add(offset)
	drawRect = drawRect.Intersect(sentinelHand.Bounds())
	drawRect = drawRect.Sub(offset)
	draw.Draw(combinedImg, drawRect, sentinelHand, sentinelHand.Bounds().Min, draw.Over)

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, combinedImg); err != nil {
		return nil, fmt.Errorf("Error Encoding Image: %w", err)
	}

	return buff, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
