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
//

package ysql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/docker"

	_ "github.com/yugabyte/pgx/v4/stdlib"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {

	if version == "" {
		version = "latest"
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "yugabytedb/yugabyte",
		Cmd:           []string{"./bin/yugabyted", "start", "--daemon=false"},
		ImageTag:      version,
		Env:           []string{"YSQL_DB=testdb", "YSQL_PASSWORD=testsecret", "POSTGRES_DB=testdb", "POSTGRES_PASSWORD=testsecret"},
		Ports:         []string{"5433/tcp"},
		ContainerName: "yugabyte",
	})
	if err != nil {
		t.Fatalf("Could not start docker YugabyteDB: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectYugabyteDB)
	if err != nil {
		t.Fatalf("Could not start docker YugabyteDB: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}

func connectYugabyteDB(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("yugabyte", "testsecret"),
		Host:     fmt.Sprintf("%s:%d", host, port),
		Path:     "yugabyte",
		RawQuery: "sslmode=disable",
	}

	u_conn := u.String() + "&load_balance=true&yb_servers_refresh_interval=0"

	db, err := sql.Open("pgx", u_conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return docker.NewServiceURL(u), nil
}
