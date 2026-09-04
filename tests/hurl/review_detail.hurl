# Review Detail E2E (GraphQL /query) — create review → detail lifecycle.
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"ReviewDetail\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Dependencies: category + merchant + product + review ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id}}}","variables":{"input":{"name":"Hurl RD Category {{uuid}}","description":"Category for review detail E2E","slug_category":"hurl-rd-category-{{uuid}}","image_category":null}}}
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
    "query": "mutation { createMerchant(input: {user_id: {{userId}}, name: \"Hurl RD Merchant {{uuid}}\", description: \"Merchant for review detail E2E\", address: \"Jl. Review Detail No. 1\", contact_email: \"rd.{{uuid}}@example.com\", contact_phone: \"081200000007\", status: \"active\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateProductInput!){createProduct(input: $input){status message data{id}}}","variables":{"input":{"merchantId":{{merchantId}},"categoryId":{{categoryId}},"name":"Hurl RD Product {{uuid}}","description":"Product for review detail E2E","price":65000,"countInStock":15,"brand":"HurlBrand","weight":250,"rating":5,"slugProduct":"hurl-rd-product-{{uuid}}","imageProduct":null}}}
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
    "query": "mutation { createReview(input: {user_id: {{userId}}, product_id: {{productId}}, name: \"Hurl Reviewer\", comment: \"Review dengan detail\", rating: 5}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createReview.status" == "success"
[Captures]
reviewId: jsonpath "$.data.createReview.data.id"

# ── Create Review Detail ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createReviewDetail(input: {reviewId: {{reviewId}}, type: \"photo\", url: \"https://example.com/photo-{{uuid}}.jpg\", caption: \"Foto produk\"}) { status message data { id reviewId type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createReviewDetail.status" == "success"
[Captures]
reviewDetailId: jsonpath "$.data.createReviewDetail.data.id"

# ── Find All Review Details ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllReviewDetails(input: {page: 1, pageSize: 10}) { status message data { id reviewId type } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllReviewDetails.status" == "success"

# ── Find Review Detail by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findReviewDetailById(input: {id: {{reviewDetailId}}}) { status message data { id reviewId type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findReviewDetailById.status" == "success"

# ── Find Active Review Details ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveReviewDetails(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveReviewDetails.status" == "success"

# ── Update Review Detail ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateReviewDetail(input: {reviewDetailId: {{reviewDetailId}}, type: \"video\", url: \"https://example.com/video-{{uuid}}.mp4\", caption: \"Video produk\"}) { status message data { id type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateReviewDetail.status" == "success"

# ── Trash Review Detail ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReviewDetail(input: {id: {{reviewDetailId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReviewDetail.status" == "success"

# ── Find Trashed Review Details ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedReviewDetails(input: {page: 1, pageSize: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedReviewDetails.status" == "success"

# ── Restore Review Detail ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreReviewDetail(input: {id: {{reviewDetailId}}}) { status message data { id deletedAt } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreReviewDetail.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReviewDetail(input: {id: {{reviewDetailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReviewDetail.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllReviewDetails { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllReviewDetails.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReviewDetail(input: {id: {{reviewDetailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReviewDetail.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllReviewDetailsPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllReviewDetailsPermanent.status" == "success"

# ── Detail kedua — coverage deleteReviewDetailPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createReviewDetail(input: {reviewId: {{reviewId}}, type: \"photo\", url: \"https://example.com/perm-{{uuid}}.jpg\", caption: \"Detail permanent\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createReviewDetail.status" == "success"
[Captures]
permDetailId: jsonpath "$.data.createReviewDetail.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedReviewDetail(input: {id: {{permDetailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedReviewDetail.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteReviewDetailPermanent(input: {id: {{permDetailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteReviewDetailPermanent.status" == "success"
