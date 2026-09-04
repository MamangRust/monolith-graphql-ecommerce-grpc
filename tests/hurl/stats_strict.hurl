# Strict statistics E2E (GraphQL /query).
# Menutup semua query stats: category pricing, order revenue, transaction status/method.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Stats\", lastname: \"Tester\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"{{testEmail}}\", password: \"{{testPassword}}\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"

# -----------------------------------------------------------------------------
# Category statistics
# -----------------------------------------------------------------------------
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPrices(input: {year: 2026, month: 4}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPrices.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPrices(input: {year: 2026}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPrices.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPricesByMerchant(input: {year: 2026, month: 4, merchant_id: 1}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPricesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPricesByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPricesByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPricesById(input: {year: 2026, month: 4, category_id: 1}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPricesById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPricesById(input: {year: 2026, category_id: 1}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPricesById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPrice(input: {year: 2026}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPrice.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPrice(input: {year: 2026}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPrice.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPriceByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPriceByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPriceByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPriceByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPriceById(input: {year: 2026, category_id: 1}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPriceById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPriceById(input: {year: 2026, category_id: 1}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPriceById.status" == "success"

# -----------------------------------------------------------------------------
# Order statistics
# -----------------------------------------------------------------------------
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalRevenue(input: {year: 2026, month: 4}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalRevenue(input: {year: 2026}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalRevenueByMerchant(input: {year: 2026, month: 4, merchant_id: 1}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalRevenueByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyRevenue(input: {year: 2026, month: 4}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyRevenue(input: {year: 2026}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyRevenue.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyRevenueByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyRevenueByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyRevenueByMerchant.status" == "success"

# -----------------------------------------------------------------------------
# Transaction statistics
# -----------------------------------------------------------------------------
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusSuccess(input: {year: 2026, month: 4}) { status message data { month total_success total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusSuccess(input: {year: 2026, month: 4}) { status message data { year total_success total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusFailed(input: {year: 2026, month: 4}) { status message data { month total_failed total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusFailed(input: {year: 2026, month: 4}) { status message data { year total_failed total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusSuccessByMerchant(input: {year: 2026, month: 4, merchant_id: 1}) { status message data { month total_success } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusSuccessByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusSuccessByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { year total_success } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusSuccessByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusFailedByMerchant(input: {year: 2026, month: 4, merchant_id: 1}) { status message data { month total_failed } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusFailedByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusFailedByMerchant(input: {year: 2026, merchant_id: 1}) { status message data { year total_failed } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusFailedByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodSuccess(input: {year: 2026, month: 4}) { status message data { month payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodSuccess(input: {year: 2026}) { status message data { year payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodByMerchantSuccess(input: {year: 2026, month: 4, merchant_id: 1}) { status message data { month payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodByMerchantSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodByMerchantSuccess(input: {year: 2026, merchant_id: 1}) { status message data { year payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodByMerchantSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodFailed(input: {year: 2026, month: 4}) { status message data { month payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodFailed(input: {year: 2026}) { status message data { year payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthMethodByMerchantFailed(input: {year: 2026, month: 4, merchant_id: 1}) { status message data { month payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodByMerchantFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodByMerchantFailed(input: {year: 2026, merchant_id: 1}) { status message data { year payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodByMerchantFailed.status" == "success"
