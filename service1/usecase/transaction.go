package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"service1/pkg/types"
)

func (u *Usecase) Invoice(transactions *types.Transactions) error {
	msg, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	err = u.kafka.SendMessage(msg, "Transactions")
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) GetBalance(ctx context.Context, id string) (*types.Wallet, error) {
	wallet, err := u.repository.GetWalletBy(ctx, map[string]any{"id": id})
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (u *Usecase) CreateTransaction(ctx context.Context, transactions *types.Transactions) (*string, *string, error) {
	if transactions.Amount < 0 {
		wallet, err := u.repository.GetWalletBy(ctx, map[string]any{"id": transactions.IDSender})
		if err != nil {
			return nil, nil, err
		}
		actual := wallet.Actual
		if actual+transactions.Amount < 0 {
			return nil, nil, fmt.Errorf("insufficient funds")
		}
		frozen := wallet.Frozen

		actual += transactions.Amount

		err = u.repository.UpdateWallet(ctx, actual, frozen, map[string]any{"id": transactions.IDSender})
		if err != nil {
			return nil, nil, err
		}
	}

	wallet, err := u.repository.GetWalletBy(ctx, map[string]any{"id": transactions.Wallet})
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	actual := wallet.Actual
	frozen := wallet.Frozen

	frozen += int(math.Abs(float64(transactions.Amount)))
	err = u.repository.UpdateWallet(ctx, actual, frozen, map[string]any{"id": transactions.Wallet})
	if err != nil {
		return nil, nil, err
	}
	id, status, err := u.repository.CreateTransaction(ctx, transactions)
	if err != nil {
		return nil, nil, err
	}
	transactions.ID = *id
	msg, err := json.Marshal(transactions)
	if err != nil {
		return nil, nil, err
	}

	err = u.kafka.SendMessage(msg, "Transactions")
	if err != nil {
		return nil, nil, err
	}

	return id, status, nil
}

func (u *Usecase) ProcessingTransaction(ctx context.Context, transactions *types.Transactions) error {
	wallet, err := u.repository.GetWalletBy(ctx, map[string]any{"id": transactions.Wallet})
	if err != nil {
		return err
	}
	actual := wallet.Actual
	frozen := wallet.Frozen

	if transactions.Amount < 0 {
		actual -= transactions.Amount
		frozen += transactions.Amount
	} else {
		actual += transactions.Amount
		frozen -= transactions.Amount
	}

	err = u.repository.UpdateWallet(ctx, actual, frozen, map[string]any{"id": transactions.Wallet})
	if err != nil {
		return err
	}

	err = u.repository.UpdateTransaction(ctx, "Success", map[string]any{"id": transactions.ID})
	if err != nil {
		return err
	}

	return nil
}
