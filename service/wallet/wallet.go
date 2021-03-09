package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/shopspring/decimal"
)

type Config struct {
	Pin       string   `valid:"required"`
	Members   []string `valid:"required"`
	Threshold uint8    `valid:"required"`
}

func New(client *mixin.Client, cfg Config) core.WalletService {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	return &walletService{
		client:    client,
		members:   cfg.Members,
		threshold: cfg.Threshold,
		pin:       cfg.Pin,
	}
}

type walletService struct {
	client    *mixin.Client
	members   []string
	threshold uint8
	pin       string
}

func (s *walletService) Pull(ctx context.Context, offset time.Time, limit int) ([]*core.Output, error) {
	outputs, err := s.client.ReadMultisigOutputs(ctx, s.members, s.threshold, offset, limit)
	if err != nil {
		return nil, err
	}

	results := make([]*core.Output, 0, len(outputs))
	for _, output := range outputs {
		result := convertToOutput(output)
		results = append(results, result)
	}

	return results, nil
}

// Spend 消费指定的 UTXO
// 如果 transfer 是 nil，则合并这些 UTXO
func (s *walletService) Spend(ctx context.Context, outputs []*core.Output, transfer *core.Transfer) (*core.RawTransaction, error) {
	state, tx, err := s.signTransaction(ctx, outputs, transfer)
	if err != nil {
		return nil, err
	}

	switch state {
	case mixin.UTXOStateSpent:
		// 已经付款出去了，说明 m 节点里面有 n 个节点签名，Output 已经全部花出去了。此处做简单验证:
		//	1. memo 与预期的一致
		//	2. 金额与预期的一致 (第一个 output 为支出，第二个 output 为找零)
		//	3. 第一个 output 为普通转账，且转给单个地址
		//	4. 如果有第二个 output , 则其应该也是 n/m 的多签
		//
		// 如果检验不通过，则说明本地数据可能与其他节点不一致，或者其他节点串通作恶等，直接报错等待人工确认处理。
		tx, err := mixin.TransactionFromRaw(tx)
		if err != nil {
			return nil, err
		}

		if err := s.validateTransaction(tx, transfer); err != nil {
			return nil, fmt.Errorf("validateTransaction failed: %w", err)
		}
	case mixin.UTXOStateSigned:
		sig, err := s.client.CreateMultisig(ctx, mixin.MultisigActionSign, tx)
		if err != nil {
			return nil, fmt.Errorf("CreateMultisig %s failed: %w", mixin.MultisigActionSign, err)
		}

		if !govalidator.IsIn(s.client.ClientID, sig.Signers...) {
			if valiErr := s.validateMultisig(sig, transfer); valiErr != nil {
				// unlock multisig
				// unlock, err := s.client.CreateMultisig(ctx, mixin.MultisigActionUnlock, tx)
				// if err != nil {
				// 	return nil, fmt.Errorf("CreateMultisig %s failed: %w", mixin.MultisigActionUnlock, err)
				// }
				//
				// if err := s.client.UnlockMultisig(ctx, unlock.RequestID, s.pin); err != nil {
				// 	return nil, fmt.Errorf("UnlockMultisig failed: %w", err)
				// }

				// 消费失败
				return nil, valiErr
			}

			sig, err = s.client.SignMultisig(ctx, sig.RequestID, s.pin)
			if err != nil {
				return nil, fmt.Errorf("SignMultisig failed: %w", err)
			}
		}

		// 签名数量达到要求，返回 raw transaction，将异步提交到主网
		if len(sig.Signers) >= int(sig.Threshold) {
			tx, err := mixin.TransactionFromRaw(sig.RawTransaction)
			if err != nil {
				return nil, err
			}

			if len(tx.Signatures) == 0 {
				return nil, fmt.Errorf("generate raw transaction failed, invalid signatures")
			}

			return &core.RawTransaction{
				TraceID: transfer.TraceID,
				Data:    sig.RawTransaction,
			}, nil
		}
	default:
		// 理论上程序逻辑不会走到这里
		return nil, errors.New("cannot consume unsigned utxo")
	}

	return nil, nil
}

func (s *walletService) ReqTransfer(ctx context.Context, transfer *core.Transfer) (string, error) {
	input := mixin.TransferInput{
		AssetID: transfer.AssetID,
		Amount:  transfer.Amount,
		TraceID: transfer.TraceID,
		Memo:    transfer.Memo,
	}

	input.OpponentMultisig.Receivers = transfer.Opponents
	input.OpponentMultisig.Threshold = transfer.Threshold

	payment, err := s.client.VerifyPayment(ctx, input)
	if err != nil {
		return "", err
	}

	return payment.CodeID, nil
}

func (s *walletService) HandleTransfer(ctx context.Context, transfer *core.Transfer) error {
	input := mixin.TransferInput{
		AssetID: transfer.AssetID,
		Amount:  transfer.Amount,
		TraceID: transfer.TraceID,
		Memo:    transfer.Memo,
	}

	input.OpponentMultisig.Receivers = transfer.Opponents
	input.OpponentMultisig.Threshold = transfer.Threshold

	_, err := s.client.Transaction(ctx, &input, s.pin)
	return err
}

// signTransaction 根据输入的 Output 计算出 Transaction Hash
func (s *walletService) signTransaction(ctx context.Context, outputs []*core.Output, transfer *core.Transfer) (string, string, error) {
	if len(outputs) == 0 {
		return mixin.UTXOStateSpent, "", nil
	}

	input := &mixin.TransactionInput{
		Memo: transfer.Memo,
		Hint: transfer.TraceID,
	}

	state := outputs[0].State
	signedTx := outputs[0].SignedTx
	sum := decimal.Zero

	for _, output := range outputs[0:] {
		st := output.State
		tx := output.SignedTx
		sum = sum.Add(output.Amount)

		if st == state && tx == signedTx {
			input.AppendUTXO(convertToUTXO(output, s.members, s.threshold))
			continue
		}

		return "", "", errors.New("state not match")
	}

	if signedTx != "" {
		return state, signedTx, nil
	}

	input.AppendOutput(transfer.Opponents, transfer.Threshold, transfer.Amount)
	tx, err := s.client.MakeMultisigTransaction(ctx, input)
	if err != nil {
		return "", "", err
	}

	signedTx, _ = tx.DumpTransaction()
	return mixin.UTXOStateSigned, signedTx, nil
}

// validateMultisig validate multisig request
func (s *walletService) validateMultisig(req *mixin.MultisigRequest, transfer *core.Transfer) error {
	if req.AssetID != transfer.AssetID {
		return fmt.Errorf("asset id not match, expect %q got %q", transfer.AssetID, req.AssetID)
	}

	if req.Memo != transfer.Memo {
		return fmt.Errorf("memo not match, expect %q got %q", transfer.Memo, req.Memo)
	}

	if !req.Amount.Equal(transfer.Amount) {
		return fmt.Errorf("amount not match, expect %s got %s", transfer.Amount, req.Amount)
	}

	if mixin.HashMembers(req.Receivers) != mixin.HashMembers(transfer.Opponents) {
		return errors.New("receivers not match")
	}

	return nil
}

// validateTransaction validate spent Tx
func (s *walletService) validateTransaction(tx *mixin.Transaction, transfer *core.Transfer) error {
	if string(tx.Extra) != transfer.Memo {
		return fmt.Errorf("memo not match, expect %q got %q", transfer.Memo, string(tx.Extra))
	}

	for idx, output := range tx.Outputs {
		switch idx {
		case 0: // 检查 output 和 transfer
			if output.Type != 0 {
				return fmt.Errorf("first output type not matched, expect %d got %d", 0, output.Type)
			}

			if expect, got := mixin.NewIntegerFromDecimal(transfer.Amount).String(), output.Amount.String(); expect != got {
				return fmt.Errorf("amount not match, expect %s got %s", expect, got)
			}

			if expect, got := mixin.NewThresholdScript(transfer.Threshold).String(), output.Script.String(); expect != got {
				return fmt.Errorf("first output script not matched, expect %s got %s", expect, got)
			}

			if len(output.Keys) != len(transfer.Opponents) {
				return errors.New("receivers not match")
			}
		default: // 检查找零
			if expect, got := mixin.NewThresholdScript(s.threshold).String(), output.Script.String(); expect != got {
				return fmt.Errorf("first output script not matched, expect %s got %s", expect, got)
			}

			if len(output.Keys) != len(s.members) {
				return errors.New("receivers not match")
			}
		}
	}

	return nil
}
