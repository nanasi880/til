package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/spanner"
)

var (
	projectID = flag.String("p", "", "GCP Project ID (required)")
	instance  = flag.String("i", "", "Instance Name (required)")
	database  = flag.String("d", "", "Database Name (required)")
)

func init() {
	flag.Parse()

	if *projectID == "" || *instance == "" || *database == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {

	ctx := context.Background()
	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", *projectID, *instance, *database)

	client, err := spanner.NewClient(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	txn1Open := make(chan struct{})
	defer close(txn1Open)

	txn2Done := make(chan struct{})

	// txn1 read
	txn1Done := make(chan struct{})
	go func() {
		defer close(txn1Done)

		var val string
		_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			txn.Query(ctx, spanner.NewStatement(`SELECT 1`)).Stop()

			txn1Open <- struct{}{}
			<-txn2Done
			return txn.Query(ctx, spanner.NewStatement(`SELECT Val FROM test WHERE ID = 1`)).Do(func(r *spanner.Row) error {
				return r.ColumnByName("Val", &val)
			})
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("got: ", val)
	}()

	// txn2 write
	go func() {
		defer close(txn2Done)

		<-txn1Open

		_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			return txn.BufferWrite([]*spanner.Mutation{
				spanner.Update("test", []string{
					"ID",
					"Val",
				}, []interface{}{
					1,
					"2",
				}),
			})
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-txn1Done
	<-txn2Done
}
