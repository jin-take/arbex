package domain

import "time"

type VenueID string
type AccountID string
type AssetID string
type InstrumentID string
type RiskFactorID string
type LiquiditySourceID string
type ExecutionGroupID string
type ChildOrderID string
type FillID string
type ReservationID string

type Token struct {
	ChainID  uint64
	Address  string
	Symbol   string
	Decimals uint8
	AssetID  AssetID
	RiskID   RiskFactorID
}

type PoolSnapshot struct {
	ChainID      uint64
	PoolAddress  string
	BlockNumber  uint64
	BlockHash    string
	StateRevision uint64
	ObservedAt   time.Time
}

type Route struct {
	ID         string
	StartAsset AssetID
	EndAsset   AssetID
	Sources    []LiquiditySourceID
}

type Quote struct {
	RouteID       string
	Input         Amount
	Output        Amount
	StateRevision uint64
	QuotedAt      time.Time
}

type Opportunity struct {
	ID        string
	Quote     Quote
	ExpiresAt time.Time
}

type ExecutionAttempt struct {
	GroupID        ExecutionGroupID
	ConfigRevision uint64
	CreatedAt      time.Time
}

type ExecutionGroup struct {
	ID             ExecutionGroupID
	ConfigRevision uint64
	CreatedAt      time.Time
}

type ChildOrder struct {
	ID        ChildOrderID
	GroupID   ExecutionGroupID
	VenueID   VenueID
	AccountID AccountID
	AssetID   AssetID
}

type Fill struct {
	ID       FillID
	OrderID  ChildOrderID
	Quantity Amount
	Fee      Amount
}

type Reservation struct {
	ID        ReservationID
	VenueID   VenueID
	AccountID AccountID
	AssetID   AssetID
	Amount    Amount
	ExpiresAt *time.Time
}

type Exposure struct {
	RiskID   RiskFactorID
	Known    Amount
	Possible Amount
}
