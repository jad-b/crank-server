users
=====

## Auth Workflow
1. User POSTs username, password Basic-Auth header to `authenticate` endpoint
1. Server validates by password; creates token entry and return token in
   Authorization header
1. User include `Authorization: token token=<Token>` in API requests
1. Server verifies token
    1. Token exists for a user(Authentication)
1. If valid, request proceeds. Else, request is rejected with 401 Unauthorized.
1. Actual handler can now check if the authenticated user is allowed to perform
   the requested action (Authorization)
