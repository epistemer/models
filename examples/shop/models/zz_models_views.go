// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"context"

	"github.com/epistemer/models/column"
	"github.com/epistemer/models/listing"
)

var _ column.Kind // the import stays when no view has a column of a kind

// Listings is the listing methods of Queries, one per view.
type Listings interface {
	ListShopOrderListView(ctx context.Context, src listing.Source) ([]ShopOrderListView, int64, error)
}

var _ Listings = (*Queries)(nil)

// ListShopOrderListView returns the rows of shop_order_list_view
// that src selects, and how many rows its search and filters match.
func (q *Queries) ListShopOrderListView(ctx context.Context, src listing.Source) ([]ShopOrderListView, int64, error) {
	lq := src.ListingQuery()
	rowsSQL, countSQL, args, err := listing.Build(lq, "shop_order_list_view", []string{
		"id",
		"status",
		"placed_at",
		"email",
	})
	if err != nil {
		return nil, 0, err
	}
	rows, err := q.db.Query(ctx, rowsSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ShopOrderListView
	for rows.Next() {
		var i ShopOrderListView
		if err := rows.Scan(
			&i.ID,
			&i.Status,
			&i.PlacedAt,
			&i.Email,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, i)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	// The rows tell the total when they end the listing.
	n := len(items)
	if (lq.Limit == 0 || n < lq.Limit) && (lq.Offset == 0 || n > 0) {
		return items, int64(lq.Offset + n), nil
	}
	var total int64
	if err := q.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if n > 0 && total < int64(lq.Offset+n) {
		total = int64(lq.Offset + n)
	}
	return items, total, nil
}

// ShopOrderListViewColumns are shop_order_list_view's columns, for a table package (A4).
var ShopOrderListViewColumns = struct {
	ID       column.Column[ShopOrderListView]
	Status   column.Column[ShopOrderListView]
	PlacedAt column.Column[ShopOrderListView]
	Email    column.Column[ShopOrderListView]
}{
	ID: column.Column[ShopOrderListView]{
		Name:  "id",
		Label: L10nOrderListView_Id,
		Kind:  column.Integer,
		Get:   column.Field(column.Integer, func(r ShopOrderListView) any { return r.ID }),
	},
	Status: column.Column[ShopOrderListView]{
		Name:    "status",
		Label:   L10nOrderListView_Status,
		Kind:    column.Enum,
		Options: AllOrderStatusOptions(),
		Get: func(r ShopOrderListView) column.Value {
			return column.Value{Enum: r.Status}
		},
	},
	PlacedAt: column.Column[ShopOrderListView]{
		Name:  "placed_at",
		Label: L10nOrderListView_PlacedAt,
		Kind:  column.Time,
		Get:   column.Field(column.Time, func(r ShopOrderListView) any { return r.PlacedAt }),
	},
	Email: column.Column[ShopOrderListView]{
		Name:  "email",
		Label: L10nOrderListView_Email,
		Kind:  column.Text,
		Get:   column.Field(column.Text, func(r ShopOrderListView) any { return r.Email }),
	},
}
