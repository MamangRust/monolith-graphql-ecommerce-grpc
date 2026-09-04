# Strict authorization/rules E2E (GraphQL /query).
# - Akses tanpa token harus ditolak (401) oleh auth middleware.
# - User ter-register (auto ROLE_ADMIN) bisa mengakses data role.
# - User seeder tanpa role tetap bisa mengakses (middleware GraphQL hanya cek JWT).
# Runner meng-inject adminEmail/adminPassword/userEmail/userPassword.

# Unauthenticated access must be rejected.
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "query { findAllRole(input: {page: 1, page_size: 10}) { status message data { id name } } }"
}
HTTP 401

# Register + login the admin (register auto-assigns ROLE_ADMIN).
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Rules\", lastname: \"Admin\", email: \"{{adminEmail}}\", password: \"{{adminPassword}}\", confirm_password: \"{{adminPassword}}\"}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"

POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"{{adminEmail}}\", password: \"{{adminPassword}}\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
adminToken: jsonpath "$.data.loginUser.data.access_token"

# Admin can read protected role data.
POST {{baseUrl}}/query
Authorization: Bearer {{adminToken}}
Content-Type: application/json
{
    "query": "query { findAllRole(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllRole.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{adminToken}}
Content-Type: application/json
{
    "query": "query { findByActiveRole(input: {page: 1, page_size: 10}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByActiveRole.status" == "success"

POST {{baseUrl}}/query
Authorization: Bearer {{adminToken}}
Content-Type: application/json
{
    "query": "query { findByTrashedRole(input: {page: 1, page_size: 10}) { status message data { id } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findByTrashedRole.status" == "success"

# Seeded role-less user: login works (middleware hanya cek JWT valid).
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"{{userEmail}}\", password: \"{{userPassword}}\"}) { status message data { access_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
userToken: jsonpath "$.data.loginUser.data.access_token"

POST {{baseUrl}}/query
Authorization: Bearer {{userToken}}
Content-Type: application/json
{
    "query": "query { findAllRole(input: {page: 1, page_size: 10}) { status message data { id name } pagination { current_page } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.findAllRole.status" == "success"
