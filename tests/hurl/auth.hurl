# Auth E2E (GraphQL /query) — register → login → me → refresh → forgot/reset/verify.
# Runner meng-inject testEmail/testPassword; register auto-assign ROLE_ADMIN dan
# login tidak menuntut is_verified=true (FindByEmailAndVerify = FindByEmail).

# ── Register (public) ──
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { registerUser(input: {firstname: \"Hurl\", lastname: \"Auth\", email: \"{{testEmail}}\", password: \"{{testPassword}}\", confirm_password: \"{{testPassword}}\"}) { status message data { id email } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.registerUser.status" == "success"
[Captures]
userId: jsonpath "$.data.registerUser.data.id"

# ── Login (public) ──
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { loginUser(input: {email: \"{{testEmail}}\", password: \"{{testPassword}}\"}) { status message data { access_token refresh_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.loginUser.status" == "success"
[Captures]
accessToken: jsonpath "$.data.loginUser.data.access_token"
refreshToken: jsonpath "$.data.loginUser.data.refresh_token"

# ── Get Me (authorized) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "query { getMe(input: {access_token: \"{{accessToken}}\"}) { status message data { id email } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.getMe.status" == "success"

# ── Get Me tanpa token → 401 (auth middleware) ──
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "query { getMe(input: {access_token: \"invalid\"}) { status message } }"
}
HTTP 401

# ── Refresh Token (public) ──
POST {{baseUrl}}/query
Content-Type: application/json
{
    "query": "mutation { refreshToken(input: {refresh_token: \"{{refreshToken}}\"}) { status message data { access_token refresh_token } } }"
}
HTTP 200
[Asserts]
jsonpath "$.data.refreshToken.status" == "success"

# ── Forgot Password (email terdaftar; butuh auth middleware) ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { forgotPassword(input: {email: \"{{testEmail}}\"}) { status message } }"
}
HTTP 200

# ── Reset Password dengan token tidak valid → error ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { resetPassword(input: {reset_token: \"invalid-token\", password: \"newpassword123\", confirm_password: \"newpassword123\"}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.errors[0].message" exists

# ── Verify Code — kode tidak dikenal → error ──
POST {{baseUrl}}/query
Authorization: Bearer {{accessToken}}
Content-Type: application/json
{
    "query": "mutation { verifyCode(input: {code: \"invalid-code\"}) { status message } }"
}
HTTP 200
[Asserts]
jsonpath "$.errors[0].message" exists
