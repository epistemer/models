// Code generated from labels by github.com/epistemer/i18n. DO NOT EDIT.

package models

import "github.com/epistemer/i18n"

var (
	L10nOrderStatusPending     = i18n.Key("enum.order_status.pending", "Pending")
	L10nOrderStatusPaid        = i18n.Key("enum.order_status.paid", "Paid")
	L10nOrderStatusShipped     = i18n.Key("enum.order_status.shipped", "Shipped")
	L10nOrderStatusCancelled   = i18n.Key("enum.order_status.cancelled", "Cancelled")
	L10nCustomer               = i18n.Key("table.shop_customer", "Customer")
	L10nCustomer_Id            = i18n.Key("table.shop_customer.Customer_Id", "Id")
	L10nCustomer_Email         = i18n.Key("table.shop_customer.Customer_Email", "Email address")
	L10nCustomer_FullName      = i18n.Key("table.shop_customer.Customer_FullName", "Full Name")
	L10nOrder                  = i18n.Key("table.shop_order", "Order")
	L10nOrder_Id               = i18n.Key("table.shop_order.Order_Id", "Id")
	L10nOrder_CustomerId       = i18n.Key("table.shop_order.Order_CustomerId", "Customer Id")
	L10nOrder_Status           = i18n.Key("table.shop_order.Order_Status", "Status")
	L10nOrder_PlacedAt         = i18n.Key("table.shop_order.Order_PlacedAt", "Placed At")
	L10nOrderListView          = i18n.Key("table.shop_order_list_view", "Orders")
	L10nOrderListView_Id       = i18n.Key("table.shop_order_list_view.OrderListView_Id", "Id")
	L10nOrderListView_Status   = i18n.Key("table.shop_order_list_view.OrderListView_Status", "Status")
	L10nOrderListView_PlacedAt = i18n.Key("table.shop_order_list_view.OrderListView_PlacedAt", "Placed At")
	L10nOrderListView_Email    = i18n.Key("table.shop_order_list_view.OrderListView_Email", "Email")
)

// l10n holds every label that has a path.
var l10n = map[string]i18n.TKer{
	"enum.order_status.pending":            L10nOrderStatusPending,
	"enum.order_status.paid":               L10nOrderStatusPaid,
	"enum.order_status.shipped":            L10nOrderStatusShipped,
	"enum.order_status.cancelled":          L10nOrderStatusCancelled,
	"table.shop_customer":                  L10nCustomer,
	"table.shop_customer.id":               L10nCustomer_Id,
	"table.shop_customer.email":            L10nCustomer_Email,
	"table.shop_customer.full_name":        L10nCustomer_FullName,
	"table.shop_order":                     L10nOrder,
	"table.shop_order.id":                  L10nOrder_Id,
	"table.shop_order.customer_id":         L10nOrder_CustomerId,
	"table.shop_order.status":              L10nOrder_Status,
	"table.shop_order.placed_at":           L10nOrder_PlacedAt,
	"table.shop_order_list_view":           L10nOrderListView,
	"table.shop_order_list_view.id":        L10nOrderListView_Id,
	"table.shop_order_list_view.status":    L10nOrderListView_Status,
	"table.shop_order_list_view.placed_at": L10nOrderListView_PlacedAt,
	"table.shop_order_list_view.email":     L10nOrderListView_Email,
}

// L10n returns the label at path, or nil when there is none.
func L10n(path string) i18n.TKer { return l10n[path] }
