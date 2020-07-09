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

package openwtester

import (
	"github.com/blocktree/openwallet/v2/openw"
	"testing"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

func testGetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) {
	balance, err := tm.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return
	}
	log.Info("balance:", balance)
}

func testGetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) {
	balance, err := tm.GetAssetsAccountTokenBalance(testApp, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return
	}
	log.Info("token balance:", balance.Balance)
}

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "", contract)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testCreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract,
	feeSupportAccount *openwallet.FeesSupportAccount) ([]*openwallet.RawTransactionWithError, error) {

	rawTxArray, err := tm.CreateSummaryRawTransactionWithError(testApp, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract, feeSupportAccount)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

func testSignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer_TRX(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WLHdqGtGGZkBHEyXmv1w82s2iZjWJjgWF8"
	accountID := "4pF3jRC2XokaaLZWiiLvxXrD8SKRYNuzcVCFkJdu6rkt"
	to := "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	//to := "TCa3csiJd8Xhx75GPWP9S3kyXX9PonMx7n"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.1", "", nil)
	if err != nil {
		return
	}

	log.Std.Info("rawTx: %+v", rawTx)

	_, err = testSignTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testSubmitTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

}

func TestTransfer_TRC20(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WLHdqGtGGZkBHEyXmv1w82s2iZjWJjgWF8"
	accountID := "4pF3jRC2XokaaLZWiiLvxXrD8SKRYNuzcVCFkJdu6rkt"
	//to := "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	to := "TCa3csiJd8Xhx75GPWP9S3kyXX9PonMx7n"

	contract := openwallet.SmartContract{
		Address:  "THvZvKPLHKLJhEFYKiyqj6j8G8nGgfg7ur",
		Symbol:   "TRX",
		Name:     "TRONdice",
		Token:    "DICE",
		Decimals: 6,
		Protocol: "trc20",
	}

	//contract := openwallet.SmartContract{
	//	Address:  "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	//	Symbol:   "TRX",
	//	Name:     "Tether USD",
	//	Token:    "USDT",
	//	Decimals: 6,
	//	Protocol: "trc20",
	//}

	testGetAssetsAccountBalance(tm, walletID, accountID)

	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "1", "", &contract)
	if err != nil {
		return
	}
	//log.Infof("rawHex: %+v", rawTx.RawHex)
	_, err = testSignTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	//_, err = testSubmitTransactionStep(tm, rawTx)
	//if err != nil {
	//	return
	//}

}

func TestTransfer_TRC10(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WLHdqGtGGZkBHEyXmv1w82s2iZjWJjgWF8"
	accountID := "4pF3jRC2XokaaLZWiiLvxXrD8SKRYNuzcVCFkJdu6rkt"
	//to := "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	to := "TCa3csiJd8Xhx75GPWP9S3kyXX9PonMx7n"

	contract := openwallet.SmartContract{
		Address:  "1002000",
		Symbol:   "TRX",
		Name:     "BitTorrent",
		Token:    "BTT",
		Decimals: 6,
		Protocol: "trc10",
	}

	testGetAssetsAccountBalance(tm, walletID, accountID)

	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.001", "", &contract)
	if err != nil {
		return
	}
	log.Infof("rawHex: %+v", rawTx.RawHex)
	_, err = testSignTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	//_, err = testSubmitTransactionStep(tm, rawTx)
	//if err != nil {
	//	return
	//}

}

func TestSummary(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WGVsUfTTVaCwAMRTqeJiDQsZ3vrWp9DzMA"
	//accountID := "4pF3jRC2XokaaLZWiiLvxXrD8SKRYNuzcVCFkJdu6rkt"
	accountID := "C31rHUi8FJpwhWC2KTb5mMx9LUCSRpNnS1cG2QVMixYN"
	summaryAddress := "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	//summaryAddress := "TS11WZyPnT8qidwREcR8VDzNULCXxFeBMa"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, nil, nil)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}

func TestSummary_TRC10(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WGVsUfTTVaCwAMRTqeJiDQsZ3vrWp9DzMA"
	//accountID := "4pF3jRC2XokaaLZWiiLvxXrD8SKRYNuzcVCFkJdu6rkt"
	accountID := "C31rHUi8FJpwhWC2KTb5mMx9LUCSRpNnS1cG2QVMixYN"
	summaryAddress := "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	//summaryAddress := "TS11WZyPnT8qidwREcR8VDzNULCXxFeBMa"

	feesSupport := openwallet.FeesSupportAccount{
		AccountID:        "5Tm3sqFap329wj3Du4DVXMkjAe85FVH3MaB6HSV8joj1",
		FixSupportAmount: "1",
		//FeesSupportScale: "1.3",
	}

	contract := openwallet.SmartContract{
		Address:  "1002000",
		Symbol:   "TRX",
		Name:     "BitTorrent",
		Token:    "BTT",
		Decimals: 6,
		Protocol: "trc10",
	}

	testGetAssetsAccountBalance(tm, walletID, accountID)

	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, &contract, &feesSupport)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}

func TestSummary_TRC20(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WGVsUfTTVaCwAMRTqeJiDQsZ3vrWp9DzMA"
	//accountID := "4pF3jRC2XokaaLZWiiLvxXrD8SKRYNuzcVCFkJdu6rkt"
	accountID := "C31rHUi8FJpwhWC2KTb5mMx9LUCSRpNnS1cG2QVMixYN"
	summaryAddress := "TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5"
	//summaryAddress := "TS11WZyPnT8qidwREcR8VDzNULCXxFeBMa"

	feesSupport := openwallet.FeesSupportAccount{
		AccountID:        "5Tm3sqFap329wj3Du4DVXMkjAe85FVH3MaB6HSV8joj1",
		FixSupportAmount: "0.5",
		//FeesSupportScale: "1.3",
	}

	contract := openwallet.SmartContract{
		Address:  "THvZvKPLHKLJhEFYKiyqj6j8G8nGgfg7ur",
		Symbol:   "TRX",
		Name:     "TRONdice",
		Token:    "DICE",
		Decimals: 6,
		Protocol: "trc20",
	}

	//contract := openwallet.SmartContract{
	//	Address:  "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	//	Symbol:   "TRX",
	//	Name:     "Tether USD",
	//	Token:    "USDT",
	//	Decimals: 6,
	//	Protocol: "trc20",
	//}

	testGetAssetsAccountBalance(tm, walletID, accountID)

	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, &contract, &feesSupport)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}
