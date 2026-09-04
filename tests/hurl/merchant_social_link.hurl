# Merchant Social Link E2E (GraphQL /query) — create merchant+detail → social link.
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"MerchantSocial\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Dependencies: merchant + merchant detail (upload cover & logo) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchant(input: {user_id: {{ownerId}}, name: \"Hurl Social Merchant {{uuid}}\", description: \"Merchant for social link E2E\", address: \"Jl. Sosial No. 1\", contact_email: \"social.{{uuid}}@example.com\", contact_phone: \"081200000011\", status: \"active\"}) { status message data { id } } }"
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
operations: {"query":"mutation($input: CreateMerchantDetailInput!){createMerchantDetail(input: $input){status message data{id}}}","variables":{"input":{"merchant_id":{{merchantId}},"display_name":"Hurl Social Store","cover_image":null,"logo":null,"short_description":"Toko untuk social link E2E","website_url":"https://example.com/{{uuid}}"}}}
map: {"0":["variables.input.cover_image"],"1":["variables.input.logo"]}
0: file,{{testImage}}; image/png
1: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantDetail.status" == "success"
[Captures]
detailId: jsonpath "$.data.createMerchantDetail.data.id"

# ── Create Merchant Social Link ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createMerchantSocialLink(input: {merchant_detail_id: {{detailId}}, platform: \"instagram\", url: \"https://instagram.com/store-{{uuid}}\"}) { status message data { id platform url } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createMerchantSocialLink.status" == "success"
[Captures]
socialId: jsonpath "$.data.createMerchantSocialLink.data.id"

# ── Update Merchant Social Link ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateMerchantSocialLink(input: {id: {{socialId}}, merchant_detail_id: {{detailId}}, platform: \"tiktok\", url: \"https://tiktok.com/@store-{{uuid}}\"}) { status message data { id platform url } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateMerchantSocialLink.status" == "success"
