package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	prefixCmd string = "uc!"

	// Commands
	cmdPlay  = prefixCmd + "play"
	cmdJoin  = prefixCmd + "join"
	cmdStart = prefixCmd + "start"
	cmdVote  = prefixCmd + "vote"
	cmdKick  = prefixCmd + "kick"
)

var (
	game  Game
	botID string
)

func CheckError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	go startServer()

	discord, err := discordgo.New("Bot " + os.Args[1])
	CheckError("Error creating discord session", err)
	user, err := discord.User("@me")
	CheckError("Error retrieving account", err)

	botID = user.ID
	discord.AddHandler(CommandHandler)
	discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Undercover^^")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
	})

	err = discord.Open()
	CheckError("Error opening connection to Discord", err)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = discord.Close()
}

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
		return
	}

	// If a command is called
	if strings.HasPrefix(m.Content, prefixCmd) {
		args := strings.Split(m.Content, " ")
		// Check game state
		switch game.state {
		// If the game is off
		case Off:
			switch args[0] {
			// If the play command is called
			case cmdPlay:
				// Get the channel where the game is being played
				game.channel, _ = s.State.GuildChannel(m.GuildID, m.ChannelID)
				game.session = s
				// Set game state to "waiting for players"
				game.SetState(Waiting)
				msg := game.AddPlayer(m.Author)

				// TODO: Remove, this is for easy tests
				for i := 0; i < 3; i++ {
					dummy := discordgo.User{
						ID: fmt.Sprintf("dummy%d", i),
					}
					msg += game.AddPlayer(&dummy)
				}

				// Sends message to channel
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s started an **Undercover** game!\n%d players are required to start, type `%s` in %s to join.\nOnly %s can type `%s` to start the game.\n%s",
					m.Author.Mention(),
					playerMin,
					cmdJoin,
					game.channel.Mention(),
					game.GetCreator().Mention(),
					cmdStart,
					msg))
			// If any other command is called, send error message
			case cmdJoin:
				fallthrough
			case cmdStart:
				fallthrough
			case cmdVote:
				fallthrough
			case cmdKick:
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No game is currently running. Start one by typing `%s` in any channel", cmdPlay))
			}
		// Is the game is waiting for players
		case Waiting:
			switch args[0] {
			// Using play command while another game is started results in an error
			case cmdPlay:
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A game is currently waiting for players. Join it by typing `%s` in %s", cmdJoin, game.channel.Mention()))
			// Join a game currently running
			case cmdJoin:
				if game.IsOnSameChannel(s, m) {
					_, _ = s.ChannelMessageSend(m.ChannelID, game.AddPlayer(m.Author))
				} else {
					_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A game is waiting for players in %s, you can join there using `%s`", game.channel.Mention(), cmdJoin))
				}
			case cmdStart:
				if game.IsOnSameChannel(s, m) {
					if m.Author.ID == game.GetCreator().ID {
						if len(args) != 3 {
							_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("To start the game, use `%s <number of %s> <number of %s>`", cmdStart, Undercover.String(), MrWhite.String()))
						} else if game.IsReady() {
							undercoverNumber, err1 := strconv.ParseInt(args[1], 10, 32)
							mrWhiteNumber, err2 := strconv.ParseInt(args[2], 10, 32)
							if err1 != nil || err2 != nil || undercoverNumber < 0 || mrWhiteNumber < 0 {
								_, _ = s.ChannelMessageSend(m.ChannelID, "The arguments must be whole positive numbers")
							} else {
								game.Start(undercoverNumber, mrWhiteNumber)
							}
						} else {
							_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d players are required to start a game, %d have joined", playerMin, len(game.players)))
						}
					} else {
						_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Only %s can start the game", game.GetCreator().Mention()))
					}
				} else {
					_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A game is waiting for players in %s, you can join there using `%s`", game.channel.Mention(), cmdJoin))
				}
			case cmdVote:
				fallthrough
			case cmdKick:
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A game is currently waiting for players. Join it by typing `%s` in %s", cmdJoin, game.channel.Mention()))
			}
		case Running:
			switch args[0] {
			case cmdPlay:
				fallthrough
			case cmdJoin:
				fallthrough
			case cmdStart:
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A game is currently running in %s. You must wait for it to end before starting a new one with `%s`", game.channel.Mention(), cmdPlay))
			case cmdVote:
				_, _ = s.ChannelMessageSend(m.ChannelID, "Votes have not been implemented yet!")
			case cmdKick:
				_, _ = s.ChannelMessageSend(m.ChannelID, "Kicking players has not been implemented yet!")
			}
		}
	}
}
