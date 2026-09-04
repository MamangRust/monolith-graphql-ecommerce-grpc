# Merchant Business E2E (GraphQL /query) — create merchant → business info lifecycle.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"MerchantBusiness\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Hurl Business Merchant {{uuid}}\", description: \"Merchant for business E2E\", address: \"Jl. Bisnis No. 1\", contact_email: \"biz.{{uuid}}@example.com\", contact_phone: \"081200000009\", status: \"active\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchant.status" == "success"
[Captures]
merchantId: jsonpath "$.data.createMerchant.data.id"

# ── Find All Merchant Businesses ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllMerchantBusinesses(input: {page: 1, page_size: 10}) { status message data { id merchant_id business_type } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllMerchantBusinesses.status" == "success"

# ── Create Merchant Business ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantBusiness(input: {merchant_id: {{merchantId}}, business_type: \"retail\", tax_id: \"TAX-{{uuid}}\", established_year: 2020, number_of_employees: 10, website_url: \"https://example.com/{{uuid}}\"}) { status message data { id merchant_id business_type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantBusiness.status" == "success"
[Captures]
businessId: jsonpath "$.data.createMerchantBusiness.data.id"

# ── Find Merchant Business by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMerchantBusinessById(input: {id: {{businessId}}}) { status message data { id business_type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMerchantBusinessById.status" == "success"

# ── Find Active Merchant Businesses ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveMerchantBusinesses(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveMerchantBusinesses.status" == "success"

# ── Update Merchant Business ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchantBusiness(input: {merchant_business_info_id: {{businessId}}, business_type: \"wholesale\", tax_id: \"TAX-UPD-{{uuid}}\", established_year: 2019, number_of_employees: 25, website_url: \"https://example.com/upd-{{uuid}}\"}) { status message data { id business_type } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchantBusiness.status" == "success"

# ── Trash Merchant Business ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantBusiness(input: {id: {{businessId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantBusiness.status" == "success"

# ── Find Trashed Merchant Businesses ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedMerchantBusinesses(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedMerchantBusinesses.status" == "success"

# ── Restore Merchant Business ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreMerchantBusiness(input: {id: {{businessId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreMerchantBusiness.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantBusiness(input: {id: {{businessId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantBusiness.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllMerchantBusinesses { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllMerchantBusinesses.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantBusiness(input: {id: {{businessId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantBusiness.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllMerchantBusinessesPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllMerchantBusinessesPermanent.status" == "success"

# ── Business kedua — coverage deleteMerchantBusinessPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantBusiness(input: {merchant_id: {{merchantId}}, business_type: \"retail\", tax_id: \"TAX-PERM-{{uuid}}\", established_year: 2021, number_of_employees: 5, website_url: \"https://example.com/perm-{{uuid}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantBusiness.status" == "success"
[Captures]
permBizId: jsonpath "$.data.createMerchantBusiness.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashMerchantBusiness(input: {id: {{permBizId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashMerchantBusiness.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteMerchantBusinessPermanent(input: {id: {{permBizId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteMerchantBusinessPermanent.status" == "success"
