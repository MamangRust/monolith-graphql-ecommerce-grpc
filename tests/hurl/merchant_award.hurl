# Merchant Award E2E (GraphQL /query) — create merchant → award lifecycle.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"MerchantAward\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Hurl Award Merchant {{uuid}}\", description: \"Merchant for award E2E\", address: \"Jl. Penghargaan No. 1\", contact_email: \"award.{{uuid}}@example.com\", contact_phone: \"081200000008\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Find All Merchant Awards ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllMerchantAwards(input: {page: 1, page_size: 10}) { status message data { id merchant_id title } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllMerchantAwards.status" == "success"

# ── Create Merchant Award ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantAward(input: {merchant_id: {{merchantId}}, title: \"Best Store {{uuid}}\", description: \"Penghargaan toko terbaik\", issued_by: \"SanEdge Awards\", issue_date: \"2026-01-01\", expiry_date: \"2027-01-01\", certificate_url: \"https://example.com/cert-{{uuid}}.pdf\"}) { status message data { id merchant_id title } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantAward.status" == "success"
[Captures]
awardId: jsonpath "$.data.createMerchantAward.data.id"

# ── Find Merchant Award by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMerchantAwardById(input: {id: {{awardId}}}) { status message data { id title } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMerchantAwardById.status" == "success"

# ── Find Active Merchant Awards ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveMerchantAwards(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveMerchantAwards.status" == "success"

# ── Update Merchant Award ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchantAward(input: {merchant_certification_id: {{awardId}}, title: \"Best Store Updated {{uuid}}\", description: \"Penghargaan diperbarui\", issued_by: \"SanEdge Awards\", issue_date: \"2026-02-01\", expiry_date: \"2027-02-01\", certificate_url: \"https://example.com/cert-upd-{{uuid}}.pdf\"}) { status message data { id title } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchantAward.status" == "success"

# ── Trash Merchant Award ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantAward(input: {id: {{awardId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantAward.status" == "success"

# ── Find Trashed Merchant Awards ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedMerchantAwards(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedMerchantAwards.status" == "success"

# ── Restore Merchant Award ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreMerchantAward(input: {id: {{awardId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreMerchantAward.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantAward(input: {id: {{awardId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantAward.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllMerchantAwards { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllMerchantAwards.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantAward(input: {id: {{awardId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantAward.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllMerchantAwardsPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllMerchantAwardsPermanent.status" == "success"

# ── Award kedua — coverage deleteMerchantAwardPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantAward(input: {merchant_id: {{merchantId}}, title: \"Perm Award {{uuid}}\", description: \"Award untuk permanent\", issued_by: \"SanEdge Awards\", issue_date: \"2026-03-01\", expiry_date: \"2027-03-01\", certificate_url: \"https://example.com/perm-{{uuid}}.pdf\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantAward.status" == "success"
[Captures]
permAwardId: jsonpath "$.data.createMerchantAward.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantAward(input: {id: {{permAwardId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantAward.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteMerchantAwardPermanent(input: {id: {{permAwardId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteMerchantAwardPermanent.status" == "success"
