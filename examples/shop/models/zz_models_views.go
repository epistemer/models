// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"context"

	"github.com/epistemer/models/listing"
)

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
