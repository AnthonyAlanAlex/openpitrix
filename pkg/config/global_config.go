// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/koding/multiconfig"

	"openpitrix.io/openpitrix/pkg/logger"
)

type Config struct {
	Log   LogConfig
	Mysql MysqlConfig
}

type LogConfig struct {
	Level string `default:"info"` // debug, info, warn, error, fatal
}
type MysqlConfig struct {
	Host     string `default:"openpitrix-db"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"password"`
	Database string `default:"openpitrix"`
}

func (m *MysqlConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

func PrintUsage() {
	fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, "\nSupported environment variables:\n")
	e := newLoader("openpitrix")
	e.PrintEnvs(new(Config))
	fmt.Println("")
}

func GetFlagSet() *flag.FlagSet {
	flag.CommandLine.Usage = PrintUsage
	return flag.CommandLine
}

func ParseFlag() {
	GetFlagSet().Parse(os.Args[1:])
}

// TODO: load config from etcd
func LoadConf() *Config {
	ParseFlag()

	config := new(Config)
	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader("openpitrix"))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		logger.Panicf("Failed to load config: %+v", err)
		panic(err)
	}
	logger.SetLevelByString(config.Log.Level)
	logger.Debugf("LoadConf: %+v", config)
	return config
}
