package model

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/common/utils/conv"
	"go.mongodb.org/mongo-driver/bson"
)

const CollectionUserAccounts = "user_accounts"

var (
	_ Document         = (*UserAccount)(nil)
	_ bson.Marshaler   = (*UserAccount)(nil)
	_ bson.Unmarshaler = (*UserAccount)(nil)
)

type UserAccount struct {
	UserId            string          `bson:"_id" json:"userId"`
	Password          string          `bson:"password" json:"password"`
	Owner             *common.Address `bson:"owner,omitempty" json:"owner,omitempty"`
	Account           *common.Address `bson:"account,omitempty" json:"account,omitempty"`
	CreatedAt         int64           `bson:"createdAt" json:"createdAt"`
	LastFaucet        int64           `bson:"lastFaucet" json:"lastFaucet"`
	LastUsedNonce     uint64          `bson:"lastUsedNonce" json:"lastUsedNonce"`
	LastTxnHash       string          `bson:"lastTxnHash" json:"lastTxnHash"`
	Pending           bool            `bson:"pending" json:"pending"`
	SyncedBlockNumber uint64          `bson:"syncedBlockNumber" json:"syncedBlockNumber"`
}

func (d UserAccount) ID() string         { return d.UserId }
func (d UserAccount) Collection() string { return CollectionUserAccounts }

func (d UserAccount) MarshalBSON() ([]byte, error) {
	if d.UserId == "" || d.Password == "" {
		return nil, ErrInsufficientField
	}
	var owner *string
	if d.Owner != nil {
		owner = conv.ToPtr(d.Owner.Hex())
	}
	var account *string
	if d.Account != nil {
		account = conv.ToPtr(d.Account.Hex())
	}
	var createdAt int64
	if d.CreatedAt == 0 {
		createdAt = time.Now().Unix()
	}
	type UserAccount struct {
		UserId            string  `bson:"_id"`
		Password          string  `bson:"password"`
		Owner             *string `bson:"owner,omitempty"`
		Account           *string `bson:"account,omitempty"`
		CreatedAt         int64   `bson:"createdAt"`
		LastFaucet        int64   `bson:"lastFaucet"`
		LastUsedNonce     uint64  `bson:"lastUsedNonce"`
		LastTxnHash       string  `bson:"lastTxnHash"`
		Pending           bool    `bson:"pending"`
		SyncedBlockNumber uint64  `bson:"syncedBlockNumber"`
	}
	return bson.Marshal(&UserAccount{
		UserId:            d.UserId,
		Password:          d.Password,
		Owner:             owner,
		Account:           account,
		CreatedAt:         createdAt,
		LastFaucet:        d.LastFaucet,
		LastUsedNonce:     d.LastUsedNonce,
		LastTxnHash:       d.LastTxnHash,
		Pending:           d.Pending,
		SyncedBlockNumber: d.SyncedBlockNumber,
	})
}

func (d *UserAccount) UnmarshalBSON(data []byte) error {
	type UserAccount struct {
		UserId            string  `bson:"_id"`
		Password          string  `bson:"password"`
		Owner             *string `bson:"owner,omitempty"`
		Account           *string `bson:"account,omitempty"`
		CreatedAt         int64   `bson:"createdAt"`
		LastFaucet        int64   `bson:"lastFaucet"`
		LastUsedNonce     uint64  `bson:"lastUsedNonce"`
		LastTxnHash       string  `bson:"lastTxnHash"`
		Pending           bool    `bson:"pending"`
		SyncedBlockNumber uint64  `bson:"syncedBlockNumber"`
	}
	var dec UserAccount
	if err := bson.Unmarshal(data, &dec); err != nil {
		return err
	}
	d.UserId = dec.UserId
	d.Password = dec.Password
	if dec.Owner != nil {
		d.Owner = conv.ToPtr(common.HexToAddress(*dec.Owner))
	}
	if dec.Account != nil {
		d.Account = conv.ToPtr(common.HexToAddress(*dec.Account))
	}
	d.CreatedAt = dec.CreatedAt
	d.LastFaucet = dec.LastFaucet
	d.LastUsedNonce = dec.LastUsedNonce
	d.LastTxnHash = dec.LastTxnHash
	d.Pending = dec.Pending
	d.SyncedBlockNumber = dec.SyncedBlockNumber
	return nil
}
