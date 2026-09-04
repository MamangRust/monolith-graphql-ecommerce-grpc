# User E2E (GraphQL /query) — lifecycle + active/trashed + restore-all/permanent-all.
# Runner meng-inject testEmail/testPassword.

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"User\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id } } }"
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

# ── Find All Users ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findAllUsers(input: {page: 1, page_size: 10}) { status message data { id email } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllUsers.status" == "success"

# ── Create User ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createUser(input: {firstname: \"Alice\", lastname: \"Smith\", email: \"alice.{{testEmail}}\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id email } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createUser.status" == "success"
[Captures]
createdUserId: jsonpath "$.data.createUser.data.id"

# ── Find User by ID ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByIdUser(input: {id: {{createdUserId}}}) { status message data { id email } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByIdUser.status" == "success"

# ── Update User ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { updateUser(input: {id: {{createdUserId}}, firstname: \"Alicia\", lastname: \"Smith\", email: \"alice.{{testEmail}}\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id firstname } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.updateUser.status" == "success"

# ── Find Active Users ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByActiveUsers(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveUsers.status" == "success"

# ── Trash User ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{createdUserId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

# ── Find Trashed Users ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedUsers(input: {page: 1, page_size: 10}) { status message data { id } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedUsers.status" == "success"

# ── Restore User ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreUser(input: {id: {{createdUserId}}}) { status message data { id deleted_at } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreUser.status" == "success"

# ── Trash lagi → restoreAll → trash lagi → deleteAllPermanent ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{createdUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { restoreAllUser { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.restoreAllUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{createdUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteAllUserPermanent { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteAllUserPermanent.status" == "success"

# ── User kedua — coverage deleteUserPermanent (permanent/:id) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { createUser(input: {firstname: \"Bob\", lastname: \"Smith\", email: \"bob.{{testEmail}}\", password: \"password123\", confirm_password: \"password123\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.createUser.status" == "success"
[Captures]
permUserId: jsonpath "$.data.createUser.data.id"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { trashedUser(input: {id: {{permUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.trashedUser.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { deleteUserPermanent(input: {id: {{permUserId}}}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.deleteUserPermanent.status" == "success"
