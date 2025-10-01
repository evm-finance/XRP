package xrp

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"qc-defi-graphql-server/internal/models"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxInput struct {
	MethodSigData    []byte
	InputsSigData    []byte
	MethodSigDataStr string
	InputsSigDataStr string
	FunctionName     string
	FullFunctionSig  string
	InputsMap        map[string]interface{}
	ArgsMap          map[string]interface{}
}

type EvmTransaction struct {
	Hash           string
	IsPending      bool
	Status         string
	ChainId        int
	Block          int64
	Timestamp      int64
	Value          float64
	From           string
	To             string
	GssLimit       float64
	GasPrice       float64
	TransactionFee float64
	Nonce          float64
	TxDataHex      string
	Input          *TxInput
	LogEvents      *models.TransactionEvents
}

type ParamQueue struct {
	items []map[string]interface{}
}

func (q *ParamQueue) IsEmpty() bool {
	return len(q.items) == 0
}
func (q *ParamQueue) Enqueue(item map[string]interface{}) {
	q.items = append(q.items, item)
}
func (q *ParamQueue) Dequeue() map[string]interface{} {
	if q.IsEmpty() {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func getPoolsFromSwapEvents(events []*models.TransactionEvent) []string {
	var addresses []string
	for _, e := range events {
		if e.Name == models.Swap {
			addresses = append(addresses, strings.ToLower(e.Address))
		}
	}
	return addresses
}

// GetTransactionSender returns the sender address of a transaction
func GetTransactionSender(tx *types.Transaction, chainID *big.Int) (common.Address, error) {
	signer := types.LatestSignerForChainID(chainID)
	return types.Sender(signer, tx)
}

type TxService struct {
	ABIs      map[models.ContractType]abi.ABI
	EthClient *ethclient.Client
}

func NewTxService(ethClient *ethclient.Client) *TxService {
	// For now, return empty ABI map - you'll need to implement SupportedABIs()
	return &TxService{ABIs: make(map[models.ContractType]abi.ABI), EthClient: ethClient}
}

func (svs *TxService) SetClient(ethClient *ethclient.Client) { svs.EthClient = ethClient }

func (svs *TxService) getEventFromLog(vLog *types.Log) *models.TransactionEvent {
	var (
		dexEvent      = models.TransactionEvent{}
		indexedParams = map[string]interface{}{}
		queue         = ParamQueue{}
	)

	if len(vLog.Topics) > 1 {
		for i, param := range vLog.Topics[1:] {
			p := map[string]interface{}{"id": i, "decoded": common.HexToAddress(param.Hex()), "dex": param.String()}
			indexedParams[strconv.Itoa(i)] = p
			queue.Enqueue(p)
		}
	}

	if len(vLog.Topics) > 0 {
		for k, abiVal := range svs.ABIs {
			event, err := abiVal.EventByID(vLog.Topics[0])
			if err != nil {
				continue
			}
			outputDataMap := make(map[string]interface{})
			if len(vLog.Data) > 0 {
				err = abiVal.UnpackIntoMap(outputDataMap, event.Name, vLog.Data)

				if err != nil {
					continue
				}
			}
			finalMap := map[string]interface{}{}
			for key, v := range event.Inputs {
				val := make(map[string]interface{})
				val["name"] = v.Name
				val["type"] = v.Type.String()
				val["indexed"] = v.Indexed
				if v.Indexed {
					elem := queue.Dequeue()
					if elem != nil {
						if value, ok := elem["decoded"]; ok {
							val["value"] = value
						}
					}
				} else {
					if value, ok := outputDataMap[v.Name]; ok {
						val["value"] = value
					}
				}
				finalMap[strconv.Itoa(key)] = val
			}

			dexEvent.Contract = k
			dexEvent.Topic = vLog.Topics[0].String()
			dexEvent.Name = models.EventType(event.Name)
			dexEvent.Address = vLog.Address.String()
			dexEvent.Signature = event.Sig
			dexEvent.OutputDataMap = outputDataMap
			dexEvent.AllFunctionParams = finalMap
			dexEvent.OutputDataMapHex = hex.EncodeToString(vLog.Data)
		}
	}

	dexEvent.IndexedParams = indexedParams
	return &dexEvent
}

func txStatus(receipt *types.Receipt) string {
	if receipt.Status == 1 {
		return "success"
	} else if receipt.Status == 0 {
		return "failed!"
	}
	return "pending"
}

func (svs *TxService) txTimestamp(receipt *types.Receipt) int64 {
	blockHeader, err := svs.EthClient.HeaderByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		return 0
	}
	return int64(blockHeader.Time)
}

func (svs *TxService) parseTransactionBaseInfo(tx *types.Transaction, receipt *types.Receipt) *EvmTransaction {
	gasPrice := tx.GasPrice()
	gasLimit := tx.Gas()
	transactionFee := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasLimit))

	// Get sender address
	sender, err := GetTransactionSender(tx, tx.ChainId())
	if err != nil {
		sender = common.Address{} // Use zero address if error
	}

	return &EvmTransaction{
		Hash:           tx.Hash().Hex(),
		ChainId:        int(tx.ChainId().Int64()),
		Block:          int64(receipt.BlockNumber.Uint64()),
		Timestamp:      svs.txTimestamp(receipt),
		IsPending:      false,
		Status:         txStatus(receipt),
		Value:          float64(tx.Value().Uint64()) / 1e18, // Convert from wei to ether
		From:           sender.Hex(),
		To:             tx.To().Hex(),
		GssLimit:       float64(gasLimit),
		GasPrice:       float64(gasPrice.Uint64()),
		TransactionFee: float64(transactionFee.Uint64()) / 1e18, // Convert from wei to ether
		Nonce:          float64(tx.Nonce()),
		TxDataHex:      hex.EncodeToString(tx.Data()),
	}
}

func (svs *TxService) DecodeTransactionInputData(data []byte) *TxInput {
	if len(data) < 4 {
		return &TxInput{MethodSigData: data}
	}
	var (
		methodSigData         = data[:4]
		inputsSigData         = data[4:]
		functionName          string
		fullFunctionSignature string
		inputsMap             = make(map[string]interface{})
		argsMap               = make(map[string]interface{})
	)

	for _, contractABI := range svs.ABIs {
		method, err := contractABI.MethodById(methodSigData)
		if err != nil {
			continue
		}

		if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
			continue
		}

		functionName = method.Name
		fullFunctionSignature = method.String()

		for i, v := range method.Inputs {
			m := make(map[string]interface{})
			m[v.Name] = v.Type.String()
			argsMap[strconv.Itoa(i)] = m
		}
	}

	return &TxInput{
		MethodSigData:    methodSigData,
		InputsSigData:    inputsSigData,
		MethodSigDataStr: hex.EncodeToString(methodSigData),
		InputsSigDataStr: hex.EncodeToString(inputsSigData),
		InputsMap:        inputsMap,
		ArgsMap:          argsMap,
		FunctionName:     functionName,
		FullFunctionSig:  fullFunctionSignature,
	}
}

func (svs *TxService) ProcessTXEvents(tx *types.Transaction, receipt *types.Receipt) []*models.TransactionEvent {
	var (
		events []*models.TransactionEvent
		wg     = &sync.WaitGroup{}
		mu     = &sync.Mutex{}
	)

	for _, vLog := range receipt.Logs {
		wg.Add(1)
		go func(v *types.Log) {
			defer wg.Done()
			e := svs.getEventFromLog(v)
			mu.Lock()
			events = append(events, e)
			mu.Unlock()
		}(vLog)
	}
	wg.Wait()
	return events
}

func (svs *TxService) GetTx(hash string) (*EvmTransaction, error) {
	ctx := context.Background()
	txHash := common.HexToHash(hash)
	tx, isPending, err := svs.EthClient.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	receipt, err := svs.EthClient.TransactionReceipt(ctx, tx.Hash())

	if err != nil {
		return nil, fmt.Errorf("can't get tx receipt. %s", err.Error())
	}

	info := svs.parseTransactionBaseInfo(tx, receipt)
	info.IsPending = isPending
	info.Input = svs.DecodeTransactionInputData(tx.Data())
	logEvents := svs.ProcessTXEvents(tx, receipt)
	info.LogEvents = &models.TransactionEvents{Items: logEvents}
	return info, nil
}
