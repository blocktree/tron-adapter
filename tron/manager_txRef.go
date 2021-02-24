/*
* Copyright 2018 The openwallet Authors
* This file is part of the openwallet library.
*
* The openwallet library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The openwallet library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
 */

package tron

import (
	"encoding/hex"
	"fmt"
	"github.com/blocktree/openwallet/v2/common"
	"github.com/blocktree/openwallet/v2/openwallet"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/tron-adapter/tron/grpc-gateway/core"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
)

type AddrBalance struct {
	Address      string
	TronBalance  *big.Int
	TokenBalance *big.Int
	Index        int
}

func getTxHash1(txHex string) ([]byte, error) {
	tx := &core.Transaction{}
	/*
		 if txRawBts, err := hex.DecodeString(tx.GetRawData()); err != nil {
			 return nil, err
		 } else {
			 if err := proto.Unmarshal(txRawBts, tx); err != nil {
				 return signedTxRaw, err
			 }
		 }
	*/
	txByte, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, fmt.Errorf("get Tx hex failed;unexpected error: %v", err)
	}
	if err := proto.Unmarshal(txByte, tx); err != nil {
		return nil, fmt.Errorf("unmarshal RawData failed; unexpected error: %v", err)
	}
	txRaw, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("marshal RawData failed;unexpected error:%v", err)
	}
	txHash := owcrypt.Hash(txRaw, 0, owcrypt.HASH_ALG_SHA256)
	return txHash, nil
}

func getTxHash(tx *core.Transaction) (txHash []byte, err error) {

	txRaw, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("marshal RawData failed; unexpected error: %v", err)
	}
	txHash = owcrypt.Hash(txRaw, 0, owcrypt.HASH_ALG_SHA256)
	return txHash, nil
}

//Deprecated
func (wm *WalletManager) CreateTransactionRef(toAddress, ownerAddress string, amount string) (txRawHex string, err error) {

	// addressEncoder.AddressDecode return 20 bytes of the center of Address
	toAddressBytes, err := addressEncoder.AddressDecode(toAddress, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		wm.Log.Info("toAddress decode failed failed;unexpected error:%v", err)
		return "", err
	}
	toAddressBytes = append([]byte{0x41}, toAddressBytes...)

	ownerAddressBytes, err := addressEncoder.AddressDecode(ownerAddress, addressEncoder.TRON_mainnetAddress)
	if err != nil {
		wm.Log.Info("ownerAddress decode failed failed;unexpected error:%v", err)
		return "", err
	}
	ownerAddressBytes = append([]byte{0x41}, ownerAddressBytes...)

	amountDec := common.StringNumToBigIntWithExp(amount, Decimals)

	// Check amount: amount * 1000000
	// ******** Generate TX Contract ********
	tc := &core.TransferContract{
		OwnerAddress: ownerAddressBytes,
		ToAddress:    toAddressBytes,
		Amount:       amountDec.Int64(),
	}

	tcRaw, err := proto.Marshal(tc)
	if err != nil {
		wm.Log.Info("marshal tc failed;unexpected error:%v", err)
		return "", err
	}

	txContact := &core.Transaction_Contract{
		Type:         core.Transaction_Contract_TransferContract,
		Parameter:    &any.Any{Value: tcRaw, TypeUrl: "type.googleapis.com/protocol.TransferContract"},
		Provider:     nil,
		ContractName: nil,
	}

	// ******** Get Reference Block ********
	block, err := wm.GetNowBlock()
	if err != nil {
		wm.Log.Info("get current block failed;unexpected error:%v", err)
		return "", err
	}
	blockID, err := hex.DecodeString(block.GetBlockHashID())
	if err != nil {
		wm.Log.Info("conver BlockHashID from hex to byte failed;unexpected error:%v", err)
		return txRawHex, err
	}
	refBlockBytes, refBlockHash := blockID[6:8], blockID[8:16]

	// ********Set timestamp ********
	/*
	 According to RFC-3339 date strings
	 Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000).setNanos((int) ((millis % 1000) * 1000000)).build();
	*/
	timestamp := time.Now().UnixNano() / 1000000 // <int64

	// ******** Create Traction ********
	txRaw := &core.TransactionRaw{
		RefBlockBytes: refBlockBytes,
		RefBlockHash:  refBlockHash,
		Contract:      []*core.Transaction_Contract{txContact},
		Expiration:    timestamp + 10*60*60*60,
		// Timestamp:     timestamp,
	}
	tx := &core.Transaction{
		RawData: txRaw,
		// Signature: nil,
		// Ret:       nil,
	}

	// ******** TX Encoding ********
	if x, err := proto.Marshal(tx); err != nil {
		wm.Log.Info("marshal tx failed;unexpected error:%v", err)
		return "", err
	} else {
		txRawHex = hex.EncodeToString(x)
	}
	return txRawHex, nil
}

func (wm *WalletManager) CreateTokenTransaction(toAddress, ownerAddress string, amount string, contract openwallet.SmartContract) (txRawHex string, err error) {

	// addressEncoder.AddressDecode return 20 bytes of the center of Address
	toAddrHex, toAddressBytes, err := DecodeAddress(toAddress, wm.Config.IsTestNet)
	if err != nil {
		wm.Log.Info("toAddress decode failed failed;unexpected error:%v", err)
		return "", err
	}

	_, ownerAddressBytes, err := DecodeAddress(ownerAddress, wm.Config.IsTestNet)
	if err != nil {
		wm.Log.Info("ownerAddress decode failed failed;unexpected error:%v", err)
		return "", err
	}

	if contract.Address == "" {
		amountDec := common.StringNumToBigIntWithExp(amount, wm.Decimal())

		tc := &core.TransferContract{
			OwnerAddress: ownerAddressBytes,
			ToAddress:    toAddressBytes,
			Amount:       amountDec.Int64(),
		}
		return wm.createAssetsTransaction(tc, core.Transaction_Contract_TransferContract, TransferContract)

	} else if strings.EqualFold(contract.Protocol, TRC10) {
		amountDec := common.StringNumToBigIntWithExp(amount, int32(contract.Decimals))
		tc := &core.TransferAssetContract{
			OwnerAddress: ownerAddressBytes,
			ToAddress:    toAddressBytes,
			Amount:       amountDec.Int64(),
			AssetName:    []byte(contract.Address),
		}

		return wm.createAssetsTransaction(tc, core.Transaction_Contract_TransferAssetContract, TransferAssetContract)
	} else if strings.EqualFold(contract.Protocol, TRC20) {
		amountDec := common.StringNumToBigIntWithExp(amount, int32(contract.Decimals))
		_, contractAddressBytes, err := DecodeAddress(contract.Address, wm.Config.IsTestNet)
		if err != nil {
			return "", err
		}

		var funcParams []SolidityParam
		funcParams = append(funcParams, SolidityParam{
			ParamType:  SOLIDITY_TYPE_ADDRESS,
			ParamValue: toAddrHex,
		})

		funcParams = append(funcParams, SolidityParam{
			ParamType:  SOLIDITY_TYPE_UINT256,
			ParamValue: amountDec,
		})

		//fmt.Println("make token transfer data, amount:", amount.String())
		dataHex, err := makeTransactionParameter(TRC20_TRANSFER_METHOD_ID, funcParams)
		if err != nil {
			return "", err
		}

		data, err := hex.DecodeString(dataHex)
		if err != nil {
			return "", err
		}

		tc := &core.TriggerSmartContract{
			OwnerAddress:    ownerAddressBytes,
			ContractAddress: contractAddressBytes,
			Data:            data,
		}
		return wm.createAssetsTransaction(tc, core.Transaction_Contract_TriggerSmartContract, TriggerSmartContract)
	}

	return "", fmt.Errorf("%s is not supported", contract.Protocol)
}

func (wm *WalletManager) createAssetsTransaction(message proto.Message, contractType core.Transaction_Contract_ContractType, typeName string) (string, error) {

	msgBytes, err := proto.Marshal(message)
	if err != nil {
		return "", err
	}

	txContact := &core.Transaction_Contract{
		Type:         contractType,
		Parameter:    &any.Any{Value: msgBytes, TypeUrl: "type.googleapis.com/protocol." + typeName},
		Provider:     nil,
		ContractName: nil,
	}

	// ******** Get Reference Block ********
	block, err := wm.GetNowBlock()
	if err != nil {
		wm.Log.Info("get current block failed;unexpected error:%v", err)
		return "", err
	}
	blockID, err := hex.DecodeString(block.GetBlockHashID())
	if err != nil {
		return "", err
	}
	refBlockBytes, refBlockHash := blockID[6:8], blockID[8:16]

	// ********Set timestamp ********
	/*
	 According to RFC-3339 date strings
	 Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000).setNanos((int) ((millis % 1000) * 1000000)).build();
	*/
	timestamp := time.Now().UnixNano() / 1000000 // <int64

	// ******** Create Traction ********
	txRaw := &core.TransactionRaw{
		RefBlockBytes: refBlockBytes,
		RefBlockHash:  refBlockHash,
		Contract:      []*core.Transaction_Contract{txContact},
		Expiration:    timestamp + 10*60*60*60,
		//FeeLimit:      FeeLimit,
	}

	if wm.Config.FeeLimit > 0 {
		txRaw.FeeLimit = wm.Config.FeeLimit
	}

	tx := &core.Transaction{
		RawData: txRaw,
	}

	x, err := proto.Marshal(tx)
	if err != nil {
		wm.Log.Info("marshal tx failed;unexpected error:%v", err)
		return "", err
	}
	txRawHex := hex.EncodeToString(x)
	return txRawHex, nil

}

func (wm *WalletManager) SignTransactionRef(hash string, privateKey string) (signedTxRaw string, err error) {

	txHash, err := hex.DecodeString(hash)

	if err != nil {
		wm.Log.Info("conver hash from hex to byte failed;unexpected error:%v", err)
		return "", err
	}
	pk, err := hex.DecodeString(privateKey)
	if err != nil {
		wm.Log.Info("conver privatekey from hex to byte failed;unexpected error:%v", err)
		return "", err
	}
	//sign, ret := signatureSet.TronSignature(pk, txHash)
	//if ret != owcrypt.SUCCESS {
	//	wm.Log.Infof("sign txHash failed;unexpected error:%v", ret)
	//	return "", fmt.Errorf("sign txHash failed")
	//}

	signature, v, sigErr := owcrypt.Signature(pk, nil, txHash, owcrypt.ECC_CURVE_SECP256K1)
	if sigErr != owcrypt.SUCCESS {
		return "", fmt.Errorf("transaction hash sign failed")
	}
	signature = append(signature, v)

	return hex.EncodeToString(signature), nil
}

func (wm *WalletManager) ValidSignedTokenTransaction(txHex string) error {

	tx := &core.Transaction{}
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		wm.Log.Info("conver txhex from hex to byte failed;unexpected error:%v", err)
		return err
	}
	if err := proto.Unmarshal(txBytes, tx); err != nil {
		wm.Log.Info("unmarshal txBytes failed;unexpected error:%v", err)
		return err
	}
	listContracts := tx.RawData.GetContract()
	countSignature := len(tx.Signature)
	txHash, err := getTxHash1(txHex)
	if err != nil {
		wm.Log.Info("get txHex hash failed;unexpected error:%v", err)
		return err
	}

	if countSignature == 0 {
		return fmt.Errorf("not found signature")
	}

	for i, contract := range listContracts {

		err = wm.validSignedTokenTransaction(txHash, tx.Signature[i], contract)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wm *WalletManager) validSignedTokenTransaction(txHash []byte, signature []byte, contract *core.Transaction_Contract) error {

	codeType := addressEncoder.TRON_mainnetAddress
	if wm.Config.IsTestNet {
		codeType = addressEncoder.TRON_testnetAddress
	}

	var ownerAddressHex string
	switch contract.Type {
	case core.Transaction_Contract_TransferContract:
		tc := &core.TransferContract{}
		err := proto.Unmarshal(contract.Parameter.GetValue(), tc)
		if err != nil {
			return err
		}
		ownerAddressHex = hex.EncodeToString(tc.GetOwnerAddress())
	case core.Transaction_Contract_TransferAssetContract:
		tc := &core.TransferAssetContract{}
		err := proto.Unmarshal(contract.Parameter.GetValue(), tc)
		if err != nil {
			return err
		}
		ownerAddressHex = hex.EncodeToString(tc.GetOwnerAddress())
	case core.Transaction_Contract_TriggerSmartContract:
		tc := &core.TriggerSmartContract{}
		err := proto.Unmarshal(contract.Parameter.GetValue(), tc)
		if err != nil {
			return err
		}
		ownerAddressHex = hex.EncodeToString(tc.GetOwnerAddress())
	}

	pkBytes, ret := owcrypt.RecoverPubkey(signature, txHash, wm.CurveType())

	if ret != owcrypt.SUCCESS {
		return fmt.Errorf("verify SignedTransactionRef faild: recover Pubkey error")
	}
	//if owcrypt.SUCCESS != owcrypt.Verify(pkBytes, nil, 0, txHash, 32, signature[0:len(signature)-1], wm.CurveType()) {
	//	return fmt.Errorf("verify SignedTransactionRef failed:verify signature failed")
	//}
	pkHash := owcrypt.Hash(pkBytes, 0, owcrypt.HASH_ALG_KECCAK256)[12:32]
	pkgenAddress := append(codeType.Prefix, pkHash...)
	pkgenAddressHex := hex.EncodeToString(pkgenAddress)
	if pkgenAddressHex != ownerAddressHex {
		return fmt.Errorf("verify SignedTransactionRef failed: signed address is not the owner address")
	}
	return nil
}

//Deprecated
func (wm *WalletManager) ValidSignedTransactionRef(txHex string) error {

	tx := &core.Transaction{}
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		wm.Log.Info("conver txhex from hex to byte failed;unexpected error:%v", err)
		return err
	}
	if err := proto.Unmarshal(txBytes, tx); err != nil {
		wm.Log.Info("unmarshal txBytes failed;unexpected error:%v", err)
		return err
	}
	listContracts := tx.RawData.GetContract()
	countSignature := len(tx.Signature)
	txHash, err := getTxHash1(txHex)
	if err != nil {
		wm.Log.Info("get txHex hash failed;unexpected error:%v", err)
		return err
	}

	if countSignature == 0 {
		return fmt.Errorf("not found signature")
	}

	for i := 0; i < countSignature; i++ {
		contract := listContracts[i]

		// Get the instance of TransferContract to get Owner Address for validate signature
		tc := &core.TransferContract{}
		err := proto.Unmarshal(contract.Parameter.GetValue(), tc)
		if err != nil {
			wm.Log.Info("unmarshal contract (value) failed;unexpected error:%v", err)
			return err
		}
		ownerAddressHex := hex.EncodeToString(tc.GetOwnerAddress())
		pkBytes, ret := owcrypt.RecoverPubkey(tx.Signature[i], txHash, wm.CurveType())

		if ret != owcrypt.SUCCESS {
			return fmt.Errorf("verify SignedTransactionRef faild: recover Pubkey error")
		}
		//if owcrypt.SUCCESS != owcrypt.Verify(pkBytes, nil, 0, txHash, 32, tx.Signature[i][0:len(tx.Signature[i])-1], wm.CurveType()) {
		//	return fmt.Errorf("verify SignedTransactionRef failed:verify signature failed")
		//}
		pkHash := owcrypt.Hash(pkBytes, 0, owcrypt.HASH_ALG_KECCAK256)[12:32]
		pkgenAddress := append([]byte{0x41}, pkHash...)
		pkgenAddressHex := hex.EncodeToString(pkgenAddress)
		if pkgenAddressHex != ownerAddressHex {
			return fmt.Errorf("verify SignedTransactionRef failed: signed address is not the owner address")
		}
	}
	return nil
}

func (wm *WalletManager) BroadcastTransaction(raw string) (string, error) {
	tx := &core.Transaction{}
	if txBytes, err := hex.DecodeString(raw); err != nil {
		wm.Log.Info("conver raw from hex to byte failed;unexpected error:%v", err)
		return "", err
	} else {
		if err := proto.Unmarshal(txBytes, tx); err != nil {
			wm.Log.Info("unmarshal txBytes failed;unexpected error:%v", err)
			return "", err
		}
	}

	/* Generate Params */
	var (
		signature []string
		txID      string
		contracts []map[string]interface{}
		raw_data  map[string]interface{}
	)
	for _, x := range tx.GetSignature() {
		signature = append(signature, hex.EncodeToString(x)) // base64
	}
	if txHash, err := getTxHash1(raw); err != nil {
		wm.Log.Info("get raw hash failed;unexpected error:%v", err)
		return "", err
	} else {
		txID = hex.EncodeToString(txHash)
	}
	rawData := tx.GetRawData()
	contracts = []map[string]interface{}{}
	for _, c := range rawData.GetContract() {
		any := c.GetParameter().GetValue()
		var contract map[string]interface{}
		switch c.Type {
		case core.Transaction_Contract_TransferContract:
			tc := &core.TransferContract{}
			if err := proto.Unmarshal(any, tc); err != nil {
				wm.Log.Info("unmarshal contract value failed;unexpected error:%v", err)
				return "", err
			}

			contract = map[string]interface{}{
				"type": c.GetType().String(),
				"parameter": map[string]interface{}{
					"type_url": c.GetParameter().GetTypeUrl(),
					"value": map[string]interface{}{
						"amount":        tc.Amount,
						"owner_address": hex.EncodeToString(tc.GetOwnerAddress()),
						"to_address":    hex.EncodeToString(tc.GetToAddress()),
					},
				},
			}

		case core.Transaction_Contract_TransferAssetContract:
			tc := &core.TransferAssetContract{}
			if err := proto.Unmarshal(any, tc); err != nil {
				wm.Log.Info("unmarshal contract value failed;unexpected error:%v", err)
				return "", err
			}

			contract = map[string]interface{}{
				"type": c.GetType().String(),
				"parameter": map[string]interface{}{
					"type_url": c.GetParameter().GetTypeUrl(),
					"value": map[string]interface{}{
						"amount":        tc.Amount,
						"owner_address": hex.EncodeToString(tc.GetOwnerAddress()),
						"to_address":    hex.EncodeToString(tc.GetToAddress()),
						"asset_name":    hex.EncodeToString(tc.AssetName),
					},
				},
			}

		case core.Transaction_Contract_TriggerSmartContract:
			tc := &core.TriggerSmartContract{}
			if err := proto.Unmarshal(any, tc); err != nil {
				wm.Log.Info("unmarshal contract value failed;unexpected error:%v", err)
				return "", err
			}

			contract = map[string]interface{}{
				"type": c.GetType().String(),
				"parameter": map[string]interface{}{
					"type_url": c.GetParameter().GetTypeUrl(),
					"value": map[string]interface{}{
						"data":             hex.EncodeToString(tc.Data),
						"owner_address":    hex.EncodeToString(tc.GetOwnerAddress()),
						"contract_address": hex.EncodeToString(tc.GetContractAddress()),
					},
				},
			}
		}

		contracts = append(contracts, contract)
	}
	raw_data = map[string]interface{}{
		"ref_block_bytes": hex.EncodeToString(rawData.GetRefBlockBytes()),
		"ref_block_hash":  hex.EncodeToString(rawData.GetRefBlockHash()),
		"expiration":      rawData.GetExpiration(),
		"timestamp":       rawData.GetTimestamp(),
		"contract":        contracts,
	}

	if wm.Config.FeeLimit > 0 {
		raw_data["fee_limit"] = wm.Config.FeeLimit
	}

	params := req.Param{
		"signature": signature,
		"txID":      txID,
		"raw_data":  raw_data,
	}
	// Call api to broadcast transaction
	r, err := wm.WalletClient.Call("/wallet/broadcasttransaction", params)

	if err != nil {
		wm.Log.Info("broadcast transaction failed;unexpected error:%v", err)
		return "", err
	} else {
		if r.Get("result").Bool() != true {
			var err error
			if r.Get("message").String() != "" {
				msg, _ := hex.DecodeString(r.Get("message").String())
				err = fmt.Errorf("BroadcastTransaction error message: %+v", string(msg))
			} else {
				err = fmt.Errorf("BroadcastTransaction return error: %+v", r)
			}
			return "", err
		}
	}
	return txID, nil
}

//SendTransaction 发送交易
func (wm *WalletManager) SendTransaction(walletID, to string, amount decimal.Decimal, password string, feesInSender bool) ([]string, error) {

	return nil, nil
}

//deprecated
func (wm *WalletManager) Getbalance(address string) (*AddrBalance, error) {
	account, err := wm.GetAccount(address)
	if err != nil {
		return nil, err
	}
	var balance int64
	if account.Balance == "" {
		balance = 0
	} else {
		balance, err = strconv.ParseInt(account.Balance, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	ret := &AddrBalance{TronBalance: big.NewInt(balance)}
	ret.Address = address
	return ret, nil
}

/*
// ------------------------------------------------------------------------------------------------------
func debugPrintTx(txRawhex string) {

	tx := &core.Transaction{}
	if txRawBts, err := hex.DecodeString(txRawhex); err != nil {
		fmt.Println(err)
	} else {
		if err := proto.Unmarshal(txRawBts, tx); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv Print Test vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")

	txHash, err := getTxHash(tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Tx Hash = ", hex.EncodeToString(txHash))

	txRawD := tx.RawData
	txC := txRawD.GetContract()
	fmt.Println("txRawD.Contract = ")
	for _, c := range txC {
		fmt.Println("\tc.ContractName=", c.ContractName)
		fmt.Println("\tc.Provider   =", c.Provider)
		fmt.Println("\tc.Type       =", c.Type)
		fmt.Println("\tc.Parameter   =", c.Parameter)

		ts := &core.TransferContract{}
		proto.Unmarshal(c.Parameter.Value, ts)
		fmt.Println("\tts.OwnerAddress =", hex.EncodeToString(ts.OwnerAddress))
		fmt.Println("\tts.ToAddress =", hex.EncodeToString(ts.ToAddress))
		fmt.Println("\tts.Amount =", ts.Amount)
	}
	fmt.Println("txRawD.Data =  ", txRawD.Data)
	fmt.Println("txRawD.Auths =   ", txRawD.Auths)
	fmt.Println("txRawD.Scripts =   ", txRawD.Scripts)
	fmt.Println("txRawD.RefBlockBytes = ", hex.EncodeToString(txRawD.RefBlockBytes))
	fmt.Println("txRawD.RefBlockHash Bts = ", txRawD.RefBlockHash, "Len:", len(txRawD.RefBlockHash))
	fmt.Println("txRawD.RefBlockHash Hex = ", hex.EncodeToString(txRawD.RefBlockHash), "Len:", len(hex.EncodeToString(txRawD.RefBlockHash)))
	// dst := make([]byte, 32)
	// bs, err := base64.StdEncoding.Decode(dst, txRawD.RefBlockHash)
	// fmt.Println("txRawD.RefBlockHash base64Bytes = ", bs, "XX = ", dst)

	fmt.Println("")

	fmt.Println("txRawD.RefBlockNum =  ", txRawD.RefBlockNum)
	fmt.Println("txRawD.Expiration =  ", txRawD.Expiration)
	fmt.Println("txRawD.Timestamp =   ", txRawD.Timestamp)
	fmt.Println("tx.Signature[0]     = ", hex.EncodeToString(tx.Signature[0]))
	fmt.Println("tx.Ret          =     ", tx.Ret)

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ Print Test ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ End")
}
*/

//GetTransactionFeeEstimated 预计手续费
func (wm *WalletManager) GetTransactionFeeEstimated(from string, data string) (*txFeeInfo, error) {

	//:计算矿工费
	feeInfo := &txFeeInfo{
		GasUsed:  0,
		GasPrice: decimal.Zero,
		Fee:      decimal.Zero,
	}

	//:检查地址账户可用带宽是否足够
	/*
		accountNet, err := wm.GetAccountNet(from)
		if err != nil {
			return nil, err
		}
		log.Debugf("length of data: %d", len(data)/2)
		//先判断冻结的带宽是否可用
		if accountNet.NetUsed >= accountNet.NetLimit {
			//再判断免费的带宽是否可用
			if accountNet.FreeNetUsed >= accountNet.FreeNetLimit {
				//:矿工费 = 字节长度 * 10 SUN
				feeInfo.Fee = decimal.New(GasPrice, 0).Shift(-Decimals)
				feeInfo.GasUsed = int64(len(data) / 2)
				feeInfo.CalcFee()
			}
		}
	*/
	return feeInfo, nil
}

//IsEnoughEnergyToTransferTRC20 是否足够能量转账TRC20
func (wm *WalletManager) IsEnoughEnergyToTransferTRC20(address string, trxBalance *big.Int) (flag bool, energyRest int64, feeMini int64) {
	feeMini = wm.Config.FeeMini
	res, err := wm.GetAccountResource(address)
	if err != nil {
		return false, 0, feeMini
	}

	energyRest = res.EnergyLimit - res.EnergyUsed
	//1 Energy = 140 SUN
	trxEnergy := trxBalance.Div(trxBalance, big.NewInt(140)).Int64()
	energyRest = energyRest + trxEnergy
	//能量少于下限
	if energyRest < feeMini {
		return false, energyRest, feeMini
	}

	return true, energyRest, feeMini
}
