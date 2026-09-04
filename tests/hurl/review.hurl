# Review E2E (GraphQL /query) — create category+merchant+product → review lifecycle.
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Review\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id}}}","variables":{"input":{"name":"Hurl Review Category {{uuid}}","description":"Category for review E2E","slug_category":"hurl-review-category-{{uuid}}","image_category":null}}}
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
    "query": "mutation { createMerchant(input: {user_id: {{userId}}, name: \"Hurl Review Merchant {{uuid}}\", description: \"Merchant for review E2E\", address: \"Jl. Ulasan No. 1\", contact_email: \"review.{{uuid}}@example.com\", contact_phone: \"081200000005\", status: \"active\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Hurl Review Product {{uuid}}","description":"Product for review E2E","price":75000,"countInStock":20,"brand":"HurlBrand","weight":300,"rating":5,"slugProduct":"hurl-review-product-{{uuid}}","imageProduct":null}}}
map: {"0":["variables.input.imageProduct"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createProduct.status" == "success"
[Captures]
productId: jsonpath "$.data.createProduct.data.id"

# ── Create Review ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createReview(input: {user_id: {{userId}}, product_id: {{productId}}, name: \"Hurl Reviewer\", comment: \"Produk bagus sekali!\", rating: 5}) { status message data { id product_id rating } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createReview.status" == "success"
[Captures]
reviewId: jsonpath "$.data.createReview.data.id"

# ── Find All Reviews ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllReviews(input: {page: 1, page_size: 10}) { status message data { id product_id rating } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllReviews.status" == "success"

# ── Find Reviews by Product ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findReviewsByProduct(input: {product_id: {{productId}}, page: 1, page_size: 10}) { status message data { id product_id rating } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findReviewsByProduct.status" == "success"

# ── Find Reviews by Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findReviewsByMerchant(input: {merchant_id: {{merchantId}}, page: 1, page_size: 10}) { status message data { id rating } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findReviewsByMerchant.status" == "success"

# ── Find Active Reviews ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveReviews(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveReviews.status" == "success"

# ── Update Review ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateReview(input: {review_id: {{reviewId}}, name: \"Hurl Reviewer\", comment: \"Ulasan diperbarui\", rating: 4}) { status message data { id rating } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateReview.status" == "success"

# ── Trash Review ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReview(input: {id: {{reviewId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReview.status" == "success"

# ── Find Trashed Reviews ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedReviews(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedReviews.status" == "success"

# ── Restore Review ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreReview(input: {id: {{reviewId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreReview.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReview(input: {id: {{reviewId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReview.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllReviews { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllReviews.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReview(input: {id: {{reviewId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReview.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllReviewsPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllReviewsPermanent.status" == "success"

# ── Review kedua — coverage deleteReviewPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createReview(input: {user_id: {{userId}}, product_id: {{productId}}, name: \"Hurl Reviewer\", comment: \"Review untuk permanent\", rating: 3}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createReview.status" == "success"
[Captures]
permReviewId: jsonpath "$.data.createReview.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReview(input: {id: {{permReviewId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReview.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteReviewPermanent(input: {id: {{permReviewId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteReviewPermanent.status" == "success"
