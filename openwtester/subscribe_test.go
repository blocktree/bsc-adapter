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
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/v2/common/file"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openw"
	"github.com/blocktree/openwallet/v2/openwallet"
	"path/filepath"
	"testing"
)

////////////////////////// 测试单个扫描器 //////////////////////////

type subscriberSingle struct {
}

//BlockScanNotify 新区块扫描完成通知
func (sub *subscriberSingle) BlockScanNotify(header *openwallet.BlockHeader) error {
	log.Notice("header:", header)
	return nil
}

//BlockTxExtractDataNotify 区块提取结果通知
func (sub *subscriberSingle) BlockExtractDataNotify(sourceKey string, data *openwallet.TxExtractData) error {
	log.Notice("account:", sourceKey)

	for i, input := range data.TxInputs {
		log.Std.Notice("data.TxInputs[%d]: %+v", i, input)
	}

	for i, output := range data.TxOutputs {
		log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
	}

	log.Std.Notice("data.Transaction: %+v", data.Transaction)

	return nil
}

//BlockExtractSmartContractDataNotify 区块提取智能合约交易结果通知
func (sub *subscriberSingle) BlockExtractSmartContractDataNotify(sourceKey string, data *openwallet.SmartContractReceipt) error {

	log.Notice("sourceKey:", sourceKey)
	log.Std.Notice("data.ContractTransaction: %+v", data)

	for i, event := range data.Events {
		log.Std.Notice("data.Events[%d]: %+v", i, event)
	}

	return nil
}

func TestSubscribeAddress_BSC(t *testing.T) {

	var (
		endRunning = make(chan bool, 1)
		symbol     = "BSC"
	)

	scanner := testBlockScanner(symbol)

	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}
	scanner.SetBlockScanTargetFuncV2(testScanTargetFunc(symbol))
	scanner.SetRescanBlockHeight(6807927)
	scanner.Run()

	<-endRunning
}

func TestBlockScanner_ExtractTransactionAndReceiptData(t *testing.T) {

	var (
		symbol = "BSC"
		txid   = "0x0462f8db1152c6f1d184fb7947e021a10bb1265758ac24d59a7bcb237bee08e6"
	)

	scanner := testBlockScanner(symbol)

	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}

	result, contractResult, err := scanner.ExtractTransactionAndReceiptData(txid, testScanTargetFunc(symbol))
	if err != nil {
		t.Errorf("ExtractTransactionData unexpected error %v", err)
		return
	}

	for sourceKey, keyData := range result {
		log.Notice("account:", sourceKey)
		for _, data := range keyData {

			for i, input := range data.TxInputs {
				log.Std.Notice("data.TxInputs[%d]: %+v", i, input)
			}

			for i, output := range data.TxOutputs {
				log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
			}

			log.Std.Notice("data.Transaction: %+v", data.Transaction)
		}
	}

	for sourceKey, keyData := range contractResult {
		log.Notice("sourceKey:", sourceKey)
		log.Std.Notice("data.ContractTransaction: %+v", keyData)

		for i, event := range keyData.Events {
			log.Std.Notice("data.Contract[%d]: %+v", i, event.Contract)
			log.Std.Notice("data.Events[%d]: %+v", i, event)
		}
	}
}

func testScanTargetFunc(symbol string) openwallet.BlockScanTargetFuncV2 {
	var (
		addrs  = make(map[string]openwallet.ScanTargetResult)
		contracts   = make(map[string]openwallet.ScanTargetResult)
	)

	//添加监听的合约地址
	contract := &openwallet.SmartContract{
		Symbol:   symbol,
		Address:  "0x8076c74c5e3f5852037f31ff0093eeb8c8add8d3",
		Decimals: 9,
	}
	contract.ContractID = openwallet.GenContractID(contract.Symbol, contract.Address)
	//contract.SetABI(quorum.ERC20_ABI_JSON)
	contracts[contract.Address] = openwallet.ScanTargetResult{SourceKey: contract.ContractID, Exist: true, TargetInfo: contract}

	//添加监听的外部地址
	addrs["0x69aeba03354894e9136f60df0c4d2b3001ce1e3c"] = openwallet.ScanTargetResult{SourceKey: "sender", Exist: true}
	addrs["0x7321b4b889f2fc7c1735be59a684ef376d959290"] = openwallet.ScanTargetResult{SourceKey: "receiver", Exist: true}

	scanTargetFunc := func(target openwallet.ScanTargetParam) openwallet.ScanTargetResult {
		if target.ScanTargetType == openwallet.ScanTargetTypeContractAddress {
			return contracts[target.ScanTarget]
		} else if target.ScanTargetType == openwallet.ScanTargetTypeAccountAddress {
			return addrs[target.ScanTarget]
		}
		return openwallet.ScanTargetResult{SourceKey: "", Exist: false, TargetInfo: nil,}
	}

	return scanTargetFunc
}

func testBlockScanner(symbol string) openwallet.BlockScanner {
	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return nil
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()
	if scanner.SupportBlockchainDAI() {
		dbFilePath := filepath.Join("data", "db")
		dbFileName := "blockchain.db"
		file.MkdirAll(dbFilePath)
		dai, err := openwallet.NewBlockchainLocal(filepath.Join(dbFilePath, dbFileName), false)
		if err != nil {
			log.Error("NewBlockchainLocal err: %v", err)
			return nil
		}

		scanner.SetBlockchainDAI(dai)
	}
	sub := subscriberSingle{}
	scanner.AddObserver(&sub)

	return scanner
}
