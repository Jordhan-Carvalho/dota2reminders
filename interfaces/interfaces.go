package interfaces

/* DOTA_GameState
        /// Undefined
        Undefined,

        /// Disconnected
        DOTA_GAMERULES_STATE_DISCONNECT,

        /// Game is in progress
        DOTA_GAMERULES_STATE_GAME_IN_PROGRESS,

        /// Players are currently selecting heroes
        DOTA_GAMERULES_STATE_HERO_SELECTION,

        /// Game is starting
        DOTA_GAMERULES_STATE_INIT,

        /// Game is ending
        DOTA_GAMERULES_STATE_LAST,

        /// Game has ended, post game scoreboard
        DOTA_GAMERULES_STATE_POST_GAME,

        /// Game has started, pre game preparations
        DOTA_GAMERULES_STATE_PRE_GAME,

        /// Players are selecting/banning heroes
        DOTA_GAMERULES_STATE_STRATEGY_TIME,

        /// Waiting for everyone to connect and load
        DOTA_GAMERULES_STATE_WAIT_FOR_PLAYERS_TO_LOAD,

        /// Game is a custom game
        DOTA_GAMERULES_STATE_CUSTOM_GAME_SETUP */

type mapEvent struct {
	Name      string `json:"name"`
	ClockTime int    `json:"clock_time"`
	GameState string `json:"game_state"` 
  WardPurchaseCooldown int `json:"ward_purchase_cooldown"`
  RadiantScore int `json:"radiant_score"`
  DireScore int `json:"dire_score"`
  Paused bool `json:"paused"`
}

type GameEvents struct {
	Map mapEvent `json:"map"`
}
