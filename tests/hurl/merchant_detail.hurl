# Merchant Detail E2E (GraphQL /query) — create merchant → detail lifecycle (upload cover & logo).
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"MerchantDetail\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Dependencies: merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Hurl Detail Merchant {{uuid}}\", description: \"Merchant for detail E2E\", address: \"Jl. Detail No. 1\", contact_email: \"detail.{{uuid}}@example.com\", contact_phone: \"081200000012\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Find All Merchant Details ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllMerchantDetails(input: {page: 1, page_size: 10}) { status message data { id merchant_id display_name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllMerchantDetails.status" == "success"

# ── Create Merchant Detail (multipart + cover & logo upload) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateMerchantDetailInput!){createMerchantDetail(input: $input){status message data{id display_name cover_image_url logo_url}}}","variables":{"input":{"merchant_id":{{merchantId}},"display_name":"Hurl Detail Store {{uuid}}","cover_image":null,"logo":null,"short_description":"Toko untuk detail E2E","website_url":"https://example.com/detail-{{uuid}}"}}}
map: {"0":["variables.input.cover_image"],"1":["variables.input.logo"]}
0: file,{{testImage}}; image/png
1: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantDetail.status" == "success"
[Captures]
detailId: jsonpath "$.data.createMerchantDetail.data.id"

# ── Find Merchant Detail by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMerchantDetailById(input: {id: {{detailId}}}) { status message data { id display_name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMerchantDetailById.status" == "success"

# ── Find Active Merchant Details ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveMerchantDetails(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveMerchantDetails.status" == "success"

# ── Update Merchant Detail (multipart + upload) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: UpdateMerchantDetailInput!){updateMerchantDetail(input: $input){status message data{id display_name}}}","variables":{"input":{"merchant_detail_id":{{detailId}},"display_name":"Hurl Detail Store Updated","cover_image":null,"logo":null,"short_description":"Toko detail diperbarui","website_url":"https://example.com/detail-upd-{{uuid}}"}}}
map: {"0":["variables.input.cover_image"],"1":["variables.input.logo"]}
0: file,{{testImage}}; image/png
1: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchantDetail.status" == "success"

# ── Trash Merchant Detail ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantDetail(input: {id: {{detailId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantDetail.status" == "success"

# ── Find Trashed Merchant Details ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedMerchantDetails(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedMerchantDetails.status" == "success"

# ── Restore Merchant Detail ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreMerchantDetail(input: {id: {{detailId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreMerchantDetail.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantDetail(input: {id: {{detailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantDetail.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllMerchantDetails { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllMerchantDetails.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantDetail(input: {id: {{detailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantDetail.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllMerchantDetailsPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllMerchantDetailsPermanent.status" == "success"

# ── Detail kedua — coverage deleteMerchantDetailPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateMerchantDetailInput!){createMerchantDetail(input: $input){status message data{id}}}","variables":{"input":{"merchant_id":{{merchantId}},"display_name":"Perm Detail Store {{uuid}}","cover_image":null,"logo":null,"short_description":"Detail untuk permanent","website_url":"https://example.com/perm-{{uuid}}"}}}
map: {"0":["variables.input.cover_image"],"1":["variables.input.logo"]}
0: file,{{testImage}}; image/png
1: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantDetail.status" == "success"
[Captures]
permDetailId: jsonpath "$.data.createMerchantDetail.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantDetail(input: {id: {{permDetailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantDetail.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteMerchantDetailPermanent(input: {id: {{permDetailId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteMerchantDetailPermanent.status" == "success"
