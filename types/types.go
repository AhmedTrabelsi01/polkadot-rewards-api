package types

type SubscanRewardSlashRequest struct {
	Address  string `json:"address"`
	Page     int    `json:"page"`
	Row      int    `json:"row"`
	Category string `json:"category,omitempty"`
}

type SubscanRewardItem struct {
	ID             int64  `json:"id"`
	BlockNum       int64  `json:"block_num"`
	ExtrinsicIdx   int    `json:"extrinsic_idx"`
	Stash          string `json:"stash"`
	Account        string `json:"account"`
	ModuleID       string `json:"module_id"`
	EventID        string `json:"event_id"`
	EventMethod    string `json:"event_method"`
	Params         string `json:"params"`
	ExtrinsicHash  string `json:"extrinsic_hash"`
	EventIdx       int    `json:"event_idx"`
	Amount         string `json:"amount"`
	BlockTimestamp int64  `json:"block_timestamp"`
	EventIndex     string `json:"event_index"`
	ExtrinsicIndex string `json:"extrinsic_index"`
}

type SubscanRewardSlashData struct {
	Count int                 `json:"count"`
	List  []SubscanRewardItem `json:"list"`
}

type SubscanRewardSlashResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    SubscanRewardSlashData `json:"data"`
}

type EraReward struct {
	Era    int    `json:"era"`
	Amount string `json:"amount"`
}

type DailyReward struct {
	Date   string `json:"date"`
	Amount string `json:"amount"`
}

type MonthlyReward struct {
	Month  string `json:"month"`
	Amount string `json:"amount"`
}

type EraRewardsResponse struct {
	Validator string      `json:"validator"`
	Rewards   []EraReward `json:"rewards"`
}

type DailyRewardsResponse struct {
	Validator string        `json:"validator"`
	Rewards   []DailyReward `json:"rewards"`
}

type MonthlyRewardsResponse struct {
	Validator string          `json:"validator"`
	Rewards   []MonthlyReward `json:"rewards"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

