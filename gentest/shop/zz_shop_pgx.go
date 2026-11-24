// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package shop

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
		"public.payment_method", "public._payment_method",
		"public.id_kind", "public._id_kind",
		"public.no_values", "public._no_values",
		"public.2fa_method", "public._2fa_method",
	})
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterTypes(types)
	return nil
}
