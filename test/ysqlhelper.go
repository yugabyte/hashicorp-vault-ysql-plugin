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
	//	return func() {}, "postgres://yugabyte:yugbayte@127.0.0.1:5433/yugabyte?sslmode=disable"
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
		Env:           []string{"YSQL_DB=lbcat", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=database"},
		Ports:         []string{"5433/tcp"},
		ContainerName: "yugabyte",
	})
	if err != nil {
		fmt.Println("Cound not start")
		t.Fatalf("Could not start docker Postgres: %s", err)
	}
	fmt.Println("Check the Docker")
	//time.Sleep(time.Second * 20)

	svc, err := runner.StartService(context.Background(), connectPostgres)
	if err != nil {
		t.Fatalf("Could not start docker Postgres: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}

func connectPostgres(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "secret"),
		Host:     fmt.Sprintf("%s:%d", host, port),
		Path:     "postgres",
		RawQuery: "sslmode=disable",
	}

	fmt.Println("Connection URL::", u.String())
	db, err := sql.Open("postgres", u.String())
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		fmt.Print("SQL.Open  success ")
	}

	///
	//	Trying a new SQL query::

	var dropStmt = `DROP TABLE IF EXISTS employee`
	if _, err := db.Exec(dropStmt); err != nil {
		fmt.Println("Unable to drop the table")
		fmt.Println(err)
	} else {
		fmt.Println("Drop the table")
	}

	var createStmt = `CREATE TABLE employee (id int PRIMARY KEY,
                                             name varchar,
                                             age int,
                                             language varchar)`
	if _, err := db.Exec(createStmt); err != nil {
		fmt.Println("Unable to create the table")
		fmt.Println(err)
	} else {
		fmt.Println("Created the table")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error here")
		fmt.Println(err)
		return nil, err
	}

	return docker.NewServiceURL(u), nil
}
