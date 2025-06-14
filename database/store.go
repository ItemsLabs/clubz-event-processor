package database

import (
	"time"

	"github.com/volatiletech/sqlboiler/types"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
)

const (
	MatchStatusUnknown   = "u"
	MatchStatusWaiting   = "w"
	MatchStatusLineups   = "l"
	MatchStatusGame      = "g"
	MatchStatusEnded     = "e"
	MatchStatusCancelled = "c"

	MatchPeriodPregame    = "p"
	MatchPeriodFirstHalf  = "f"
	MatchPeriodHalfTime   = "h"
	MatchPeriodSecondHalf = "s"
	MatchPeriodBreakX1    = "bx1"
	MatchPeriodFirstExt   = "x1"
	MatchPeriodBreakX2    = "bx2"
	MatchPeriodSecondExt  = "x2"
	MatchPeriodBreakP     = "bp"
	MatchPeriodPenalties  = "pe"
	MatchPeriodPostGame   = "pg"

	GameStatusWaiting  = "w"
	GameStatusGameplay = "g"
	GameStatusFinished = "f"

	MatchEventStatusActive    = 1
	MatchEventStatusCancelled = 2
	MatchEventStatusIgnored   = 3

	PositionGoalkeeper = "g"
	PositionDefender   = "d"
	PositionMidfielder = "m"
	PositionForward    = "f"

	SubscriptionTierNone    = 0
	SubscriptionTierPremium = 1
	SubscriptionTierLite    = 2
)

type LineupPickInfo struct {
	UserID         string
	GameID         string
	NonLineupPicks types.StringArray `boil:""`
}

type TeamWithValue struct {
	Team  *schema.Team
	Value float64
}

type PlayerWithNFTDetails struct {
	Player *schema.AssignedPlayer
	NFT    *schema.NFTBucket
}

type Store interface {
	Transaction(func(s Store) error) error

	GetOrCreateMatchEventProcessor(matchID string, typ int) (*schema.MatchEventProcessor, error)
	GetMatchEventProcessorLastProcessedID(processorID string) (int, error)
	SetMatchEventProcessorLastProcessedID(processorID string, lastProcessedID int) error

	GetActivePicksAtTime(matchID string, playerID string, timestamp time.Time) (schema.GamePickSlice, error)
	GetActivePicksAtMinSec(matchID string, playerID string, min, sec int) (schema.GamePickSlice, error)
	UpdatePickScore(pick *schema.GamePick) error
	UpdatePickEndedAt(pick *schema.GamePick) error
	InsertPick(pick *schema.GamePick) (*schema.GamePick, error)
	GetGamePicks(gameID string) (schema.GamePickSlice, error)
	GetGamesWithoutPowerUps(matchID string) (schema.GameSlice, error)
	GetPowerUpCountForGame(matchID string) (map[string]int, error)

	GetUserByID(userID string) (*schema.User, error)

	IncGameVersion(gameID string) error

	UpdateMatchGamesStatus(matchID string, state string) error
	SyncGamePremiumFlags(matchID string) error

	InsertGameEvent(ev *schema.GameEvent) error
	DeleteGameEventsByMatchEventID(id int) error

	InsertUserTransaction(t *schema.Transaction) (*schema.Transaction, error)
	GetUserTransaction(userID, matchID string) (*schema.Transaction, error)

	GetMatchEvents(matchID string) (schema.MatchEventSlice, error)
	GetMatchEventByID(id int) (*schema.MatchEvent, error)
	GetMatchEventByMatchEventID(matchID string, matchEventID int) (*schema.MatchEvent, error)
	UpdateMatchEventStatus(ev *schema.MatchEvent) error
	GetLatestMatchEvent(matchID string) (*schema.MatchEvent, error)

	GetMatchByIDWithTeams(matchID string) (*schema.Match, error)
	GetMatchByID(matchID string) (*schema.Match, error)
	UpdateMatch(match *schema.Match, updatedFields []string) (*schema.Match, error)
	UpdateMatchTime(matchID string, minute, second int) error
	IncHomeScore(match *schema.Match) (*schema.Match, error)
	IncAwayScore(match *schema.Match) (*schema.Match, error)
	DecHomeScore(match *schema.Match) (*schema.Match, error)
	DecAwayScore(match *schema.Match) (*schema.Match, error)

	GetPlayerByID(playerID string) (*schema.Player, error)

	GetMatchPlayer(matchID string, teamID string, playerID string) (*schema.MatchPlayer, error)
	GetMatchPlayerWithoutTeam(matchID string, playerID string) (*schema.MatchPlayer, error)
	GetOrCreateMatchPlayer(matchID string, teamID string, playerID string) (*schema.MatchPlayer, error)
	InsertMatchPlayer(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error)
	UpdateMatchPlayer(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error)
	GetMatchPlayers(matchID string) (schema.MatchPlayerSlice, error)
	UpdateMatchPlayerScore(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error)
	UpdateMatchPlayerPosition(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error)
	UpdateMatchPlayerPlayedSeconds(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error)
	IncMatchPlayerScoreFast(matchID string, teamID string, playerID string, score float64) error

	GetTeamByID(teamID string) (*schema.Team, error)

	InsertAMQPEvent(ev *schema.AmqpEvent) (*schema.AmqpEvent, error)

	GetPowerUpActions() (schema.PowerupActionSlice, error)
	GetActiveGamePowerUpsAtTime(matchID string, t time.Time) (schema.GamePowerupSlice, error)
	UpdateEndedAt(*schema.GamePowerup) error

	UpdateMatchLeaderboard(matchID string) error
	SetTransactionIDForMatchLeaderboard(userID, matchID, transactionID string) error
	GetMatchLeaderboard(matchID string) (schema.MatchLeaderboardSlice, error)
	GetTopMatchLeaderboard(matchID string, maxPosition int) (schema.MatchLeaderboardSlice, error)
	GetFullMatchLeaderboard(matchID string) (schema.MatchLeaderboardSlice, error)
	GetMatchLeaderboardAtTime(matchID string, t time.Time) (schema.MatchLeaderboardSlice, error)
	GetMatchRewards(matchID string) (schema.MatchRewardSlice, error)

	GetGameCount(matchID string) (int, error)
	GetRandomUserAvatars(matchID string, count int) ([]string, error)
	GetMostPickedPlayer(matchID string, isStar bool) (*schema.Player, error)
	GetHighestTeamPickPercent(matchID string) (*schema.Team, float64, error)
	GetAvgTeamScore(teamID string) (float64, error)
	GetNewPlayersCount(matchID string) (int, error)
	GetNewPlayerAvatars(matchID string, count int) ([]string, error)
	GetMatchesForHeadlines(t time.Time) (schema.MatchSlice, error)
	GetCurrentMatchTime(matchID string) (int, int, error)
	GetCurrentMatchMinute(matchID string) (int, error)
	GetTopPlayerByAvg(matchID string) (*schema.Player, error)
	GetAvgPlayerGoals(playerID string) (float64, error)
	GetTopPlayer(matchID string, fromMin int, toMin int) (*schema.Player, float64, error)
	GetWorstPlayer(matchID string, fromMin int, toMin int) (*schema.Player, float64, error)
	GetMostTraded(matchID string, fromMin, toMin int) (*schema.Player, int, error)
	GetPowerUpUsages(matchID string, fromMin, toMin int) (int, error)
	GetMostUsedPowerUp(matchID string, fromMin, toMin int) (*schema.Powerup, int, error)
	GetBestUser(matchID string, fromMin, toMin int) (*schema.User, float64, error)
	GetTeamsWithScore(matchID string) ([]*TeamWithValue, error)
	GetTeamsWithPossession(matchID string) ([]*TeamWithValue, error)
	GetTeamsWithShots(matchID string) ([]*TeamWithValue, error)
	GetMatchReward(matchID string) (float64, error)
	GetMatchWinner(matchID string) (*schema.User, float64, error)
	UpdateHeadlines(matchID string, headlines schema.MatchHeadlineSlice) error

	GetLineupsPickInfo(matchID string) ([]*LineupPickInfo, error)

	GetMatchNotification(matchID string, userID string, typ int) (*schema.MatchNotification, error)
	CreateMatchNotification(notification *schema.MatchNotification) error
	GetActionNameByActionID(actionID int) (string, error)
	GetUserIDsByMatchID(matchID string) ([]string, error)
	UpdateFinishedGamesByUserID(userID, gameID string) error
	GetGameByUserIDMatchID(userID, matchID string) (*schema.Game, error)
	GetUserIDsNotInMatch(matchID string) ([]string, error)

	// Newly added functions
	CreateNFTBucket(nb *schema.NFTBucket) (*schema.NFTBucket, error)
	GetNFTBucketByID(id string) (*schema.NFTBucket, error)
	GetAllNFTBuckets() (schema.NFTBucketSlice, error)
	UpdateNFTBucket(nb *schema.NFTBucket) error
	DeleteNFTBucket(id string) error
	GetNFTBucketsByTeamID(teamID string) (schema.NFTBucketSlice, error)
	CountNFTBuckets() (int64, error)
	GetNFTsByRarity(rarity string) (schema.NFTBucketSlice, error)
	GetNFTsByTeamAndRarity(teamID, rarity string) (schema.NFTBucketSlice, error)
	GetRandomNFTByTeamAndRarity(teamID string, rarity string) (*schema.NFTBucket, error)
	GetRandomNFTBucketByTeamID(teamID string) (*schema.NFTBucket, error)
	GetNFTBucketsByPosition(position string) (schema.NFTBucketSlice, error)
	GetNFTBucketsByGamePosition(gamePosition string) (schema.NFTBucketSlice, error)
	GetNFTBucketsByPositionAndTeam(teamID, position string) (schema.NFTBucketSlice, error)
	GetNFTBucketsByTeamAndNationality(teamID, nationality string) (schema.NFTBucketSlice, error)
	GetNFTBucketsByNationality(nationality string) (schema.NFTBucketSlice, error)
	GetNFTBucketsByAgeRange(minAge, maxAge int) (schema.NFTBucketSlice, error)
	GetNFTBucketsByRarityAndPosition(rarity, position string) (schema.NFTBucketSlice, error)
	CountNFTBucketsByRarity() (map[string]int64, error)
	AssignNFTToTeam(nftID, teamID string) (*schema.NFTBucket, error)
	GetRandomNFTBucket() (*schema.NFTBucket, error)
	GetNFTBucketByName(name string) (*schema.NFTBucket, error)
	GetRandomNFTBucketByTeamAndRarityAndStarRating(teamID string, rarity string, starRating int) (*schema.NFTBucket, error)

	// Assigned player functions
	CreateAssignedPlayer(ap *schema.AssignedPlayer) (*schema.AssignedPlayer, error)
	GetAssignedPlayerByID(id string) (*schema.AssignedPlayer, error)
	GetAllAssignedPlayers() (schema.AssignedPlayerSlice, error)
	UpdateAssignedPlayer(ap *schema.AssignedPlayer) error
	DeleteAssignedPlayer(id string) error
	GetAssignedPlayersByUserID(userID string) (schema.AssignedPlayerSlice, error)
	GetAssignedPlayersByPlayerID(playerID string) (schema.AssignedPlayerSlice, error)
	PlayerPopularity() (map[string]int, error)
	UserEngagement() (map[string]float64, error)
	TeamBalanceAnalysis() (map[string]map[string]int, error)
	CountAllEntries() (int64, error)
	GetAssignedPlayersByUUIDs(uuids []string) (schema.AssignedPlayerSlice, error)
	GetAssignedPlayersWithNFTDetails(uuids []string) ([]struct {
		Player *schema.AssignedPlayer
		NFT    *schema.NFTBucket
	}, error)
	GetAssignedPlayersByNFTIDs(nftIDs []string) ([]PlayerWithNFTDetails, error)

	GetActionsByTypeID(typeID int) (schema.ActionSlice, error)
	GetAllUserIDsFromTable() ([]string, error)

	CreateReward(reward *schema.Reward) (*schema.Reward, error)

	CreateAppInbox(inbox *schema.AppInbox) (*schema.AppInbox, error)

	GetCurrentGameWeek() (*schema.GameWeek, error)
	InsertAssignedCardPack(acp *schema.AssignedCardPack) (*schema.AssignedCardPack, error)

	InsertPushNotification(pushNotification *schema.PushNotification) error
	CountUserNotificationsLastHour(userID string) (int, error)
}
