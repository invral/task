package service

import (
	"context"
	"errors"
	"task/internal/domain/transaction/entity"
	"task/internal/domain/transaction/service/mocks"
	"testing"
)

//type transaction struct {
//	id        int
//	accountID int
//	amount    int
//	currency  string
//	toAccount int
//}

func TestService_CreateDepositTransaction(t *testing.T) {
	//mock.NewRepository_acc_dto(t)
	//mock.NewRepository_transaction(t)

	cases := []struct {
		name      string
		ex        entity.Transaction
		wantError bool
	}{
		{
			name: "Invalid currency",
			ex: entity.Transaction{
				ID:        1,
				Status:    "created",
				AccountID: 1,
				Amount:    100,
				Currency:  "RUR",
				ToAccount: 0,
			},
			wantError: true,
		},
		{
			name: "Succes",
			ex: entity.Transaction{
				ID:        2,
				Status:    "created",
				AccountID: 1,
				Amount:    100,
				Currency:  "RUB",
				ToAccount: 0,
			},
			wantError: false,
		},
		{
			name: "No needed field toAccount",
			ex: entity.Transaction{
				ID:        2,
				Status:    "created",
				AccountID: 1,
				Amount:    100,
				Currency:  "USD",
				ToAccount: 1,
			},
			wantError: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			ctx := context.Background()
			repTr := mocks.NewRepository_transaction(t)

			if !tc.wantError {
				repTr.On("CreateDepositTransaction", ctx, &tc.ex).Return(nil)
			} else {
				repTr.On("CreateDepositTransaction", ctx, &tc.ex).
					Return(errors.New("Failed to create transaction"))
			}

			err := repTr.CreateDepositTransaction(ctx, &tc.ex)

			if (err != nil) != tc.wantError {
				t.Errorf("CreateDepositTransaction error = %v, wantErr %v", err, tc.wantError)
				return
			}

			//require.Error(t, err)

		})

	}

}
