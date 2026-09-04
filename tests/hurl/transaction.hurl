# Transaction E2E (GraphQL /query) — create category+merchant+product+order → transaction lifecycle + stats.
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Transaction\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
userId: jsonpath "$.data.registerUser.data.id"

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

# ── Dependencies: category + merchant + product + order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id}}}","variables":{"input":{"name":"Hurl Txn Category {{uuid}}","description":"Category for transaction E2E","slug_category":"hurl-txn-category-{{uuid}}","image_category":null}}}
map: {"0":["variables.input.image_category"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
categoryId: jsonpath "$.data.createCategory.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{userId}}, name: \"Hurl Txn Merchant {{uuid}}\", description: \"Merchant for transaction E2E\", address: \"Jl. Transaksi No. 1\", contact_email: \"txn.{{uuid}}@example.com\", contact_phone: \"081200000004\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Hurl Txn Product {{uuid}}","description":"Product for transaction E2E","price":50000,"countInStock":20,"brand":"HurlBrand","weight":200,"rating":5,"slugProduct":"hurl-txn-product-{{uuid}}","imageProduct":null}}}
map: {"0":["variables.input.imageProduct"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createOrder(input: {merchant_id: {{merchantId}}, user_id: {{userId}}, total_price: 50000, items: [{product_id: {{productId}}, quantity: 1, price: 50000}], shipping: {alamat: \"Jl. Antah Berantah No. 9\", provinsi: \"DKI Jakarta\", kota: \"Jakarta Selatan\", courier: \"JNE\", shipping_method: \"REG\", shipping_cost: 10000, negara: \"Indonesia\"}}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createOrder.status" == "success"
[Captures]
orderId: jsonpath "$.data.createOrder.data.id"

# ── Find All Transactions ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllTransaction(input: {page: 1, page_size: 10}) { status message data { id order_id amount } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllTransaction.status" == "success"

# ── Create Transaction ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createTransaction(input: {order_id: {{orderId}}, merchant_id: {{merchantId}}, payment_method: \"bank_transfer\", amount: 100000}) { status message data { id order_id amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createTransaction.status" == "success"
[Captures]
transactionId: jsonpath "$.data.createTransaction.data.id"

# ── Find Transaction by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTransactionById(input: {id: {{transactionId}}}) { status message data { id order_id amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTransactionById.status" == "success"

# ── Find Transactions by Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTransactionByMerchant(input: {merchant_id: {{merchantId}}, page: 1, page_size: 10}) { status message data { id merchant_id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTransactionByMerchant.status" == "success"

# ── Update Transaction ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateTransaction(input: {transaction_id: {{transactionId}}, order_id: {{orderId}}, merchant_id: {{merchantId}}, payment_method: \"e_wallet\", amount: 100000}) { status message data { id amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateTransaction.status" == "success"

# ── Stats: status & method success/failed (monthly & yearly) ──
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
    "query": "query { findMonthStatusSuccessByMerchant(input: {year: 2026, month: 4, merchant_id: {{merchantId}}}) { status message data { month total_success } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusSuccessByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusSuccessByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year total_success } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearStatusSuccessByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthStatusFailedByMerchant(input: {year: 2026, month: 4, merchant_id: {{merchantId}}}) { status message data { month total_failed } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthStatusFailedByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearStatusFailedByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year total_failed } } }"
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
    "query": "query { findMonthMethodByMerchantSuccess(input: {year: 2026, month: 4, merchant_id: {{merchantId}}}) { status message data { month payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodByMerchantSuccess.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodByMerchantSuccess(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year payment_method total_amount } } }"
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
    "query": "query { findMonthMethodByMerchantFailed(input: {year: 2026, month: 4, merchant_id: {{merchantId}}}) { status message data { month payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthMethodByMerchantFailed.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearMethodByMerchantFailed(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year payment_method total_amount } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearMethodByMerchantFailed.status" == "success"

# ── Find Active Transactions ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActive(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActive.status" == "success"

# ── Trash Transaction ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{transactionId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

# ── Find Trashed Transactions ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashed(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashed.status" == "success"

# ── Restore Transaction ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreTransaction(input: {id: {{transactionId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreTransaction.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{transactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllTransaction { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{transactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllTransactionPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllTransactionPermanent.status" == "success"

# ── Transaction kedua — coverage deleteTransactionPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createTransaction(input: {order_id: {{orderId}}, merchant_id: {{merchantId}}, payment_method: \"bank_transfer\", amount: 100000}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createTransaction.status" == "success"
[Captures]
permTransactionId: jsonpath "$.data.createTransaction.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedTransaction(input: {id: {{permTransactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedTransaction.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteTransactionPermanent(input: {id: {{permTransactionId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteTransactionPermanent.status" == "success"
