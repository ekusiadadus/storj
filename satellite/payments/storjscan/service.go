// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package storjscan

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/common/uuid"
	"storj.io/storj/private/blockchain"
	"storj.io/storj/satellite/payments"
	"storj.io/storj/satellite/payments/billing"
)

// ErrService is storjscan service error class.
var ErrService = errs.Class("storjscan service")

type storjscanMetadata struct {
	ReferenceID string
	Wallet      string
	BlockNumber int64
	LogIndex    int
}

// ensures that storjscan implements payments.DepositWallets.
var _ payments.DepositWallets = (*Service)(nil)

// ensures that storjscan implements billing.PaymentType.
var _ billing.PaymentType = (*Service)(nil)

// Service exposes API to interact with storjscan payments provider.
type Service struct {
	log        *zap.Logger
	walletsDB  WalletsDB
	paymentsDB PaymentsDB
	client     *Client
}

// NewService creates new storjscan service instance.
func NewService(log *zap.Logger, walletsDB WalletsDB, paymentsDB PaymentsDB, client *Client) *Service {
	return &Service{
		log:        log,
		walletsDB:  walletsDB,
		paymentsDB: paymentsDB,
		client:     client,
	}
}

// Claim gets a new crypto wallet and associates it with a user.
func (service *Service) Claim(ctx context.Context, userID uuid.UUID) (_ blockchain.Address, err error) {
	defer mon.Task()(&ctx)(&err)

	address, err := service.client.ClaimNewEthAddress(ctx)
	if err != nil {
		return blockchain.Address{}, ErrService.Wrap(err)
	}
	err = service.walletsDB.Add(ctx, userID, address)
	if err != nil {
		return blockchain.Address{}, ErrService.Wrap(err)
	}

	return address, nil
}

// Get returns the crypto wallet address associated with the given user.
func (service *Service) Get(ctx context.Context, userID uuid.UUID) (_ blockchain.Address, err error) {
	defer mon.Task()(&ctx)(&err)

	address, err := service.walletsDB.GetWallet(ctx, userID)
	return address, ErrService.Wrap(err)
}

// Payments retrieves payments for specific wallet.
func (service *Service) Payments(ctx context.Context, wallet blockchain.Address, limit int, offset int64) (_ []payments.WalletPayment, err error) {
	defer mon.Task()(&ctx)(&err)

	cachedPayments, err := service.paymentsDB.ListWallet(ctx, wallet, limit, offset)
	if err != nil {
		return nil, ErrService.Wrap(err)
	}

	var walletPayments []payments.WalletPayment
	for _, pmnt := range cachedPayments {
		walletPayments = append(walletPayments, payments.WalletPayment{
			From:        pmnt.From,
			To:          pmnt.To,
			TokenValue:  pmnt.TokenValue,
			USDValue:    pmnt.USDValue,
			Status:      pmnt.Status,
			BlockHash:   pmnt.BlockHash,
			BlockNumber: pmnt.BlockNumber,
			Transaction: pmnt.Transaction,
			LogIndex:    pmnt.LogIndex,
			Timestamp:   pmnt.Timestamp,
		})
	}

	return walletPayments, nil
}

// Source defines the billing transaction source for storjscan payments.
func (service *Service) Source() string {
	return "storjscan"
}

// Type defines the billing transaction type for storjscan payments.
func (service *Service) Type() billing.TransactionType {
	return billing.TransactionTypeCredit
}

// GetNewTransactions returns the storjscan payments since the given timestamp as billing transactions type.
func (service *Service) GetNewTransactions(ctx context.Context, _ time.Time, lastPaymentMetadata []byte) ([]billing.Transaction, error) {

	var latestMetadata storjscanMetadata
	if err := json.Unmarshal(lastPaymentMetadata, &latestMetadata); err != nil {
		service.log.Error("error retrieving metadata from latest recorded billing payment")
		return nil, err
	}

	newCachedPayments, err := service.paymentsDB.ListConfirmed(ctx, latestMetadata.BlockNumber, latestMetadata.LogIndex)
	if err != nil {
		return []billing.Transaction{}, ErrService.Wrap(err)
	}

	var billingTXs []billing.Transaction
	for _, payment := range newCachedPayments {
		userID, err := service.walletsDB.GetUser(ctx, payment.To)
		if err != nil {
			service.log.Error("error retrieving user ID associated with wallet", zap.String("Wallet Address", payment.To.Hex()))
			return nil, err
		}

		metadata, err := json.Marshal(storjscanMetadata{
			ReferenceID: payment.Transaction.Hex(),
			Wallet:      payment.To.Hex(),
			BlockNumber: payment.BlockNumber,
			LogIndex:    payment.LogIndex,
		})
		if err != nil {
			service.log.Error("error populating metadata", zap.String("Reference ID", payment.Transaction.Hex()), zap.String("Wallet", payment.To.Hex()))
			return nil, err
		}

		billingTXs = append(billingTXs, billing.Transaction{
			UserID:      userID,
			Amount:      payment.USDValue,
			Description: "Storj token deposit",
			Source:      service.Source(),
			Status:      payments.PaymentStatusConfirmed,
			Type:        service.Type(),
			Metadata:    metadata,
			Timestamp:   payment.Timestamp,
		})
	}
	return billingTXs, nil
}
