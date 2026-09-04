# Merchant Policies E2E (GraphQL /query) — create merchant → policy lifecycle.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"MerchantPolicies\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Hurl Policy Merchant {{uuid}}\", description: \"Merchant for policy E2E\", address: \"Jl. Kebijakan No. 1\", contact_email: \"policy.{{uuid}}@example.com\", contact_phone: \"081200000010\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Find All Merchant Policies ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllMerchantPolicies(input: {page: 1, page_size: 10}) { status message data { id merchant_id title } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllMerchantPolicies.status" == "success"

# ── Create Merchant Policy ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantPolicy(input: {merchant_id: {{merchantId}}, policy_type: \"return\", title: \"Kebijakan Retur {{uuid}}\", description: \"Pengembalian barang dalam 7 hari\"}) { status message data { id merchant_id title } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantPolicy.status" == "success"
[Captures]
policyId: jsonpath "$.data.createMerchantPolicy.data.id"

# ── Find Merchant Policy by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMerchantPolicyById(input: {id: {{policyId}}}) { status message data { id title } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMerchantPolicyById.status" == "success"

# ── Find Active Merchant Policies ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveMerchantPolicies(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveMerchantPolicies.status" == "success"

# ── Update Merchant Policy ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchantPolicy(input: {merchant_policy_id: {{policyId}}, policy_type: \"refund\", title: \"Kebijakan Refund {{uuid}}\", description: \"Pengembalian dana dalam 3 hari\"}) { status message data { id title } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchantPolicy.status" == "success"

# ── Trash Merchant Policy ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantPolicy(input: {id: {{policyId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantPolicy.status" == "success"

# ── Find Trashed Merchant Policies ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedMerchantPolicies(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedMerchantPolicies.status" == "success"

# ── Restore Merchant Policy ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreMerchantPolicy(input: {id: {{policyId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreMerchantPolicy.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantPolicy(input: {id: {{policyId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantPolicy.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllMerchantPolicies { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllMerchantPolicies.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantPolicy(input: {id: {{policyId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantPolicy.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllMerchantPoliciesPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllMerchantPoliciesPermanent.status" == "success"

# ── Policy kedua — coverage deleteMerchantPolicyPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantPolicy(input: {merchant_id: {{merchantId}}, policy_type: \"privacy\", title: \"Kebijakan Privasi {{uuid}}\", description: \"Kebijakan privasi data\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantPolicy.status" == "success"
[Captures]
permPolicyId: jsonpath "$.data.createMerchantPolicy.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantPolicy(input: {id: {{permPolicyId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantPolicy.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteMerchantPolicyPermanent(input: {id: {{permPolicyId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteMerchantPolicyPermanent.status" == "success"
