package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"np-notify/pkg/models"
	"np-notify/pkg/utils"
	"strings"
)

type Fetcher struct {
	ApiURL string
}

func DataGetter(method string, url string, body string, token string) ([]interface{}, error) {
	req, err := http.NewRequest(
		method,
		url,
		io.NopCloser(strings.NewReader(body)),
	)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("auth=%s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var response []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	return response, nil
}

// Helper functions at package level
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getStringSliceFromMap(m map[string]interface{}, key string) []string {
	if val, ok := m[key].([]string); ok {
		return val
	}
	return []string{}
}

func (f *Fetcher) GetUnreadCount(ac models.Account, gameID string) (models.UnreadCount, error) {
	url := f.ApiURL + "/game_api/fetch_unread_count"
	body := "type=fetch_unread_count&version=&gameId=" + gameID
	var unreadCount models.UnreadCount

	response, err := DataGetter("POST", url, body, ac.Token)
	if err != nil {
		return unreadCount, err
	}

	if len(response) != 2 {
		return unreadCount, fmt.Errorf("unexpected response format")
	}

	unreadCountData, ok := response[1].(map[string]interface{})
	if !ok {
		return unreadCount, fmt.Errorf("unexpected response format")
	}

	unreadCount, err = utils.Unmarshalinterface[models.UnreadCount](unreadCountData)

	return unreadCount, nil
}

// Main function update
func (f *Fetcher) GetPlayerGames(ac models.Account) ([]models.PlayerGame, error) {
	url := fmt.Sprintf("%s/account_api/init_player", f.ApiURL)
	body := "type=init_player"

	response, err := DataGetter("POST", url, body, ac.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get player games: %v", err)
	}

	if len(response) != 2 {
		return nil, fmt.Errorf("unexpected response format")
	}

	playerData, ok := response[1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	playerGames, err := utils.Unmarshalinterface[[]models.PlayerGame](playerData["open_games"].([]interface{}))

	return playerGames, nil
}

func (f *Fetcher) GetGameEvents(ac models.Account, pg models.PlayerGame) (models.MessageWrapper, error) {
	mw := models.MessageWrapper{}

	messageCouint := "10"

	url := f.ApiURL + "/game_api/fetch_game_messages"
	body := "type=fetch_game_messages&count=" + messageCouint + "&offset=0&group=game_event&version=&gameId=" + pg.ID

	response, err := DataGetter("POST", url, body, ac.Token)
	if err != nil {
		return mw, fmt.Errorf("failed to get game events: %v", err)
	}

	if len(response) != 2 {
		return mw, fmt.Errorf("unexpected response format")
	}

	mw.ParseMessages(response)

	return mw, nil
}
func (f *Fetcher) GetGameDiplomacy(ac models.Account, pg models.PlayerGame) (models.MessageWrapper, error) {
	mw := models.MessageWrapper{}

	messageCount := "10"

	url := f.ApiURL + "/game_api/fetch_game_messages"
	body := "type=fetch_game_messages&count=" + messageCount + "&offset=0&group=game_diplomacy&version=&gameId=" + pg.ID

	response, err := DataGetter("POST", url, body, ac.Token)
	if err != nil {
		return mw, fmt.Errorf("failed to get game messages: %v", err)
	}

	if len(response) != 2 {
		return mw, fmt.Errorf("unexpected response format")
	}

	mw.ParseMessages(response)

	return mw, nil
}

func (f *Fetcher) GetGameMessageComments(ac models.Account, pg models.PlayerGame, key string) (*models.MessageCommentWrapper, error) {
	mcw := models.MessageCommentWrapper{}

	commentCount := "20"

	url := f.ApiURL + "/game_api/fetch_game_message_comments"
	body := "type=fetch_game_message_comments&key=" + key + "&count=" + commentCount + "&offset=0&version=&gameId=" + pg.ID

	response, err := DataGetter("POST", url, body, ac.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get game message comments: %v", err)
	}

	if len(response) != 2 {
		return nil, fmt.Errorf("unexpected response format")
	}

	mcw.Parse(response)

	return &mcw, nil
}
