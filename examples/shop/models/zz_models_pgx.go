// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// RegisterTypes registers every enum of the schema, and its array, with
// conn's type map, so pgx can encode and scan []OrderStatus. For
// pgxpool.Config.AfterConnect.
func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	types, err := conn.LoadTypes(ctx, []string{
		"public.order_status", "public._order_status",
	})
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterTypes(types)
	return nil
}
