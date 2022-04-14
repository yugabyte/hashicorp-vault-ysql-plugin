package ysql_helper

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {
	return prepareTestContainer(t, version, "secret", "database")
}

func PrepareTestContainerWithPassword(t *testing.T, version, password string) (func(), string) {
	return prepareTestContainer(t, version, password, "database")
}

func prepareTestContainer(t *testing.T, version, password, db string) (func(), string) {
	if os.Getenv("PG_URL") != "" {
		return func() {}, os.Getenv("PG_URL")
	}

	if version == "" {
		version = "latest"
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		Cmd:           []string{"./bin/yugabyted", "start", "--daemon=false"},
		ImageRepo:     "yugabytedb/yugabyte",
		ImageTag:      version,
		Env:           []string{"YSQL_DB=lbcat", "POSTGRES_PASSWORD=" + password, "POSTGRES_DB=" + db},
		Ports:         []string{"5433/tcp"},
		ContainerName: "yugabyte",
	})

	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	svc, err := runner.StartService(context.Background(), connectPostgres(password))
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}

func connectPostgres(password string) docker.ServiceAdapter {
	return func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword("postgres", password),
			Host:     fmt.Sprintf("%s:%d", host, port),
			Path:     "postgres",
			RawQuery: "sslmode=disable",
		}
		fmt.Println(host, port)
		db, err := sql.Open("postgres", u.String())
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
}
