# Order Item E2E (GraphQL /query) — query only (order item dibuat via createOrder).
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"OrderItem\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id}}}","variables":{"input":{"name":"Hurl OI Category {{uuid}}","description":"Category for order item E2E","slug_category":"hurl-oi-category-{{uuid}}","image_category":null}}}
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
    "query": "mutation { createMerchant(input: {user_id: {{userId}}, name: \"Hurl OI Merchant {{uuid}}\", description: \"Merchant for order item E2E\", address: \"Jl. Order Item No. 1\", contact_email: \"oi.{{uuid}}@example.com\", contact_phone: \"081200000006\", status: \"active\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Hurl OI Product {{uuid}}","description":"Product for order item E2E","price":40000,"countInStock":25,"brand":"HurlBrand","weight":150,"rating":5,"slugProduct":"hurl-oi-product-{{uuid}}","imageProduct":null}}}
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
    "query": "mutation { createOrder(input: {merchant_id: {{merchantId}}, user_id: {{userId}}, total_price: 40000, items: [{product_id: {{productId}}, quantity: 1, price: 40000}], shipping: {alamat: \"Jl. Antah Berantah No. 9\", provinsi: \"DKI Jakarta\", kota: \"Jakarta Selatan\", courier: \"JNE\", shipping_method: \"REG\", shipping_cost: 10000, negara: \"Indonesia\"}}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createOrder.status" == "success"
[Captures]
orderId: jsonpath "$.data.createOrder.data.id"

# ── Find All Order Items ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllOrderItems(input: {page: 1, page_size: 10}) { status message data { id order_id product_id quantity } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllOrderItems.status" == "success"

# ── Find Order Items by Order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findOrderItemsByOrder(input: {id: {{orderId}}}) { status message data { id order_id product_id quantity } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findOrderItemsByOrder.status" == "success"

# ── Find Active Order Items ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveOrderItems(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveOrderItems.status" == "success"

# ── Find Trashed Order Items ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedOrderItems(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedOrderItems.status" == "success"
