# Product E2E (GraphQL /query) — create merchant+category → product lifecycle → filter.
# imageProduct bersifat opsional (Upload), dilewati agar request tetap JSON.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Product\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
ownerId: jsonpath "$.data.registerUser.data.id"

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

# ── Dependencies: merchant + category ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Hurl Product Merchant {{uuid}}\", description: \"Merchant for product\", address: \"Jl. Produk No. 3\", contact_email: \"produk.{{uuid}}@example.com\", contact_phone: \"081200000031\", status: \"active\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id name}}}","variables":{"input":{"name":"Kopi-{{uuid}}","description":"Kategori kopi","slug_category":"kopi-{{uuid}}","image_category":null}}}
map: {"0":["variables.input.image_category"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
categoryId: jsonpath "$.data.createCategory.data.id"
categoryName: jsonpath "$.data.createCategory.data.name"

# ── Create Product (tanpa upload — imageProduct opsional) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id name price}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Kopi Susu {{uuid}}","description":"Es kopi susu gula aren","price":25000,"countInStock":100,"brand":"Nusantara","weight":250,"rating":5,"slugProduct":"kopi-susu-{{uuid}}","imageProduct":null}}}
map: {"0":["variables.input.imageProduct"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# ── Find All Products ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllProducts(input: {page: 1, pageSize: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllProducts.status" == "success"

# ── Find Product by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findProductById(input: {id: {{productId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findProductById.status" == "success"

# ── Find Products by Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findProductsByMerchant(input: {merchantId: {{merchantId}}, page: 1, pageSize: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findProductsByMerchant.status" == "success"

# ── Find Products by Category ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findProductsByCategory(input: {categoryName: \"{{categoryName}}\", page: 1, pageSize: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findProductsByCategory.status" == "success"

# ── Find Active Products ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveProducts(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveProducts.status" == "success"

# ── Update Product ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: UpdateProductInput!){updateProduct(input: $input){status message data{id name price}}}","variables":{"input":{"productId":{{productId}},"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Kopi Susu Premium {{uuid}}","description":"Es kopi susu premium","price":30000,"countInStock":90,"brand":"Nusantara","weight":250,"rating":5,"slugProduct":"kopi-susu-premium-{{uuid}}","imageProduct":null}}}
map: {"0":["variables.input.imageProduct"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.updateProduct.status" == "success"

# ── Trash Product ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{productId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

# ── Find Trashed Products ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedProducts(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedProducts.status" == "success"

# ── Restore Product ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreProduct(input: {id: {{productId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreProduct.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{productId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllProducts { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllProducts.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{productId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllProductsPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllProductsPermanent.status" == "success"

# ── Product kedua — coverage deleteProductPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id name price}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Kopi Perm {{uuid}}","description":"Produk untuk test permanent","price":15000,"countInStock":20,"brand":"Nusantara","weight":100,"rating":4,"slugProduct":"kopi-perm-{{uuid}}","imageProduct":null}}}
map: {"0":["variables.input.imageProduct"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
permProductId: jsonpath "$.data.createProduct.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedProduct(input: {id: {{permProductId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedProduct.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteProductPermanent(input: {id: {{permProductId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteProductPermanent.status" == "success"
