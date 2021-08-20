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
	"errors"
	"path/filepath"

	"github.com/astaxie/beego/config"
	"github.com/shopspring/decimal"
)

const (
	maxAddresNum = 10000
)

//LoadConfig 读取配置
func (wm *WalletManager) LoadConfig() error {

	var (
		c   config.Configer
		err error
	)

	//读取配置
	absFile := filepath.Join(wm.Config.configFilePath, wm.Config.configFileName)
	c, err = config.NewConfig("ini", absFile)
	if err != nil {
		return errors.New("Config is not setup. Please run 'wmd Config -s <symbol>' ")
	}

	wm.Config.ServerAPI = c.String("serverAPI")
	wm.Config.Threshold, _ = decimal.NewFromString(c.String("threshold"))
	wm.Config.SumAddress = c.String("sumAddress")
	wm.Config.RPCUser = c.String("rpcUser")
	wm.Config.RPCPassword = c.String("rpcPassword")
	wm.Config.NodeInstallPath = c.String("nodeInstallPath")
	wm.Config.FeeLimit, _ = c.Int64("feeLimit")
	wm.Config.FeeMini, _ = c.Int64("feeMini")
	wm.Config.IsTestNet, _ = c.Bool("isTestNet")
	if wm.Config.IsTestNet {
		wm.Config.WalletDataPath = c.String("testNetDataPath")
	} else {
		wm.Config.WalletDataPath = c.String("mainNetDataPath")
	}
	wm.Config.IgnoreDustTRX, _ = decimal.NewFromString(c.String("ignoreDustTRX"))
	wm.WalletClient = NewClient(wm.Config.ServerAPI, "", false)

	return nil
}

//LoadAssetsConfig 加载外部配置
func (wm *WalletManager) LoadAssetsConfig(c config.Configer) error {

	//读取配置
	//absFile := filepath.Join(wm.Config.configFilePath, wm.Config.configFileName)
	//c, err := config.NewConfig("ini", absFile)
	//if err != nil {
	//	return errors.New("Config is not setup. Please run 'wmd Config -s <symbol>' ")
	//}
	wm.Config.ServerAPI = c.String("serverAPI")
	//wm.Config.Threshold, _ = decimal.NewFromString(c.String("threshold"))
	//wm.Config.SumAddress = c.String("sumAddress")
	//wm.Config.RPCUser = c.String("rpcUser")
	//wm.Config.RPCPassword = c.String("rpcPassword")
	//wm.Config.NodeInstallPath = c.String("nodeInstallPath")
	wm.Config.IsTestNet, _ = c.Bool("isTestNet")
	wm.Config.FeeLimit, _ = c.Int64("feeLimit")
	wm.Config.FeeMini, _ = c.Int64("feeMini")
	//if wm.Config.IsTestNet {
	//	wm.Config.WalletDataPath = c.String("testNetDataPath")
	//} else {
	//	wm.Config.WalletDataPath = c.String("mainNetDataPath")
	//}

	wm.WalletClient = NewClient(wm.Config.ServerAPI, "", false)
	wm.Config.DataDir = c.String("dataDir")
	wm.Config.IgnoreDustTRX, _ = decimal.NewFromString(c.String("ignoreDustTRX"))

	//数据文件夹
	wm.Config.makeDataDir()
	return nil
}
