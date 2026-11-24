// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package shop

import (
	"context"

	"github.com/epistemer/models/column"
	"github.com/epistemer/models/listing"
)

var _ column.Kind // the import stays when no view has a column of a kind

// Listings is the listing methods of Queries, one per view.
type Listings interface {
	ListShopOrderListView(ctx context.Context, src listing.Source) ([]ShopOrderListView, int64, error)
	ListShopPageView(ctx context.Context, src listing.Source) ([]ShopPageView, int64, error)
	ListShopKindsView(ctx context.Context, src listing.Source) ([]ShopKindsView, int64, error)
}

var _ Listings = (*Queries)(nil)

// ListShopOrderListView returns the rows of shop_order_list_view
// that src selects, and how many rows its search and filters match.
func (q *Queries) ListShopOrderListView(ctx context.Context, src listing.Source) ([]ShopOrderListView, int64, error) {
	lq := src.ListingQuery()
	rowsSQL, countSQL, args, err := listing.Build(lq, "shop_order_list_view", []string{
		"id",
		"status",
		"payment",
		"email",
		"order",
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
			&i.Payment,
			&i.Email,
			&i.Order,
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

// ListShopPageView returns the rows of shop_page_view
// that src selects, and how many rows its search and filters match.
func (q *Queries) ListShopPageView(ctx context.Context, src listing.Source) ([]ShopPageView, int64, error) {
	lq := src.ListingQuery()
	rowsSQL, countSQL, args, err := listing.Build(lq, "shop_page_view", []string{
		"id",
		"path",
		"seen_at",
	})
	if err != nil {
		return nil, 0, err
	}
	rows, err := q.db.Query(ctx, rowsSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ShopPageView
	for rows.Next() {
		var i ShopPageView
		if err := rows.Scan(
			&i.ID,
			&i.Path,
			&i.SeenAt,
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

// ListShopKindsView returns the rows of shop_kinds_view
// that src selects, and how many rows its search and filters match.
func (q *Queries) ListShopKindsView(ctx context.Context, src listing.Source) ([]ShopKindsView, int64, error) {
	lq := src.ListingQuery()
	rowsSQL, countSQL, args, err := listing.Build(lq, "shop_kinds_view", []string{
		"id",
		"name",
		"note",
		"ref",
		"small",
		"medium",
		"ratio",
		"score",
		"price",
		"active",
		"day",
		"at",
		"seen",
		"span",
		"state",
		"tags",
		"meta",
	})
	if err != nil {
		return nil, 0, err
	}
	rows, err := q.db.Query(ctx, rowsSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ShopKindsView
	for rows.Next() {
		var i ShopKindsView
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Note,
			&i.Ref,
			&i.Small,
			&i.Medium,
			&i.Ratio,
			&i.Score,
			&i.Price,
			&i.Active,
			&i.Day,
			&i.At,
			&i.Seen,
			&i.Span,
			&i.State,
			&i.Tags,
			&i.Meta,
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
	ID      column.Column[ShopOrderListView]
	Status  column.Column[ShopOrderListView]
	Payment column.Column[ShopOrderListView]
	Email   column.Column[ShopOrderListView]
	Order   column.Column[ShopOrderListView]
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
	Payment: column.Column[ShopOrderListView]{
		Name:    "payment",
		Label:   L10nOrderListView_Payment,
		Kind:    column.Enum,
		Options: AllPaymentMethodOptions(),
		Get: func(r ShopOrderListView) column.Value {
			if !r.Payment.Valid {
				return column.Value{Null: true}
			}
			return column.Value{Enum: r.Payment.PaymentMethod}
		},
	},
	Email: column.Column[ShopOrderListView]{
		Name:  "email",
		Label: L10nOrderListView_Email,
		Kind:  column.Text,
		Get:   column.Field(column.Text, func(r ShopOrderListView) any { return r.Email }),
	},
	Order: column.Column[ShopOrderListView]{
		Name:  "order",
		Label: L10nOrderListView_Order,
		Kind:  column.Text,
		Get:   column.Field(column.Text, func(r ShopOrderListView) any { return r.Order }),
	},
}

// ShopPageViewColumns are shop_page_view's columns, for a table package (A4).
var ShopPageViewColumns = struct {
	ID     column.Column[ShopPageView]
	Path   column.Column[ShopPageView]
	SeenAt column.Column[ShopPageView]
}{
	ID: column.Column[ShopPageView]{
		Name:  "id",
		Label: L10nPageView_Id,
		Kind:  column.Integer,
		Get:   column.Field(column.Integer, func(r ShopPageView) any { return r.ID }),
	},
	Path: column.Column[ShopPageView]{
		Name:  "path",
		Label: L10nPageView_Path,
		Kind:  column.Text,
		Get:   column.Field(column.Text, func(r ShopPageView) any { return r.Path }),
	},
	SeenAt: column.Column[ShopPageView]{
		Name:  "seen_at",
		Label: L10nPageView_SeenAt,
		Kind:  column.Time,
		Get:   column.Field(column.Time, func(r ShopPageView) any { return r.SeenAt }),
	},
}

// ShopKindsViewColumns are shop_kinds_view's columns, for a table package (A4).
// Left out, of a type no kind takes: span (interval), tags (text[]), meta (jsonb).
var ShopKindsViewColumns = struct {
	ID     column.Column[ShopKindsView]
	Name   column.Column[ShopKindsView]
	Note   column.Column[ShopKindsView]
	Ref    column.Column[ShopKindsView]
	Small  column.Column[ShopKindsView]
	Medium column.Column[ShopKindsView]
	Ratio  column.Column[ShopKindsView]
	Score  column.Column[ShopKindsView]
	Price  column.Column[ShopKindsView]
	Active column.Column[ShopKindsView]
	Day    column.Column[ShopKindsView]
	At     column.Column[ShopKindsView]
	Seen   column.Column[ShopKindsView]
	State  column.Column[ShopKindsView]
}{
	ID: column.Column[ShopKindsView]{
		Name:  "id",
		Label: L10nKindsView_Id,
		Kind:  column.Integer,
		Get:   column.Field(column.Integer, func(r ShopKindsView) any { return r.ID }),
	},
	Name: column.Column[ShopKindsView]{
		Name:  "name",
		Label: L10nKindsView_Name,
		Kind:  column.Text,
		Get:   column.Field(column.Text, func(r ShopKindsView) any { return r.Name }),
	},
	Note: column.Column[ShopKindsView]{
		Name:  "note",
		Label: L10nKindsView_Note,
		Kind:  column.Text,
		Get:   column.Field(column.Text, func(r ShopKindsView) any { return r.Note }),
	},
	Ref: column.Column[ShopKindsView]{
		Name:  "ref",
		Label: L10nKindsView_Ref,
		Kind:  column.Text,
		Get:   column.Field(column.Text, func(r ShopKindsView) any { return r.Ref }),
	},
	Small: column.Column[ShopKindsView]{
		Name:  "small",
		Label: L10nKindsView_Small,
		Kind:  column.Integer,
		Get:   column.Field(column.Integer, func(r ShopKindsView) any { return r.Small }),
	},
	Medium: column.Column[ShopKindsView]{
		Name:  "medium",
		Label: L10nKindsView_Medium,
		Kind:  column.Integer,
		Get:   column.Field(column.Integer, func(r ShopKindsView) any { return r.Medium }),
	},
	Ratio: column.Column[ShopKindsView]{
		Name:  "ratio",
		Label: L10nKindsView_Ratio,
		Kind:  column.Number,
		Get:   column.Field(column.Number, func(r ShopKindsView) any { return r.Ratio }),
	},
	Score: column.Column[ShopKindsView]{
		Name:  "score",
		Label: L10nKindsView_Score,
		Kind:  column.Number,
		Get:   column.Field(column.Number, func(r ShopKindsView) any { return r.Score }),
	},
	Price: column.Column[ShopKindsView]{
		Name:  "price",
		Label: L10nKindsView_Price,
		Kind:  column.Number,
		Get:   column.Field(column.Number, func(r ShopKindsView) any { return r.Price }),
	},
	Active: column.Column[ShopKindsView]{
		Name:  "active",
		Label: L10nKindsView_Active,
		Kind:  column.Bool,
		Get:   column.Field(column.Bool, func(r ShopKindsView) any { return r.Active }),
	},
	Day: column.Column[ShopKindsView]{
		Name:  "day",
		Label: L10nKindsView_Day,
		Kind:  column.Date,
		Get:   column.Field(column.Date, func(r ShopKindsView) any { return r.Day }),
	},
	At: column.Column[ShopKindsView]{
		Name:  "at",
		Label: L10nKindsView_At,
		Kind:  column.Time,
		Get:   column.Field(column.Time, func(r ShopKindsView) any { return r.At }),
	},
	Seen: column.Column[ShopKindsView]{
		Name:  "seen",
		Label: L10nKindsView_Seen,
		Kind:  column.Time,
		Get:   column.Field(column.Time, func(r ShopKindsView) any { return r.Seen }),
	},
	State: column.Column[ShopKindsView]{
		Name:    "state",
		Label:   L10nKindsView_State,
		Kind:    column.Enum,
		Options: AllOrderStatusOptions(),
		Get: func(r ShopKindsView) column.Value {
			if !r.State.Valid {
				return column.Value{Null: true}
			}
			return column.Value{Enum: r.State.OrderStatus}
		},
	},
}
