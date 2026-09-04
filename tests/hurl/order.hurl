# Order E2E (GraphQL /query) — create category+merchant+product → order lifecycle + stats.
# Order create memerlukan objek shipping_address lengkap.
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Order\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Dependencies: category + merchant + product ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id}}}","variables":{"input":{"name":"Hurl Order Category {{uuid}}","description":"Category for order E2E","slug_category":"hurl-order-category-{{uuid}}","image_category":null}}}
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
    "query": "mutation { createMerchant(input: {user_id: {{userId}}, name: \"Hurl Order Merchant {{uuid}}\", description: \"Merchant for order E2E\", address: \"Jl. Pesanan No. 1\", contact_email: \"order.{{uuid}}@example.com\", contact_phone: \"081200000003\", status: \"active\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Hurl Order Product {{uuid}}","description":"Product for order E2E","price":25000,"countInStock":30,"brand":"HurlBrand","weight":100,"rating":5,"slugProduct":"hurl-order-product-{{uuid}}","imageProduct":null}}}
map: {"0":["variables.input.imageProduct"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# ── Create Order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createOrder(input: {merchant_id: {{merchantId}}, user_id: {{userId}}, total_price: 60000, items: [{product_id: {{productId}}, quantity: 2, price: 25000}], shipping: {alamat: \"Jl. Antah Berantah No. 9\", provinsi: \"DKI Jakarta\", kota: \"Jakarta Selatan\", courier: \"JNE\", shipping_method: \"REG\", shipping_cost: 10000, negara: \"Indonesia\"}}) { status message data { id merchant_id total_price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createOrder.status" == "success"
[Captures]
orderId: jsonpath "$.data.createOrder.data.id"

# ── Find All Orders ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllOrders(input: {page: 1, page_size: 10}) { status message data { id merchant_id total_price } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllOrders.status" == "success"

# ── Find Order by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findOrderById(input: {id: {{orderId}}}) { status message data { id total_price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findOrderById.status" == "success"

# ── Find Active Orders ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveOrders(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveOrders.status" == "success"

# ── Order items by order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findOrderItemsByOrder(input: {id: {{orderId}}}) { status message data { id order_id product_id quantity } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findOrderItemsByOrder.status" == "success"

# ── Stats: monthly/yearly revenue & total revenue (all variants) ──
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
    "query": "query { findMonthlyRevenueByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyRevenueByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyRevenueByMerchant.status" == "success"

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
    "query": "query { findMonthlyTotalRevenueByMerchant(input: {year: 2026, month: 4, merchant_id: {{merchantId}}}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalRevenueByMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalRevenueByMerchant(input: {year: 2026, merchant_id: {{merchantId}}}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalRevenueByMerchant.status" == "success"

# ── Update Order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateOrder(input: {order_id: {{orderId}}, user_id: {{userId}}, total_price: 70000, items: [{order_item_id: 0, product_id: {{productId}}, quantity: 2, price: 30000}], shipping: {shipping_id: 0, alamat: \"Jl. Baru No. 5\", provinsi: \"DKI Jakarta\", kota: \"Jakarta Pusat\", courier: \"JNE\", shipping_method: \"YES\", shipping_cost: 10000, negara: \"Indonesia\"}}) { status message data { id total_price } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateOrder.status" == "success"

# ── Trash Order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashOrder(input: {id: {{orderId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashOrder.status" == "success"

# ── Find Trashed Orders ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedOrders(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedOrders.status" == "success"

# ── Restore Order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreOrder(input: {id: {{orderId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreOrder.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashOrder(input: {id: {{orderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllOrders { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllOrders.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashOrder(input: {id: {{orderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllOrdersPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllOrdersPermanent.status" == "success"

# ── Order kedua — coverage deleteOrderPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createOrder(input: {merchant_id: {{merchantId}}, user_id: {{userId}}, total_price: 30000, items: [{product_id: {{productId}}, quantity: 1, price: 30000}], shipping: {alamat: \"Jl. Perm No. 1\", provinsi: \"DKI Jakarta\", kota: \"Jakarta Selatan\", courier: \"JNE\", shipping_method: \"REG\", shipping_cost: 5000, negara: \"Indonesia\"}}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createOrder.status" == "success"
[Captures]
permOrderId: jsonpath "$.data.createOrder.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashOrder(input: {id: {{permOrderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashOrder.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteOrderPermanent(input: {id: {{permOrderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteOrderPermanent.status" == "success"
