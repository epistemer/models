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
}

// L10n returns the label at path, or nil when there is none.
func L10n(path string) i18n.TKer { return l10n[path] }
