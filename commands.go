package main

import (
	"fmt"
	"github.com/pyed/transmission"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
)

// receiveTorrent gets an update that potentially has a .torrent file to add
func receiveTorrent(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	if ud.Message.Document.FileID == "" {
		return // has no document
	}

	// get the file ID and make the config
	fconfig := tgbotapi.FileConfig{
		FileID: ud.Message.Document.FileID,
	}
	file, err := bot.GetFile(fconfig)
	if err != nil {
		send(bot, "*ERROR*: "+err.Error(), ud.Message.Chat.ID)
		return
	}

	// add by file URL
	addTorrentsByURL(bot, client, ud, []string{file.Link(bot.Token)})
}

// stop takes id[s] of torrent[s] or 'all' to stop them
func stop(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	// make sure that we got at least one argument
	if len(ud.Tokens()) == 0 {
		send(bot, "*stop*: needs an argument", ud.Message.Chat.ID)
		return
	}

	// if the first argument is 'all' then stop all torrents
	if ud.Tokens()[0] == "all" {
		if err := client.StopAll(); err != nil {
			send(bot, "*stop*: error occurred while stopping some torrents", ud.Message.Chat.ID)
			return
		}
		send(bot, "*stop*: all torrents stopped", ud.Message.Chat.ID)
		return
	}

	for _, id := range ud.Tokens() {
		num, err := strconv.Atoi(id)
		if err != nil {
			send(bot, fmt.Sprintf("*stop*: `%s` is not a number", id), ud.Message.Chat.ID)
			continue
		}
		status, err := client.StopTorrent(num)
		if err != nil {
			send(bot, "*stop*: "+err.Error(), ud.Message.Chat.ID)
			continue
		}

		torrent, err := client.GetTorrent(num)
		if err != nil {
			send(bot, fmt.Sprintf("*[fail] stop*: No torrent with an ID of %d", num), ud.Message.Chat.ID)
			return
		}
		send(bot, fmt.Sprintf("*[%s] stop*: `%s`", status, torrent.Name), ud.Message.Chat.ID)
	}
}

// start takes id[s] of torrent[s] or 'all' to start them
func start(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	// make sure that we got at least one argument
	if len(ud.Tokens()) == 0 {
		send(bot, "*start*: needs an argument", ud.Message.Chat.ID)
		return
	}

	// if the first argument is 'all' then start all torrents
	if ud.Tokens()[0] == "all" {
		if err := client.StartAll(); err != nil {
			send(bot, "*start*: error occurred while starting some torrents", ud.Message.Chat.ID)
			return
		}
		send(bot, "*start*: all torrents started", ud.Message.Chat.ID)
		return

	}

	for _, id := range ud.Tokens() {
		num, err := strconv.Atoi(id)
		if err != nil {
			send(bot, fmt.Sprintf("*start*: `%s` is not a number", id), ud.Message.Chat.ID)
			continue
		}
		status, err := client.StartTorrent(num)
		if err != nil {
			send(bot, "*start*: "+err.Error(), ud.Message.Chat.ID)
			continue
		}

		torrent, err := client.GetTorrent(num)
		if err != nil {
			send(bot, fmt.Sprintf("*[fail] start*: No torrent with an ID of %d", num), ud.Message.Chat.ID)
			return
		}
		send(bot, fmt.Sprintf("*[%s] start: `%s`", status, torrent.Name), ud.Message.Chat.ID)
	}
}

// check takes id[s] of torrent[s] or 'all' to verify them
func check(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	// make sure that we got at least one argument
	if len(ud.Tokens()) == 0 {
		send(bot, "*check*: needs an argument", ud.Message.Chat.ID)
		return
	}

	// if the first argument is 'all' then start all torrents
	if ud.Tokens()[0] == "all" {
		if err := client.VerifyAll(); err != nil {
			send(bot, "*check*: error occurred while verifying some torrents", ud.Message.Chat.ID)
			return
		}
		send(bot, "*check*: verifying all torrents", ud.Message.Chat.ID)
		return

	}

	for _, id := range ud.Tokens() {
		num, err := strconv.Atoi(id)
		if err != nil {
			send(bot, fmt.Sprintf("*check*: `%s` is not a number", id), ud.Message.Chat.ID)
			continue
		}
		status, err := client.VerifyTorrent(num)
		if err != nil {
			send(bot, "*stop*: "+err.Error(), ud.Message.Chat.ID)
			continue
		}

		torrent, err := client.GetTorrent(num)
		if err != nil {
			send(bot, fmt.Sprintf("*[fail] check*: No torrent with an ID of %d", num), ud.Message.Chat.ID)
			return
		}
		send(bot, fmt.Sprintf("*[%s] check*: `%s`", status, torrent.Name), ud.Message.Chat.ID)
	}

}

// del takes an id or more, and delete the corresponding torrent/s
func del(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	// make sure that we got an argument
	if len(ud.Tokens()) == 0 {
		send(bot, "*del*: needs an ID", ud.Message.Chat.ID)
		return
	}

	// loop over ud.Tokens() to read each potential id
	for _, id := range ud.Tokens() {
		num, err := strconv.Atoi(id)
		if err != nil {
			send(bot, fmt.Sprintf("*del*: `%s` is not an ID", id), ud.Message.Chat.ID)
			return
		}

		name, err := client.DeleteTorrent(num, false)
		if err != nil {
			send(bot, "*del*: "+err.Error(), ud.Message.Chat.ID)
			return
		}

		send(bot, fmt.Sprintf("*del*: `%s`", name), ud.Message.Chat.ID)
	}
}

// deldata takes an id or more, and delete the corresponding torrent/s with their data
func deldata(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	// make sure that we got an argument
	if len(ud.Tokens()) == 0 {
		send(bot, "*deldata*: needs an ID", ud.Message.Chat.ID)
		return
	}
	// loop over ud.Tokens() to read each potential id
	for _, id := range ud.Tokens() {
		num, err := strconv.Atoi(id)
		if err != nil {
			send(bot, fmt.Sprintf("*deldata*: `%s` is not an ID", id), ud.Message.Chat.ID)
			return
		}

		name, err := client.DeleteTorrent(num, true)
		if err != nil {
			send(bot, "*deldata*: "+err.Error(), ud.Message.Chat.ID)
			return
		}

		send(bot, fmt.Sprintf("*deldata*: Deleted with data: `%s`", name), ud.Message.Chat.ID)
	}
}

// version sends transmission version + transmission-telegram version
func version(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	send(bot, fmt.Sprintf("Transmission *%s*\nTransmission-telegram *%s*", client.Version(), VERSION), ud.Message.Chat.ID)
}

// addTorrentsByURL adds torrent files or magnet links passed by rls
func addTorrentsByURL(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper, urls []string) {
	if len(urls) == 0 {
		send(bot, "*add*: needs atleast one URL", ud.Message.Chat.ID)
		return
	}

	// loop over the URL/s and add them
	for _, url := range urls {
		cmd := transmission.NewAddCmdByURL(url)

		torrent, err := client.ExecuteAddCommand(cmd)
		if err != nil {
			send(bot, "*add*: "+err.Error(), ud.Message.Chat.ID)
			continue
		}

		// check if torrent.Name is empty, then an error happened
		if torrent.Name == "" {
			send(bot, "*add*: error adding "+url, ud.Message.Chat.ID)
			continue
		}
		send(bot, fmt.Sprintf("*add*: *%d* `%s`", torrent.ID, torrent.Name), ud.Message.Chat.ID)
	}
}

// add takes an URL to a .torrent file in message to add it to transmission
func add(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	addTorrentsByURL(bot, client, ud, ud.Tokens())
}

// help sends help messsage
func help(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	send(bot, HELP, ud.Message.Chat.ID)
}

// unknownCommand sends message that command is unknown
func unknownCommand(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	send(bot, "no such command, try /help", ud.Message.Chat.ID)
}

// sort changes torrents sorting
func sort(bot *tgbotapi.BotAPI, client *transmission.TransmissionClient, ud UpdateWrapper) {
	if len(ud.Tokens()) == 0 {
		send(bot, `sort takes one of:
			(*id, name, age, size, progress, downspeed, upspeed, download, upload, ratio*)
			optionally start with (*rev*) for reversed order
			e.g. "*sort rev size*" to get biggest torrents first.`, ud.Message.Chat.ID)
		return
	}

	var reversed bool
	tokens := ud.Tokens()
	if strings.ToLower(ud.Tokens()[0]) == "rev" {
		reversed = true
		tokens = ud.Tokens()[1:]
	}

	mode := SortMethod{strings.ToLower(tokens[0]), reversed}

	if mode, ok := SortingMethods[mode]; ok {
		client.SetSort(mode)
		send(bot, fmt.Sprintf("*sort*: `%s` reversed: %t", tokens[0], reversed), ud.Message.Chat.ID)
	} else {
		send(bot, "*sort*: unkown sorting method", ud.Message.Chat.ID)
	}
}
