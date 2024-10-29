package model

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/common/utils/conv"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionMngdAccounts                 = "managed_accounts"
	KindEOA                mngdAccountKind = "eoa"
	KindContract           mngdAccountKind = "contract"

	RoleSupervisor mngdAccountRole = "supervisor" // main-ca-owner
	RoleManager    mngdAccountRole = "manager"    // exec createAccount, faucet
	RoleOperator   mngdAccountRole = "operator"   // exec userOp
)

var (
	_ Document         = (*ManagedAccount)(nil)
	_ bson.Marshaler   = (*ManagedAccount)(nil)
	_ bson.Unmarshaler = (*ManagedAccount)(nil)
)

type (
	mngdAccountKind string
	mngdAccountRole string

	ManagedAccount struct {
		Id             string            `bson:"_id" json:"_id"`
		Kind           mngdAccountKind   `bson:"kind" json:"kind"`
		Role           *mngdAccountRole  `bson:"role,omitempty" json:"role,omitempty"`
		Address        common.Address    `bson:"address" json:"address"`
		PrivateKey     *ecdsa.PrivateKey `bson:"privateKey,omitempty" json:"privateKey,omitempty"`
		DeployedTxHash *string           `bson:"deployedTxn,omitempty" json:"deployedTxn,omitempty"`
		Deployer       *common.Address   `bson:"deployer,omitempty" json:"deployer,omitempty"`
		Occupied       bool              `bson:"occupied" json:"occupied"`
		OccupiedAt     int64             `bson:"occupiedAt" json:"occupiedAt"`
	}
	MngdWallet struct {
		Supervisor ManagedAccount
		Manager    ManagedAccount
		Operators  []ManagedAccount
	}
)

func (d ManagedAccount) ID() string         { return d.Id }
func (d ManagedAccount) Collection() string { return CollectionMngdAccounts }

func (d ManagedAccount) MarshalBSON() ([]byte, error) {
	var privateKeyHex *string
	if d.Kind == KindEOA {
		pvhex, err := ethutil.PvKey2Hex(d.PrivateKey)
		if err != nil {
			return nil, err
		}
		privateKeyHex = &pvhex
	}
	var deployer *string
	if d.Deployer != nil {
		deployer = conv.ToPtr(d.Deployer.Hex())
	}
	type ManagedAccount struct {
		Id             string           `bson:"_id"`
		Kind           mngdAccountKind  `bson:"kind"`
		Role           *mngdAccountRole `bson:"role,omitempty"`
		Address        string           `bson:"address"`
		PrivateKeyHex  *string          `bson:"privateKey,omitempty"`
		DeployedTxHash *string          `bson:"deployedTxn,omitempty"`
		Deployer       *string          `bson:"deployer,omitempty"`
		Occupied       bool             `bson:"occupied"`
		OccupiedAt     int64            `bson:"occupiedAt"`
	}
	return bson.Marshal(&ManagedAccount{
		Id:             d.Id,
		Kind:           d.Kind,
		Role:           d.Role,
		Address:        d.Address.Hex(),
		PrivateKeyHex:  privateKeyHex,
		DeployedTxHash: d.DeployedTxHash,
		Deployer:       deployer,
		OccupiedAt:     d.OccupiedAt,
		Occupied:       d.Occupied,
	})
}

func (d *ManagedAccount) UnmarshalBSON(data []byte) error {
	type ManagedAccount struct {
		Id             string           `bson:"_id"`
		Kind           mngdAccountKind  `bson:"kind"`
		Role           *mngdAccountRole `bson:"role,omitempty"`
		Address        string           `bson:"address"`
		PrivateKeyHex  *string          `bson:"privateKey,omitempty"`
		DeployedTxHash *string          `bson:"deployedTxn,omitempty"`
		Deployer       *string          `bson:"deployer,omitempty"`
		Occupied       bool             `bson:"occupied"`
		OccupiedAt     int64            `bson:"occupiedAt"`
	}
	var dec ManagedAccount
	err := bson.Unmarshal(data, &dec)
	if err != nil {
		return err
	}
	d.Id = dec.Id
	d.Kind = dec.Kind
	d.Role = dec.Role
	d.Address = common.HexToAddress(dec.Address)
	if dec.Kind == KindEOA {
		if d.PrivateKey, err = ethutil.PvHex2Key(*dec.PrivateKeyHex); err != nil {
			return err
		}
	}
	d.DeployedTxHash = dec.DeployedTxHash
	if dec.Deployer != nil {
		d.Deployer = conv.ToPtr(common.HexToAddress(*dec.Deployer))
	}
	d.OccupiedAt = dec.OccupiedAt
	d.Occupied = dec.Occupied
	return nil
}
