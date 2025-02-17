package models

import (
	"fmt"
	"np-notify/pkg/utils"
	"strings"
	"time"
)

type MessageWrapper struct {
	EventType string
	Data      struct {
		Messages []Message `json:"messages"`
		Group    string    `json:"group"`
	}
}

func GetUnreadMessagesFromWrappers(wrappers []MessageWrapper) []*Message {
	var messages []*Message
	for _, wrapper := range wrappers {
		for _, message := range wrapper.Data.Messages {
			if message.Status == "unread" {
				messages = append(messages, &message)
			}
		}
	}
	return messages
}

func (m MessageWrapper) GetUnreadMessages() []*Message {
	var messages []*Message
	for _, message := range m.Data.Messages {
		if message.Status == "unread" {
			messages = append(messages, &message)
		}
	}
	return messages
}

type Message struct {
	Status   string    `json:"status"`
	Key      string    `json:"key"`
	Created  time.Time `json:"created"`
	Activity time.Time `json:"activity"`
	Group    string    `json:"group"`
	Payload  Payload   `json:"payload"`
}

// Payload is an interface for different message payload types
type Payload interface {
	Print(game PlayerGame) string
}

// DiplomaticMessagePayload represents diplomatic message events
type DiplomaticMessagePayload struct {
	Template  string `json:"template"`
	FromUID   int    `json:"from_uid"`
	FromColor string `json:"from_color"`
	ToAliases string `json:"to_aliases"`
	ToColors  string `json:"to_colors"`
	ToUIDs    []int  `json:"to_uids"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (dmp DiplomaticMessagePayload) Print(game PlayerGame) string {
	fromName := game.GetPlayerById(dmp.FromUID)
	toAliases := strings.ReplaceAll(dmp.ToAliases, ",", ", ")
	return fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n", fromName, toAliases, dmp.Subject)
}

// func (dmp *DiplomaticMessagePayload) GetToAliases() []string {
// 	return strings.Split(dmp.ToAliases, ",")
// }

// MoneySentPayload represents money transfer events
type MoneySentPayload struct {
	Template string `json:"template"`
	ToPUID   int    `json:"to_puid"`
	FromPUID int    `json:"from_puid"`
	Amount   int    `json:"amount"`
	Tick     int    `json:"tick"`
}

func (msp MoneySentPayload) Print(game PlayerGame) string {
	fromName := game.GetPlayerById(msp.FromPUID)
	toName := game.GetPlayerById(msp.ToPUID)
	return fmt.Sprintf("From: %s\nTo: %s\nAmount: %d\n", fromName, toName, msp.Amount)
}

// WarDeclaredPayload represents war declaration events
type WarDeclaredPayload struct {
	Template string `json:"template"`
	Attacker int    `json:"attacker"`
	Defender int    `json:"defender"`
	Tick     int    `json:"tick"`
}

func (wdp WarDeclaredPayload) Print(game PlayerGame) string {
	attackerName := game.GetPlayerById(wdp.Attacker)
	defenderName := game.GetPlayerById(wdp.Defender)
	return fmt.Sprintf("War declared by %s against %s\n", attackerName, defenderName)
}

// GoodbyePayload represents player leaving events
type GoodbyePayload struct {
	Template string `json:"template"`
	UID      int    `json:"uid"`
	Name     string `json:"name"`
	Tick     int    `json:"tick"`
}

func (gp GoodbyePayload) Print(game PlayerGame) string {
	playerName := game.GetPlayerById(gp.UID)
	return fmt.Sprintf("Goodbye from %s\n", playerName)
}

// SharedTechnologyPayload represents technology sharing events
type SharedTechnologyPayload struct {
	Template string `json:"template"`
	ToPUID   int    `json:"to_puid"`
	FromPUID int    `json:"from_puid"`
	Tech     int    `json:"tech"`
	Level    int    `json:"level"`
	Price    int    `json:"price"`
	Tick     int    `json:"tick"`
}

func (stp SharedTechnologyPayload) Print(game PlayerGame) string {
	fromName := game.GetPlayerById(stp.FromPUID)
	toName := game.GetPlayerById(stp.ToPUID)
	return fmt.Sprintf(
		"Technology shared from %s to %s\nTech: %s\nLevel: %d\n",
		fromName, toName, utils.GetTechNameById(stp.Tech), stp.Level,
	)
}

// CombatPayload represents combat events
type CombatPayload struct {
	Template  string           `json:"template"`
	Star      Star             `json:"star"`
	Attackers map[string]Fleet `json:"attackers"`
	Defenders map[string]Fleet `json:"defenders"`
	Dw        int              `json:"dw"`
	Aw        int              `json:"aw"`
	Loot      int              `json:"loot"`
	Looter    int              `json:"looter"`
	Tick      int              `json:"tick"`
}

func (cp CombatPayload) Print(game PlayerGame) string {
	return fmt.Sprintf(
		"Combat at %s's %s. Defender ships left: %d, Attacker ships left: %d\n",
		game.GetPlayerById(cp.Star.PUID), cp.Star.Name,
		sumShips(cp.Defenders)+cp.Star.EndingShips, sumShips(cp.Attackers),
	)
}

func sumShips(fleets map[string]Fleet) int {
	total := 0
	for _, fleet := range fleets {
		total += fleet.EndingShips
	}
	return total
}

// Star represents a star in combat
type Star struct {
	Name          string `json:"name"`
	UID           int    `json:"uid"`
	StartingShips int    `json:"ss"`
	EndingShips   int    `json:"es"`
	PUID          int    `json:"puid"`
	W             int    `json:"w"`
}

// Fleet represents a fleet in combat
type Fleet struct {
	StartingShips int `json:"ss"`
	EndingShips   int `json:"es"`
	PUID          int `json:"puid"`
	W             int `json:"w"`
}

// PeaceAcceptedPayload represents peace acceptance events
type PeaceAcceptedPayload struct {
	Template string `json:"template"`
	ToPUID   int    `json:"to_puid"`
	FromPUID int    `json:"from_puid"`
	Tick     int    `json:"tick"`
}

func (pap PeaceAcceptedPayload) Print(game PlayerGame) string {
	fromName := game.GetPlayerById(pap.FromPUID)
	toName := game.GetPlayerById(pap.ToPUID)
	return fmt.Sprintf("Peace accepted between %s and %s\n", fromName, toName)
}

// PeaceRequestedPayload represents peace request events
type PeaceRequestedPayload struct {
	Template string `json:"template"`
	ToPUID   int    `json:"to_puid"`
	FromPUID int    `json:"from_puid"`
	Price    int    `json:"price"`
	Tick     int    `json:"tick"`
}

func (prp PeaceRequestedPayload) Print(game PlayerGame) string {
	fromName := game.GetPlayerById(prp.FromPUID)
	toName := game.GetPlayerById(prp.ToPUID)
	return fmt.Sprintf("Peace requested by %s to %s\n", fromName, toName)
}

// ProductionPayload represents production events
type ProductionPayload struct {
	Template   string `json:"template"`
	UID        int    `json:"uid"`
	Economy    int    `json:"economy"`
	Banking    int    `json:"banking"`
	Cash       int    `json:"cash"`
	TechName   int    `json:"techName"`
	TechPoints int    `json:"techPoints"`
	TechLevel  int    `json:"techLevel"`
	Sfv        int    `json:"sfv"`
	Tick       int    `json:"tick"`
}

func (pp ProductionPayload) Print(game PlayerGame) string {
	playerName := game.GetPlayerById(pp.UID)
	return fmt.Sprintf(
		"Production report for %s\nCash: %d\nTech: %s\nTech Points: %d\n",
		playerName, pp.Cash, utils.GetTechNameById(pp.TechName), pp.TechPoints)
}

// TechUpPayload represents technology upgrade events
type TechUpPayload struct {
	Template string `json:"template"`
	UID      int    `json:"uid"`
	Tech     int    `json:"tech"`
	Level    int    `json:"level"`
	Tick     int    `json:"tick"`
}

func (tup TechUpPayload) Print(game PlayerGame) string {
	playerName := game.GetPlayerById(tup.UID)
	return fmt.Sprintf(
		"Technology %s %d gained by %s\n",
		utils.GetTechNameById(tup.Tech), tup.Level, playerName)
}

func (m *MessageWrapper) ParseMessages(c []interface{}) {
	m.EventType = c[0].(string)

	data := c[1].(map[string]interface{})
	m.Data.Group = data["group"].(string)

	messages := []Message{}

	unparsedMsgs := data["messages"].([]interface{})

	for _, msg := range unparsedMsgs {
		msg := msg.(map[string]interface{})
		message := Message{
			Status:   msg["status"].(string),
			Key:      msg["key"].(string),
			Created:  utils.StringToTime(msg["created"].(string)),
			Activity: utils.StringToTime(msg["activity"].(string)),
			Group:    msg["group"].(string),
		}

		unparsed := msg["payload"].(map[string]interface{})

		if m.Data.Group == "game_event" {
			switch unparsed["template"].(string) {
			case "money_sent":
				payload, err := utils.Unmarshalinterface[MoneySentPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal money sent payload: %v", err)
					continue
				}
				message.Payload = payload
			case "war_declared":
				payload, err := utils.Unmarshalinterface[WarDeclaredPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal war declared payload: %v", err)
					continue
				}
				message.Payload = payload
			case "goodbye_to_player":
				fallthrough
			case "goodbye_to_player_inactivity":
				fallthrough
			case "goodbye_to_player_defeated":
				fallthrough
			case "goodbye":
				payload, err := utils.Unmarshalinterface[GoodbyePayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal goodbye payload: %v", err)
					continue
				}
				message.Payload = payload
			case "shared_technology":
				payload, err := utils.Unmarshalinterface[SharedTechnologyPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal shared technology payload: %v", err)
					continue
				}
				message.Payload = payload
				message.Payload = payload
			case "combat_mk_ii":
				payload, err := utils.Unmarshalinterface[CombatPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal combat payload: %v", err)
					continue
				}
				message.Payload = payload
			case "peace_requested":
				payload, err := utils.Unmarshalinterface[PeaceRequestedPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal peace requested payload: %v", err)
					continue
				}
				message.Payload = payload
			case "peace_accepted":
				payload, err := utils.Unmarshalinterface[PeaceAcceptedPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal peace accepted payload: %v", err)
					continue
				}
				message.Payload = payload
			case "production":
				payload, err := utils.Unmarshalinterface[ProductionPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal production payload: %v", err)
					continue
				}
				message.Payload = payload
			case "tech_up":
				payload, err := utils.Unmarshalinterface[TechUpPayload](unparsed)
				if err != nil {
					utils.LogError("Failed to unmarshal tech up payload: %v", err)
					continue
				}
				message.Payload = payload
			default:
				utils.LogError("Unhandled message type %s", unparsed["template"].(string))
				continue
			}
		} else if m.Data.Group == "game_diplomacy" {
			payload, err := utils.Unmarshalinterface[DiplomaticMessagePayload](unparsed)
			if err != nil {
				utils.LogError("Failed to unmarshal diplomatic message payload: %v", err)
				continue
			}
			message.Payload = payload
		} else {
			utils.LogError("Unhandled message group %s", m.Data.Group)
			continue
		}

		messages = append(messages, message)
	}

	m.Data.Messages = messages
}
