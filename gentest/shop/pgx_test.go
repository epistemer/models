package shop_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/epistemer/models/gentest/shop"
)

func TestRegisterTypes(t *testing.T) {
	url := os.Getenv("MODELS_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("MODELS_TEST_DATABASE_URL isn't set")
	}
	ctx := t.Context()

	admin, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	dbname := fmt.Sprintf("models_test_%d", rand.Uint64())
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+dbname); err != nil {
		admin.Close(ctx)
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		if _, err := admin.Exec(ctx, "DROP DATABASE "+dbname); err != nil {
			t.Errorf("dropping %s: %v", dbname, err)
		}
		admin.Close(ctx)
	})

	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Database = dbname
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	ddl, err := os.ReadFile("../../testdata/shop/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, string(ddl)); err != nil {
		t.Fatal(err)
	}

	want := []shop.OrderStatus{shop.OrderStatusPaid, shop.OrderStatusShipped}
	var got []shop.OrderStatus
	if err := conn.QueryRow(ctx, "SELECT $1::order_status[]", want).Scan(&got); err == nil {
		t.Fatalf("encoding []OrderStatus without RegisterTypes unexpectedly succeeded: %v", got)
	}

	if err := shop.RegisterTypes(ctx, conn); err != nil {
		t.Fatal(err)
	}
	got = nil
	if err := conn.QueryRow(ctx, "SELECT $1::order_status[]", want).Scan(&got); err != nil {
		t.Fatalf("after RegisterTypes: %v", err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
