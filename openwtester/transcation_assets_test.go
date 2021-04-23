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

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract, extParam map[string]interface{}) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "", contract, extParam)

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
	//log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer_BSC(t *testing.T) {

	addrs := []string{
		//"0x7321b4b889f2fc7c1735be59a684ef376d959290",
		//"0x850553e000faf6e68bb2db8d04a5beeca6830193",
		//"0xc817710c3a41f03c67d2a84bbafd59deffd79a06",
		//"0xc9e0e0e939052e15d7cd3c57c8bc3130ecc2f5f7",
		//"0xfb1eecbccb02566a12c488c42de2a9eed24edcef",
		//"0xfc2019d260c58315233602542e9c8bdf4ceebc26",

		"0x2e78442558a6f4f3aef1046ddebda1a1c33e0d08", //fee support
	}

	tm := testInitWalletManager()
	walletID := "WKXw6NaV1nV65AN6QkaRtKkoWQ8UtdCMPi"
	accountID := "5zvpByUfjDtYxE9Ws9MaoMtdw79iY3MNeM7KQZF4rUGF"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	for _, to := range addrs {
		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.03", "", nil, nil)
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
}

func TestTransfer_ERC20(t *testing.T) {

	addrs := []string{
		"0x7321b4b889f2fc7c1735be59a684ef376d959290",
		"0x850553e000faf6e68bb2db8d04a5beeca6830193",
		"0xc817710c3a41f03c67d2a84bbafd59deffd79a06",
		"0xc9e0e0e939052e15d7cd3c57c8bc3130ecc2f5f7",
		"0xfb1eecbccb02566a12c488c42de2a9eed24edcef",
		"0xfc2019d260c58315233602542e9c8bdf4ceebc26",
	}

	tm := testInitWalletManager()
	walletID := "WKXw6NaV1nV65AN6QkaRtKkoWQ8UtdCMPi"
	accountID := "5zvpByUfjDtYxE9Ws9MaoMtdw79iY3MNeM7KQZF4rUGF"

	contract := openwallet.SmartContract{
		Address:  "0x8076c74c5e3f5852037f31ff0093eeb8c8add8d3",
		Symbol:   "BSC",
		Name:     "SAFEMOON",
		Token:    "SAFEMOON",
		Decimals: 9,
	}

	testGetAssetsAccountBalance(tm, walletID, accountID)

	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	for _, to := range addrs {
		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "123", "", &contract, nil)
		if err != nil {
			return
		}

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
}

func TestSummary_BSC(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WKXw6NaV1nV65AN6QkaRtKkoWQ8UtdCMPi"
	accountID := "5LXocNDt6Q2RQ56RkcpLafZ8kUd3oTCdRfgCideKF2NQ"
	summaryAddress := "0x69aeba03354894e9136f60df0c4d2b3001ce1e3c"

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

func TestSummary_ERC20(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WKXw6NaV1nV65AN6QkaRtKkoWQ8UtdCMPi"
	accountID := "5LXocNDt6Q2RQ56RkcpLafZ8kUd3oTCdRfgCideKF2NQ"
	summaryAddress := "0x69aeba03354894e9136f60df0c4d2b3001ce1e3c"

	feesSupport := openwallet.FeesSupportAccount{
		AccountID: "C9NaeqdMvi4ufvGiiTnkEcBG96tVAfhiPtiKqJtzX8Jv",
		//FixSupportAmount: "0.01",
		FeesSupportScale: "1.5",
	}

	contract := openwallet.SmartContract{
		Address:  "0x8076c74c5e3f5852037f31ff0093eeb8c8add8d3",
		Symbol:   "BSC",
		Name:     "SAFEMOON",
		Token:    "SAFEMOON",
		Decimals: 9,
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
