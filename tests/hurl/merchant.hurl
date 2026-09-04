# Merchant E2E (GraphQL /query) — lifecycle + active/trashed + restore-all/permanent-all.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Merchant\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Find All Merchants ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllMerchants(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllMerchants.status" == "success"

# ── Create Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Hurl Merchant {{uuid}}\", description: \"Merchant created by hurl\", address: \"Jl. Merchant No. 1\", contact_email: \"merchant.{{uuid}}@example.com\", contact_phone: \"081200000001\", status: \"active\"}) { status message data { id name status } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Find Merchant by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMerchantById(input: {id: {{merchantId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMerchantById.status" == "success"

# ── Find Active Merchants ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveMerchants(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveMerchants.status" == "success"

# ── Update Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchant(input: {merchant_id: {{merchantId}}, user_id: {{ownerId}}, name: \"Hurl Merchant Updated\", description: \"Merchant updated\", address: \"Jl. Merchant No. 2\", contact_email: \"merchant.upd.{{uuid}}@example.com\", contact_phone: \"081200000002\", status: \"active\"}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchant.status" == "success"

# ── Trash Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchant(input: {id: {{merchantId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchant.status" == "success"

# ── Find Trashed Merchants ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedMerchants(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedMerchants.status" == "success"

# ── Restore Merchant ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreMerchant(input: {id: {{merchantId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreMerchant.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchant(input: {id: {{merchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllMerchants { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllMerchants.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchant(input: {id: {{merchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllMerchantsPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllMerchantsPermanent.status" == "success"

# ── Merchant kedua — coverage deleteMerchantPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Perm Merchant {{uuid}}\", description: \"Merchant for permanent delete\", address: \"Jl. Merchant No. 3\", contact_email: \"perm.{{uuid}}@example.com\", contact_phone: \"081200000003\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
permMerchantId: jsonpath "$.data.createMerchant.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchant(input: {id: {{permMerchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchant.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteMerchantPermanent(input: {id: {{permMerchantId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteMerchantPermanent.status" == "success"
