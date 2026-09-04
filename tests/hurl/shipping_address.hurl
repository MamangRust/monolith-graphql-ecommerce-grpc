# Shipping Address E2E (GraphQL /query).
# Tidak ada create/update route untuk shipping address (dibuat implisit saat
# order checkout), jadi suite ini memakai data seeder + trash/restore lifecycle.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Shipping\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Find All Shipping Addresses ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllShipping(input: {page: 1, page_size: 10}) { status message data { id order_id alamat } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllShipping.status" == "success"
[Captures]
shippingId: jsonpath "$.data.findAllShipping.data[0].id"

# ── Find Shipping by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findShippingById(input: {id: {{shippingId}}}) { status message data { id alamat } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findShippingById.status" == "success"

# ── Find Shipping by Order ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findShippingByOrder(input: {id: {{shippingId}}}) { status message data { id order_id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findShippingByOrder.status" == "success"

# ── Find Active Shipping ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveShipping(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveShipping.status" == "success"

# ── Trash Shipping ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedShipping(input: {id: {{shippingId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedShipping.status" == "success"

# ── Find Trashed Shipping ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedShipping(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedShipping.status" == "success"

# ── Restore Shipping ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreShipping(input: {id: {{shippingId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreShipping.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedShipping(input: {id: {{shippingId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedShipping.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllShipping { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllShipping.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedShipping(input: {id: {{shippingId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedShipping.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllShippingPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllShippingPermanent.status" == "success"

# ── Shipping kedua — coverage deleteShippingPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllShipping(input: {page: 1, page_size: 10}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllShipping.status" == "success"
[Captures]
permShippingId: jsonpath "$.data.findAllShipping.data[0].id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedShipping(input: {id: {{permShippingId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedShipping.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteShippingPermanent(input: {id: {{permShippingId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteShippingPermanent.status" == "success"
