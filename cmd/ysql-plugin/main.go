// Copyright (c) YugaByteDB, Inc.
//
//Licensed to YugabyteDB, Inc. under one or more contributor license agreements.
//See the NOTICE file distributed with this work for additional information regarding
//copyright ownership.
//
//YugabyteDB licenses this file to you under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software distributed under the License
//is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
//or implied.  See the License for the specific language governing permissions and limitations
//under the License.

package main

import (
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	ysql "github.com/yugabyte/hashicorp-vault-ysql-plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	err := Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func Run() error {
	dbType, err := ysql.New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database))

	return nil
}
