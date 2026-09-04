# Cart E2E (GraphQL /query) — create category+merchant+product → cart lifecycle.
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Cart\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id}}}","variables":{"input":{"name":"Hurl Cart Category {{uuid}}","description":"Category for cart E2E","slug_category":"hurl-cart-category-{{uuid}}","image_category":null}}}
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
    "query": "mutation { createMerchant(input: {user_id: {{userId}}, name: \"Hurl Cart Merchant {{uuid}}\", description: \"Merchant for cart E2E\", address: \"Jl. Keranjang No. 1\", contact_email: \"cart.{{uuid}}@example.com\", contact_phone: \"081200000002\", status: \"active\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Hurl Cart Product {{uuid}}","description":"Product for cart E2E","price":50000,"countInStock":50,"brand":"HurlBrand","weight":200,"rating":5,"slugProduct":"hurl-cart-product-{{uuid}}","imageProduct":null}}}
map: {"0":["variables.input.imageProduct"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# ── Add to cart ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCart(input: {user_id: {{userId}}, product_id: {{productId}}, quantity: 2}) { status message data { id product_id quantity } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCart.status" == "success"
[Captures]
cartItemId: jsonpath "$.data.createCart.data.id"

# ── Query cart ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllCarts(input: {user_id: {{userId}}, page: 1, page_size: 10}) { status message data { id product_id quantity } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllCarts.status" == "success"

# ── Delete item from cart (needs cart_id + user_id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteCart(input: {cart_id: {{cartItemId}}, user_id: {{userId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteCart.status" == "success"

# ── Delete all (needs cart_ids) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createCart(input: {user_id: {{userId}}, product_id: {{productId}}, quantity: 1}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createCart.status" == "success"
[Captures]
cartItemId2: jsonpath "$.data.createCart.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllCarts(input: {cart_ids: [{{cartItemId2}}]}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllCarts.status" == "success"
