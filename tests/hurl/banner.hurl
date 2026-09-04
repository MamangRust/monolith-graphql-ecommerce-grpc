# Banner E2E (GraphQL /query) — lifecycle + active/trashed + restore-all/permanent-all.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Banner\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Find All Banners ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllBanners(input: {page: 1, page_size: 10}) { status message data { banner_id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllBanners.status" == "success"

# ── Create Banner ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createBanner(input: {name: \"Hurl Banner {{uuid}}\", start_date: \"2026-01-01\", end_date: \"2026-12-31\", start_time: \"00:00:00\", end_time: \"23:59:59\", is_active: true}) { status message data { banner_id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createBanner.status" == "success"
[Captures]
bannerId: jsonpath "$.data.createBanner.data.banner_id"

# ── Find Banner by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findBannerById(input: {id: {{bannerId}}}) { status message data { banner_id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findBannerById.status" == "success"

# ── Find Active Banners ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveBanners(input: {page: 1, page_size: 10}) { status message data { banner_id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveBanners.status" == "success"

# ── Update Banner ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateBanner(input: {banner_id: {{bannerId}}, name: \"Hurl Banner Updated\", start_date: \"2026-02-01\", end_date: \"2026-11-30\", start_time: \"08:00:00\", end_time: \"22:00:00\", is_active: true}) { status message data { banner_id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateBanner.status" == "success"

# ── Trash Banner ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashBanner(input: {id: {{bannerId}}}) { status message data { banner_id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashBanner.status" == "success"

# ── Find Trashed Banners ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedBanners(input: {page: 1, page_size: 10}) { status message data { banner_id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedBanners.status" == "success"

# ── Restore Banner ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreBanner(input: {id: {{bannerId}}}) { status message data { banner_id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreBanner.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashBanner(input: {id: {{bannerId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashBanner.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllBanners { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllBanners.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashBanner(input: {id: {{bannerId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashBanner.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllBannersPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllBannersPermanent.status" == "success"

# ── Banner kedua — coverage deleteBannerPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createBanner(input: {name: \"Perm Banner {{uuid}}\", start_date: \"2026-03-01\", end_date: \"2026-10-31\", start_time: \"00:00:00\", end_time: \"23:59:59\", is_active: false}) { status message data { banner_id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createBanner.status" == "success"
[Captures]
permBannerId: jsonpath "$.data.createBanner.data.banner_id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashBanner(input: {id: {{permBannerId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashBanner.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteBannerPermanent(input: {id: {{permBannerId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteBannerPermanent.status" == "success"
