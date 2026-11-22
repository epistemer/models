// Code generated from labels by github.com/epistemer/i18n. DO NOT EDIT.

package shop

import "github.com/epistemer/i18n"

var (
	L10nOrderStatusPending          = i18n.Key("enum.order_status.pending", "pending")
	L10nOrderStatusPaid             = i18n.Key("enum.order_status.paid", "paid")
	L10nOrderStatusShipped          = i18n.Key("enum.order_status.shipped", "shipped")
	L10nPaymentMethodCreditCard     = i18n.Key("enum.payment_method.credit-card", "credit-card")
	L10nPaymentMethodBankTransfer   = i18n.Key("enum.payment_method.bank/transfer", "bank/transfer")
	L10nPaymentMethodCashOnDelivery = i18n.Key("enum.payment_method.cash:on:delivery", "cash:on:delivery")
	L10nPaymentMethodPayPal         = i18n.Key("enum.payment_method.pay-pal", "pay-pal")
	L10nPaymentMethodValue4         = i18n.Key("enum.payment_method.pay_pal", "pay_pal")
	L10nPaymentMethodValue5         = i18n.Key("enum.payment_method.!!", "!!")
	L10nPaymentMethodValue6         = i18n.Key("enum.payment_method.value_5", "value_5")
	L10nPaymentMethodCaf            = i18n.Key("enum.payment_method.café", "café")
	L10nPaymentMethodValue8         = i18n.Key("enum.payment_method.", "")
	L10nIDKindID                    = i18n.Key("enum.id_kind.id", "id")
	L10nIDKindUserID                = i18n.Key("enum.id_kind.user_id", "user_id")
	L10nIDKindUuid                  = i18n.Key("enum.id_kind.uuid", "uuid")
	L10n_2faMethodTotp              = i18n.Key("enum.2fa_method.totp", "totp")
	L10n_2faMethodSms               = i18n.Key("enum.2fa_method.sms", "sms")
	L10n_2faMethodU2f               = i18n.Key("enum.2fa_method.u2f", "u2f")
	L10n_2faMethodLevel3            = i18n.Key("enum.2fa_method.level_3", "level_3")
	L10nCustomer                    = i18n.Key("table.shop_customer", "Customer")
	L10nCustomer_Id                 = i18n.Key("table.shop_customer.Customer_Id", "Id")
	L10nCustomer_Email              = i18n.Key("table.shop_customer.Customer_Email", "Where receipts go")
	L10nCustomer_IdKind             = i18n.Key("table.shop_customer.Customer_IdKind", "Id Kind")
	L10nCustomer_SecondFactor       = i18n.Key("table.shop_customer.Customer_SecondFactor", "Second Factor")
	L10nOrder                       = i18n.Key("table.shop_order", "Order")
	L10nOrder_Id                    = i18n.Key("table.shop_order.Order_Id", "Id")
	L10nOrder_CustomerId            = i18n.Key("table.shop_order.Order_CustomerId", "Customer Id")
	L10nOrder_Status                = i18n.Key("table.shop_order.Order_Status", "Status")
	L10nOrder_Payment               = i18n.Key("table.shop_order.Order_Payment", "Payment")
	L10nOrder_Ref                   = i18n.Key("table.shop_order.Order_Ref", "Ref")
	L10nOrder_AddressLine2          = i18n.Key("table.shop_order.Order_AddressLine2", "Address Line2")
	L10nOrderLine                   = i18n.Key("table.shop_order_line", "Order Line")
	L10nOrderLine_OrderId           = i18n.Key("table.shop_order_line.OrderLine_OrderId", "Order Id")
	L10nOrderLine_LineNo            = i18n.Key("table.shop_order_line.OrderLine_LineNo", "Line No")
	L10nOrderLine_Title             = i18n.Key("table.shop_order_line.OrderLine_Title", "Title")
	L10nAuditLog                    = i18n.Key("table.audit_log", "Audit Log")
	L10nAuditLog_Id                 = i18n.Key("table.audit_log.AuditLog_Id", "Id")
	L10nAuditLog_HappenedAt         = i18n.Key("table.audit_log.AuditLog_HappenedAt", "Happened At")
	L10nAuditLog_What               = i18n.Key("table.audit_log.AuditLog_What", "What")
	L10nOrderSummary                = i18n.Key("table.shop_order_summary", "Orders at a glance")
	L10nOrderSummary_Id             = i18n.Key("table.shop_order_summary.OrderSummary_Id", "Id")
	L10nOrderSummary_Status         = i18n.Key("table.shop_order_summary.OrderSummary_Status", "Status")
	L10nOrderSummary_Lines          = i18n.Key("table.shop_order_summary.OrderSummary_Lines", "Lines")
	L10nOrderListView               = i18n.Key("table.shop_order_list_view", "Order List View")
	L10nOrderListView_Id            = i18n.Key("table.shop_order_list_view.OrderListView_Id", "Id")
	L10nOrderListView_Status        = i18n.Key("table.shop_order_list_view.OrderListView_Status", "Status")
	L10nOrderListView_Payment       = i18n.Key("table.shop_order_list_view.OrderListView_Payment", "Payment")
	L10nOrderListView_Email         = i18n.Key("table.shop_order_list_view.OrderListView_Email", "Email")
	L10nOrderListView_Order         = i18n.Key("table.shop_order_list_view.OrderListView_Order", "Order")
	L10nPageView                    = i18n.Key("table.shop_page_view", "Page View")
	L10nPageView_Id                 = i18n.Key("table.shop_page_view.PageView_Id", "Id")
	L10nPageView_Path               = i18n.Key("table.shop_page_view.PageView_Path", "Path")
	L10nPageView_SeenAt             = i18n.Key("table.shop_page_view.PageView_SeenAt", "Seen At")
	L10nKindsView                   = i18n.Key("table.shop_kinds_view", "Kinds View")
	L10nKindsView_Id                = i18n.Key("table.shop_kinds_view.KindsView_Id", "Id")
	L10nKindsView_Name              = i18n.Key("table.shop_kinds_view.KindsView_Name", "Name")
	L10nKindsView_Note              = i18n.Key("table.shop_kinds_view.KindsView_Note", "Note")
	L10nKindsView_Ref               = i18n.Key("table.shop_kinds_view.KindsView_Ref", "Ref")
	L10nKindsView_Small             = i18n.Key("table.shop_kinds_view.KindsView_Small", "Small")
	L10nKindsView_Medium            = i18n.Key("table.shop_kinds_view.KindsView_Medium", "Medium")
	L10nKindsView_Ratio             = i18n.Key("table.shop_kinds_view.KindsView_Ratio", "Ratio")
	L10nKindsView_Score             = i18n.Key("table.shop_kinds_view.KindsView_Score", "Score")
	L10nKindsView_Price             = i18n.Key("table.shop_kinds_view.KindsView_Price", "Price")
	L10nKindsView_Active            = i18n.Key("table.shop_kinds_view.KindsView_Active", "Active")
	L10nKindsView_Day               = i18n.Key("table.shop_kinds_view.KindsView_Day", "Day")
	L10nKindsView_At                = i18n.Key("table.shop_kinds_view.KindsView_At", "At")
	L10nKindsView_State             = i18n.Key("table.shop_kinds_view.KindsView_State", "State")
	L10nKindsView_Tags              = i18n.Key("table.shop_kinds_view.KindsView_Tags", "Tags")
	L10nKindsView_Meta              = i18n.Key("table.shop_kinds_view.KindsView_Meta", "Meta")
)

// l10n holds every label that has a path.
var l10n = map[string]i18n.TKer{
	"enum.order_status.pending":            L10nOrderStatusPending,
	"enum.order_status.paid":               L10nOrderStatusPaid,
	"enum.order_status.shipped":            L10nOrderStatusShipped,
	"enum.payment_method.credit-card":      L10nPaymentMethodCreditCard,
	"enum.payment_method.bank/transfer":    L10nPaymentMethodBankTransfer,
	"enum.payment_method.cash:on:delivery": L10nPaymentMethodCashOnDelivery,
	"enum.payment_method.pay-pal":          L10nPaymentMethodPayPal,
	"enum.payment_method.pay_pal":          L10nPaymentMethodValue4,
	"enum.payment_method.!!":               L10nPaymentMethodValue5,
	"enum.payment_method.value_5":          L10nPaymentMethodValue6,
	"enum.payment_method.café":             L10nPaymentMethodCaf,
	"enum.payment_method.":                 L10nPaymentMethodValue8,
	"enum.id_kind.id":                      L10nIDKindID,
	"enum.id_kind.user_id":                 L10nIDKindUserID,
	"enum.id_kind.uuid":                    L10nIDKindUuid,
	"enum.2fa_method.totp":                 L10n_2faMethodTotp,
	"enum.2fa_method.sms":                  L10n_2faMethodSms,
	"enum.2fa_method.u2f":                  L10n_2faMethodU2f,
	"enum.2fa_method.level_3":              L10n_2faMethodLevel3,
	"table.shop_customer":                  L10nCustomer,
	"table.shop_customer.id":               L10nCustomer_Id,
	"table.shop_customer.email":            L10nCustomer_Email,
	"table.shop_customer.id_kind":          L10nCustomer_IdKind,
	"table.shop_customer.second_factor":    L10nCustomer_SecondFactor,
	"table.shop_order":                     L10nOrder,
	"table.shop_order.id":                  L10nOrder_Id,
	"table.shop_order.customer_id":         L10nOrder_CustomerId,
	"table.shop_order.status":              L10nOrder_Status,
	"table.shop_order.payment":             L10nOrder_Payment,
	"table.shop_order.shop_ref":            L10nOrder_Ref,
	"table.shop_order.address_line2":       L10nOrder_AddressLine2,
	"table.shop_order_line":                L10nOrderLine,
	"table.shop_order_line.order_id":       L10nOrderLine_OrderId,
	"table.shop_order_line.line_no":        L10nOrderLine_LineNo,
	"table.shop_order_line.title":          L10nOrderLine_Title,
	"table.audit_log":                      L10nAuditLog,
	"table.audit_log.id":                   L10nAuditLog_Id,
	"table.audit_log.happened_at":          L10nAuditLog_HappenedAt,
	"table.audit_log.what":                 L10nAuditLog_What,
	"table.shop_order_summary":             L10nOrderSummary,
	"table.shop_order_summary.id":          L10nOrderSummary_Id,
	"table.shop_order_summary.status":      L10nOrderSummary_Status,
	"table.shop_order_summary.lines":       L10nOrderSummary_Lines,
	"table.shop_order_list_view":           L10nOrderListView,
	"table.shop_order_list_view.id":        L10nOrderListView_Id,
	"table.shop_order_list_view.status":    L10nOrderListView_Status,
	"table.shop_order_list_view.payment":   L10nOrderListView_Payment,
	"table.shop_order_list_view.email":     L10nOrderListView_Email,
	"table.shop_order_list_view.order":     L10nOrderListView_Order,
	"table.shop_page_view":                 L10nPageView,
	"table.shop_page_view.id":              L10nPageView_Id,
	"table.shop_page_view.path":            L10nPageView_Path,
	"table.shop_page_view.seen_at":         L10nPageView_SeenAt,
	"table.shop_kinds_view":                L10nKindsView,
	"table.shop_kinds_view.id":             L10nKindsView_Id,
	"table.shop_kinds_view.name":           L10nKindsView_Name,
	"table.shop_kinds_view.note":           L10nKindsView_Note,
	"table.shop_kinds_view.ref":            L10nKindsView_Ref,
	"table.shop_kinds_view.small":          L10nKindsView_Small,
	"table.shop_kinds_view.medium":         L10nKindsView_Medium,
	"table.shop_kinds_view.ratio":          L10nKindsView_Ratio,
	"table.shop_kinds_view.score":          L10nKindsView_Score,
	"table.shop_kinds_view.price":          L10nKindsView_Price,
	"table.shop_kinds_view.active":         L10nKindsView_Active,
	"table.shop_kinds_view.day":            L10nKindsView_Day,
	"table.shop_kinds_view.at":             L10nKindsView_At,
	"table.shop_kinds_view.state":          L10nKindsView_State,
	"table.shop_kinds_view.tags":           L10nKindsView_Tags,
	"table.shop_kinds_view.meta":           L10nKindsView_Meta,
}

// L10n returns the label at path, or nil when there is none.
func L10n(path string) i18n.TKer { return l10n[path] }
