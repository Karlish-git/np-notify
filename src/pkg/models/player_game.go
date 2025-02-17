package models

import "fmt"

type PlayerGame struct {
	ID            string                 `json:"id"`
	Number        string                 `json:"number"`
	Creator       string                 `json:"creator"`
	Created       string                 `json:"created"`
	Finished      string                 `json:"finished"`
	Kind          string                 `json:"kind"`
	Status        string                 `json:"status"`
	JoinedPlayers []string               `json:"joinedPlayers"`
	Config        map[string]interface{} `json:"config"`
	Unreads       UnreadCount
	Events        MessageWrapper
	Diplomacy     MessageWrapper
}

func (pg PlayerGame) IsActive() bool {
	return pg.Status == "active"
}

func (pg PlayerGame) GetGameName() string {
	return pg.Config["name"].(string)
}

func (pg PlayerGame) GetPlayerById(id int) string {
	if id > 0 && id <= len(pg.JoinedPlayers) {
		return pg.JoinedPlayers[id-1]
	}
	return "Unknown Player Name"
}

func (pg PlayerGame) HtmlString() string {
	var html string

	html += fmt.Sprintf("%s - %s \n", pg.GetGameName(), pg.Number)

	html += fmt.Sprintf("- Global chat: %s\n", pg.Unreads.GlobalChat)
	html += fmt.Sprintf("- Diplomacy(unread threads): %s \n", pg.Unreads.Diplomacy)
	html += fmt.Sprintf("- New Events: %s \n", pg.Unreads.Events)

	new_events := pg.Events.GetUnreadMessages()

	html += "Dimpolacy\n"

	new_messages := pg.Diplomacy.GetUnreadMessages()
	if len(new_messages) > 0 {
		html += ""
		for _, message := range new_messages {
			html += fmt.Sprintf("at %s %s\n", message.Activity.Format("15:04"), message.Payload.Print(pg))
		}
		html += "\n"
	}
	html += "\n"

	html += "Events\n"
	if len(new_events) > 0 {
		html += ""
		for idx, event := range new_events {
			html += fmt.Sprintf("%d. at %s %s\n", idx, event.Created.Format("15:04"), event.Payload.Print(pg))
		}
		html += "\n"
	}

	return html
}
func (pg PlayerGame) HtmlStringDiplomacy() string {
	var html string

	return html
}

// func (pg PlayerGame) HtmlString() string {
// 	var html string
//
// 	html += fmt.Sprintf("<h2>%s - %s</h2> \n", pg.GetGameName(), pg.Number)
//
// 	html += fmt.Sprintf("- Global chat: %s</br>\n", pg.Unreads.GlobalChat)
// 	html += fmt.Sprintf("- Diplomacy(unread threads): %s \n</br>", pg.Unreads.Diplomacy)
// 	html += fmt.Sprintf("- New Events: %s \n</br>", pg.Unreads.Events)
//
// 	new_events := pg.Events.GetUnreadMessages()
//
// 	html += "<h3>Dimpolacy</h3>\n"
//
// 	new_messages := pg.Diplomacy.GetUnreadMessages()
// 	if len(new_messages) > 0 {
// 		html += "<pre>"
// 		for _, message := range new_messages {
// 			html += fmt.Sprintf("at %s %s\n", message.Created.Format("15:04"), message.Payload.Print(pg))
// 		}
// 		html += "</pre>\n"
// 	}
// 	html += "\n"
//
// 	html += "<h3>Events</h3>\n"
// 	if len(new_events) > 0 {
// 		html += "<pre>"
// 		for idx, event := range new_events {
// 			html += fmt.Sprintf("%d. at %s %s\n", idx, event.Created.Format("15:04"), event.Payload.Print(pg))
// 		}
// 		html += "</pre>\n"
// 	}
//
// 	return html
// }
