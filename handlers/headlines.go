package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/config"
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/gameon-app-inc/fanclash-event-processor/processor"
	"github.com/labstack/gommon/log"
	"github.com/palantir/stacktrace"
	"github.com/volatiletech/null/v8"
)

const (
	HeadlineTypeSlideIn        = "slide_in"
	HeadlineTypeFadeIn         = "fade_in"
	HeadlineTypeLastPickAction = "last_pick_action"
	HeadlineTypeCountdown      = "countdown"
	HeadlineTypePrize          = "prize"

	HeadlineImageTypeUser   = "user"
	HeadlineImageTypeTeam   = "team"
	HeadlineImageTypePlayer = "player"
)

type Headline struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Rotation    int      `json:"rotation"`
	Type        string   `json:"type"`
	ImageType   string   `json:"image_type"`
}

func NewHeadlinesHandler(store database.Store) *BaseEventHandler {
	return NewBaseEventHandler(store, []localHandler{
		headlinesHandler,
	}, false)
}

func headlinesHandler(_ database.Store, event *processor.Event) error {
	if event.Type == ActionLineUp || event.Type == ActionPeriodStart || event.Type == ActionPeriodEnd ||
		event.Type == ActionMatchEnd {

		GetDebouncedSendMatchHeadlines()(event.MatchID)
	}
	return nil
}

func SendAllHeadlines(store database.Store) error {
	matches, err := store.GetMatchesForHeadlines(time.Now())
	if err != nil {
		return stacktrace.Propagate(err, "cannot get matches for headlines")
	}

	for _, match := range matches {
		GetDebouncedSendMatchHeadlines()(match.ID)
	}

	return nil
}

func SendHeadlinesForMatch(store database.Store, matchID string) error {
	lobbyHeadlines, err := getLobbyHeadlines(store, matchID)
	if err != nil {
		return stacktrace.Propagate(err, "cannot get fixture headlines for match_id %s", matchID)
	}

	gameplayHeadlines, err := getGameplayHeadlines(store, matchID)
	if err != nil {
		return stacktrace.Propagate(err, "cannot get gameplay headlines for match_id %s", matchID)
	}

	fullTimeHeadlines, err := getFullTimeHeadlines(store, matchID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log the missing data issue and continue
			log.Warn("No data available for match_id:", matchID)
			return nil // Skip this match
		}
		return stacktrace.Propagate(err, "Failed to get full-time headlines for match_id %s", matchID)
	}

	var headlines schema.MatchHeadlineSlice
	headlines = append(headlines, toMatchHeadlines(matchID, 1, lobbyHeadlines)...)
	headlines = append(headlines, toMatchHeadlines(matchID, 2, gameplayHeadlines)...)
	if len(fullTimeHeadlines) > 0 {
		headlines = append(headlines, toMatchHeadlines(matchID, 3, fullTimeHeadlines)...)
	}

	// insert into db
	if err := store.UpdateHeadlines(matchID, headlines); err != nil {
		return stacktrace.Propagate(err, "cannot insert match headlines in db")
	}

	payload := map[string]interface{}{
		"match_id": matchID,
	}

	// marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		return stacktrace.Propagate(err, "cannot marshal payload")
	}

	_, err = store.InsertAMQPEvent(&schema.AmqpEvent{
		Exchange: config.RMQFCMExchange(),
		Type:     "headlines_updated",
		Data:     string(body),
	})
	if err != nil {
		return stacktrace.Propagate(err, "cannot insert amqp_event")
	}

	return nil
}

func getLobbyHeadlines(store database.Store, matchID string) ([]*Headline, error) {
	match, err := store.GetMatchByIDWithTeams(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get match by id")
	}

	result := make([]*Headline, 0)

	// 1) get fans count
	fanCount, err := store.GetGameCount(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get number of games")
	}

	if fanCount > 0 {
		avatars, err := store.GetRandomUserAvatars(matchID, 6)
		if err != nil {
			return nil, stacktrace.Propagate(err, "cannot get random user avatars")
		}

		result = append(result, &Headline{
			ID:          "10",
			Title:       fmt.Sprintf("%d Fans", fanCount),
			Description: "Have joined this game 🌎",
			Images:      avatars,
			Rotation:    1,
			Type:        HeadlineTypeFadeIn,
			ImageType:   HeadlineImageTypeUser,
		})
	}

	// 2) Countdown
	result = append(result, &Headline{
		ID:          "20",
		Title:       "Sit back",
		Description: "We will notify you on kick-off 🛎",
		Images:      []string{},
		Rotation:    1,
		Type:        HeadlineTypeCountdown,
	})

	// 4) Most Picked star player
	mostPickedStar, err := store.GetMostPickedPlayer(matchID, true)
	// if not found ignore
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, "cannot get most picked star player")
		}

		result = append(result, &Headline{
			ID:          "40",
			Title:       mostPickedStar.FullName.String,
			Description: "Most Picked Player 🌟",
			Images:      []string{mostPickedStar.ImageURL.String},
			Rotation:    1,
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 5) Most picked underdog
	mostPickedUnderdog, err := store.GetMostPickedPlayer(matchID, false)
	// if not found ignore
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, "cannot get most picked underdog")
		}

		result = append(result, &Headline{
			ID:          "50",
			Title:       mostPickedUnderdog.FullName.String,
			Description: "Most Picked Underdog ✨",
			Images:      []string{mostPickedUnderdog.ImageURL.String},
			Rotation:    1,
			Type:        HeadlineTypeFadeIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 6) Highest picked team
	highestPickedTeam, pickPercent, err := store.GetHighestTeamPickPercent(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get highest picked team")
	}

	if highestPickedTeam != nil {
		result = append(result, &Headline{
			ID:          "60",
			Title:       fmt.Sprintf("%0.1f%%", pickPercent),
			Description: fmt.Sprintf("players picked from %s ☄️", highestPickedTeam.Name),
			Images:      []string{highestPickedTeam.CrestURL.String},
			Rotation:    1,
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypeTeam,
		})
	}

	// 7) Home team scores
	homeTeam, err := store.GetTeamByID(match.HomeTeamID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get home team")
	}

	homeTeamAvgScore, err := store.GetAvgTeamScore(homeTeam.ID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get home team avg score")
	}

	result = append(result, &Headline{
		ID:          "70",
		Title:       fmt.Sprintf("%s scores %0.0f pts", homeTeam.Name, homeTeamAvgScore),
		Description: "on average",
		Images:      []string{homeTeam.CrestURL.String},
		Rotation:    1,
		Type:        HeadlineTypeSlideIn,
		ImageType:   HeadlineImageTypeTeam,
	})

	// 8) Away team scores
	awayTeam, err := store.GetTeamByID(match.AwayTeamID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get away team")
	}

	awayTeamAvgScore, err := store.GetAvgTeamScore(awayTeam.ID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get away team avg score")
	}

	result = append(result, &Headline{
		ID:          "80",
		Title:       fmt.Sprintf("%s scores %0.0f pts", awayTeam.Name, awayTeamAvgScore),
		Description: "on average",
		Images:      []string{awayTeam.CrestURL.String},
		Rotation:    1,
		Type:        HeadlineTypeFadeIn,
		ImageType:   HeadlineImageTypeTeam,
	})

	// 9) Welcome to new players
	newPlayersCount, err := store.GetNewPlayersCount(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get new players count")
	}

	if newPlayersCount > 0 {
		newPlayersAvatars, err := store.GetNewPlayerAvatars(matchID, 6)
		if err != nil {
			return nil, stacktrace.Propagate(err, "cannot get new players avatars")
		}

		result = append(result, &Headline{
			ID:          "90",
			Title:       "Welcome",
			Description: fmt.Sprintf("to %d first time players 👋", newPlayersCount),
			Images:      newPlayersAvatars,
			Rotation:    1,
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypeUser,
		})
	}

	return result, nil
}

func getGameplayHeadlines(store database.Store, matchID string) ([]*Headline, error) {
	result := make([]*Headline, 0)

	match, err := store.GetMatchByIDWithTeams(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get match by id")
	}

	currMin := match.Minute

	// New Offer Headline every 10 minutes since the start of the game
	if currMin > 0 && currMin%15 == 0 {
		result = append(result, &Headline{
			ID:          fmt.Sprintf("offer-%d", currMin), // Unique ID for each occurrence
			Title:       "NEW OFFERS",
			Description: "buy your avatars today!",
			Images:      []string{}, // No images required for this headline
			Rotation:    1,
			Type:        HeadlineTypeSlideIn,   // Adjust the headline type if necessary
			ImageType:   HeadlineImageTypeUser, // Adjust if a specific image type is needed
		})
	}

	// 1) Top player within last 10 mins
	topPlayer, points, err := store.GetTopPlayer(matchID, currMin-10, currMin)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, fmt.Sprintf("GetTopPlayer failed (match_id: %s)", matchID))
		}

		result = append(result, &Headline{
			ID:          "1010",
			Title:       fmt.Sprintf("%s top player with %0.0f pts", topPlayer.FullName.String, points),
			Description: fmt.Sprintf("from %d to %d minute 🔥", fixMin(currMin-10), fixMin(currMin)),
			Images:      []string{topPlayer.ImageURL.String},
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 2) Worst player within last 10 mins
	worstPlayer, points, err := store.GetWorstPlayer(matchID, currMin-10, currMin)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, fmt.Sprintf("GetTopPlayer failed (match_id: %s)", matchID))
		}

		result = append(result, &Headline{
			ID:          "1020",
			Title:       fmt.Sprintf("%s worst player with %0.0f pts", worstPlayer.FullName.String, points),
			Description: fmt.Sprintf("from %d to %d minute 👎", fixMin(currMin-10), fixMin(currMin)),
			Images:      []string{worstPlayer.ImageURL.String},
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 3) Most traded player within last 10 mins
	mostTraded, tradeCount, err := store.GetMostTraded(matchID, currMin-10, currMin)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, fmt.Sprintf("GetMostTraded failed (match_id: %s)", matchID))
		}

		result = append(result, &Headline{
			ID:          "1030",
			Title:       fmt.Sprintf("%s traded in %d times", mostTraded.FullName.String, tradeCount),
			Description: fmt.Sprintf("from %d to %d minute ⬆️", fixMin(currMin-10), fixMin(currMin)),
			Images:      []string{mostTraded.ImageURL.String},
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 4) Powerups used in last 5 mins
	puUseCount, err := store.GetPowerUpUsages(matchID, currMin-5, currMin)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("GetPowerUpUsages failed (match_id: %s)", matchID))
	}
	if puUseCount > 0 {
		result = append(result, &Headline{
			ID:          "1040",
			Title:       fmt.Sprintf("%d powerups used", puUseCount),
			Description: fmt.Sprintf("from %d to %d minute", fixMin(currMin-5), fixMin(currMin)),
			// TODO: default powerup icon
			Images: []string{},
			Type:   HeadlineTypeSlideIn,
		})
	}

	// 5) Most used powerup within last 5 mins
	mostUsedPu, puUserCount, err := store.GetMostUsedPowerUp(matchID, currMin-5, currMin)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, fmt.Sprintf("GetMostUsedPowerUp failed (match_id: %s)", matchID))
		}

		result = append(result, &Headline{
			ID:          "1050",
			Title:       fmt.Sprintf("%d users used %s", puUserCount, mostUsedPu.Name),
			Description: fmt.Sprintf("from %d to %d minute 💥", fixMin(currMin-5), fixMin(currMin)),
			Images:      []string{mostUsedPu.IconURL.String},
			Type:        HeadlineTypeSlideIn,
		})
	}

	// 6) Best user in period
	bestUser, points, err := store.GetBestUser(matchID, currMin-10, currMin)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, fmt.Sprintf("GetBestUser failed (match_id: %s)", matchID))
		}

		result = append(result, &Headline{
			ID:          "1060",
			Title:       fmt.Sprintf("%s best user with %0.0f pts", bestUser.Name, points),
			Description: fmt.Sprintf("from %d to %d minute 💥", fixMin(currMin-10), fixMin(currMin)),
			Images:      []string{bestUser.AvatarURL.String},
			ImageType:   HeadlineImageTypeUser,
		})
	}

	// 7) Best team for game
	scoreInfo, err := store.GetTeamsWithScore(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("GetTeamsWithScore failed (match_id: %s)", matchID))
	}

	if len(scoreInfo) >= 2 {
		result = append(result, &Headline{
			ID:          "1070",
			Title:       fmt.Sprintf("%s has %0.0f total team pts", scoreInfo[0].Team.Name, scoreInfo[0].Value),
			Description: fmt.Sprintf("%s has %0.0f pts", scoreInfo[1].Team.Name, scoreInfo[1].Value),
			Images:      []string{scoreInfo[0].Team.CrestURL.String},
			ImageType:   HeadlineImageTypeTeam,
		})
	}

	// 8) Teams with possession
	possessionInfo, err := store.GetTeamsWithPossession(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("GetTeamsWithPossession failed (match_id: %s)", matchID))
	}

	if len(possessionInfo) >= 1 {
		result = append(result, &Headline{
			ID:          "1080",
			Title:       fmt.Sprintf("%s has %0.0f%% possession", possessionInfo[0].Team.Name, possessionInfo[0].Value*100.0),
			Description: "",
			Images:      []string{possessionInfo[0].Team.CrestURL.String},
			ImageType:   HeadlineImageTypeTeam,
		})
	}

	// 9) Teams with shots
	shotsInfo, err := store.GetTeamsWithShots(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("GetTeamsWithShots failed (match_id: %s)", matchID))
	}

	if len(shotsInfo) > 0 {
		var shotsTitle = fmt.Sprintf("%s has %0.0f total shots", shotsInfo[0].Team.Name, shotsInfo[0].Value)

		var shotsDescription string
		if len(shotsInfo) > 1 {
			shotsDescription = fmt.Sprintf("%s has %0.0f shots", shotsInfo[1].Team.Name, shotsInfo[1].Value)
		}

		result = append(result, &Headline{
			ID:          "1090",
			Title:       shotsTitle,
			Description: shotsDescription,
			Images:      []string{shotsInfo[0].Team.CrestURL.String},
			ImageType:   HeadlineImageTypeTeam,
		})
	}

	if len(result) > 0 {
		result = append(result, &Headline{
			ID:          "1100",
			Title:       "",
			Description: "",
			Type:        HeadlineTypeLastPickAction,
		})
	}

	// 10) Get Involved
	result = append(result, &Headline{
		ID:          "1110",
		Title:       "Get Involved",
		Description: "Chat below with other fans 💬",
		Images:      []string{},
		Type:        HeadlineTypeSlideIn,
		ImageType:   HeadlineImageTypeUser,
	})

	return result, nil
}

func getFullTimeHeadlines(store database.Store, matchID string) ([]*Headline, error) {
	result := make([]*Headline, 0)

	// 1) match winner
	matchWinner, winnerScore, err := store.GetMatchWinner(matchID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log or handle the case where there are no rows gracefully
			fmt.Println("No match winner data available for match_id:", matchID)
			return result, nil // Return an empty slice and nil error
		}
		return nil, stacktrace.Propagate(err, "GetMatchWinner failed (match_id: %s)", matchID)
	}
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, "GetMatchWinner failed (match_id: %s)", matchID)
		}

		// // find reward for match winner
		// rewardTransaction, err := store.GetUserTransaction(matchWinner.ID, matchID)
		// if err != nil {
		// 	return nil, stacktrace.Propagate(err, "GetUserTransaction failed (user_id: %s, match_id: %s)", matchWinner.ID, matchID)
		// }

		result = append(result, &Headline{
			ID:          "2010",
			Title:       fmt.Sprintf("%s is global winner", matchWinner.Name),
			Description: fmt.Sprintf("with %0.0f Points! 🏆", winnerScore),
			Images:      []string{matchWinner.AvatarURL.String},
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypeUser,
		})
	}

	// 2) top player
	topPlayer, topPlayerScore, err := store.GetTopPlayer(matchID, 0, 999)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, "GetTopPlayer failed (match_id: %s)", matchID)
		}

		result = append(result, &Headline{
			ID:          "2020",
			Title:       fmt.Sprintf("%s top player with %0.0f pts", topPlayer.FullName.String, topPlayerScore),
			Description: "in the match 🔥",
			Images:      []string{topPlayer.ImageURL.String},
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 3) worst player
	worstPlayer, worstPlayerScore, err := store.GetWorstPlayer(matchID, 0, 999)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, "GetWorstPlayer failed (match_id: %s)", matchID)
		}

		result = append(result, &Headline{
			ID:          "2030",
			Title:       fmt.Sprintf("%s worst player with %0.0f pts", worstPlayer.FullName.String, worstPlayerScore),
			Description: "in the match 👎",
			Images:      []string{worstPlayer.ImageURL.String},
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 4) most traded
	mostTraded, mostTradedCount, err := store.GetMostTraded(matchID, 0, 999)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, "GetMostTraded failed (match_id: %s)", matchID)
		}

		result = append(result, &Headline{
			ID:          "2040",
			Title:       fmt.Sprintf("%s traded in %d times", mostTraded.FullName.String, mostTradedCount),
			Description: "in the match 📈",
			Images:      []string{mostTraded.ImageURL.String},
			Type:        HeadlineTypeSlideIn,
			ImageType:   HeadlineImageTypePlayer,
		})
	}

	// 5) number of used power-ups
	usedPowerUps, err := store.GetPowerUpUsages(matchID, 0, 999)
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, stacktrace.Propagate(err, "GetPowerUpUsages failed (match_id: %s)", matchID)
		}

		result = append(result, &Headline{
			ID:          "2050",
			Title:       fmt.Sprintf("%d powerups used", usedPowerUps),
			Description: "in the match 🚀",
			// TODO: Default Powerup Icon (lightening)
			Images: []string{},
			Type:   HeadlineTypeSlideIn,
		})
	}

	return result, nil
}

func toMatchHeadlines(matchID string, screenType int, headlines []*Headline) schema.MatchHeadlineSlice {
	var result schema.MatchHeadlineSlice
	for _, el := range headlines {
		if len(el.Images) == 0 {
			el.Images = make([]string, 0)
		}
		// serialize images
		b, _ := json.Marshal(el.Images)
		if b == nil {
			b = []byte("[]")
		}

		if el.ImageType == "" {
			el.ImageType = "raw"
		}

		if el.Type == "" {
			el.Type = HeadlineTypeSlideIn
		}

		result = append(result, &schema.MatchHeadline{
			ScreenType:  screenType,
			Title:       el.Title,
			Description: el.Description,
			Type:        el.Type,
			MatchID:     matchID,
			Images:      string(b),
			ImageType:   null.StringFrom(el.ImageType),
		})
	}

	return result
}

func fixMin(minute int) int {
	if minute <= 0 {
		return 0
	}
	return minute
}
