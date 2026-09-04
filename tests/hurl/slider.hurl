# Slider E2E (GraphQL /query) — lifecycle + active/trashed + upload image (multipart).
# Runner meng-inject testEmail/testPassword/testImage.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Slider\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Find All Sliders ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllSliders(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllSliders.status" == "success"

# ── Create Slider (multipart + image upload) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateSliderRequest!){createSlider(input: $input){status message data{id name image}}}","variables":{"input":{"name":"Hurl Slider {{testEmail}}","image":null}}}
map: {"0":["variables.input.image"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createSlider.status" == "success"
[Captures]
sliderId: jsonpath "$.data.createSlider.data.id"

# ── Find Active Sliders ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findActiveSliders(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findActiveSliders.status" == "success"

# ── Update Slider (multipart + image upload) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: UpdateSliderRequest!){updateSlider(input: $input){status message data{id name}}}","variables":{"input":{"id":{{sliderId}},"name":"Hurl Slider Updated","image":null}}}
map: {"0":["variables.input.image"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.updateSlider.status" == "success"

# ── Trash Slider ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedSlider(input: {id: {{sliderId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedSlider.status" == "success"

# ── Find Trashed Sliders ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findTrashedSliders(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findTrashedSliders.status" == "success"

# ── Restore Slider ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreSlider(input: {id: {{sliderId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreSlider.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedSlider(input: {id: {{sliderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedSlider.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllSliders { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllSliders.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedSlider(input: {id: {{sliderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedSlider.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllSlidersPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllSlidersPermanent.status" == "success"

# ── Slider kedua — coverage deleteSliderPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: multipart/form-data
[Multipart]
operations: {"query":"mutation($input: CreateSliderRequest!){createSlider(input: $input){status message data{id}}}","variables":{"input":{"name":"Perm Slider {{testEmail}}","image":null}}}
map: {"0":["variables.input.image"]}
0: file,{{testImage}}; image/png
HTTP 200
[Asserts]
jsonpath "$.data.createSlider.status" == "success"
[Captures]
permSliderId: jsonpath "$.data.createSlider.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedSlider(input: {id: {{permSliderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedSlider.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteSliderPermanent(input: {id: {{permSliderId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteSliderPermanent.status" == "success"
