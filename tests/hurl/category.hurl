# Category E2E (GraphQL /query) — lifecycle + stats + upload image_category (multipart).
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Category\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Find All Categories ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllCategories(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllCategories.status" == "success"

# ── Create Category (multipart + image upload) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id name slug_category image_category}}}","variables":{"input":{"name":"Hurl Category {{testEmail}}","description":"Category created by hurl E2E test","slug_category":"hurl-category-{{testEmail}}","image_category":null}}}
map: {"0":["variables.input.image_category"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
categoryId: jsonpath "$.data.createCategory.data.id"
categoryName: jsonpath "$.data.createCategory.data.name"

# ── Find Category by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findCategoryById(input: {id: {{categoryId}}}) { status message data { id name } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findCategoryById.status" == "success"

# ── Find Active Categories ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveCategories(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveCategories.status" == "success"

# ── Stats: monthly/yearly pricing & total pricing (all variants) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPrice(input: {year: 2026}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPrice.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPrice(input: {year: 2026}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPrice.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthPriceById(input: {year: 2026, category_id: {{categoryId}}}) { status message data { month category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthPriceById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearPriceById(input: {year: 2026, category_id: {{categoryId}}}) { status message data { year category_id total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearPriceById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPrices(input: {year: 2026, month: 4}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPrices.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPrices(input: {year: 2026}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPrices.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findMonthlyTotalPricesById(input: {year: 2026, month: 4, category_id: {{categoryId}}}) { status message data { month total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findMonthlyTotalPricesById.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findYearlyTotalPricesById(input: {year: 2026, category_id: {{categoryId}}}) { status message data { year total_revenue } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findYearlyTotalPricesById.status" == "success"

# ── Update Category (multipart + image upload) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: UpdateCategoryInput!){updateCategory(input: $input){status message data{id name}}}","variables":{"input":{"category_id":{{categoryId}},"name":"Hurl Category Updated","description":"Category updated by hurl E2E test","slug_category":"hurl-category-updated-{{testEmail}}","image_category":null}}}
map: {"0":["variables.input.image_category"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.updateCategory.status" == "success"

# ── Trash Category ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashCategory(input: {id: {{categoryId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashCategory.status" == "success"

# ── Find Trashed Categories ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedCategories(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedCategories.status" == "success"

# ── Restore Category ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreCategory(input: {id: {{categoryId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreCategory.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashCategory(input: {id: {{categoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllCategories { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllCategories.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashCategory(input: {id: {{categoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllCategoriesPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllCategoriesPermanent.status" == "success"

# ── Category kedua — coverage deleteCategoryPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateCategoryInput!){createCategory(input: $input){status message data{id}}}","variables":{"input":{"name":"Perm Category {{testEmail}}","description":"Category for permanent delete","slug_category":"perm-category-{{testEmail}}","image_category":null}}}
map: {"0":["variables.input.image_category"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createCategory.status" == "success"
[Captures]
permCategoryId: jsonpath "$.data.createCategory.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashCategory(input: {id: {{permCategoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashCategory.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteCategoryPermanent(input: {id: {{permCategoryId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteCategoryPermanent.status" == "success"
