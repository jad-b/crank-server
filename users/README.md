users
=====

## Auth Workflow
1. User POSTs username, password Basic-Auth header to `authenticate` endpoint
1. Server validates by password; creates token entry and return token in
   Authorization header
1. User include `Authorization: token token=<Token>,id=<UserID>` in API requests
1. Server verifies token & attaches User ID
    1. User exists (Authentication)
    1. Token is valid for user (Authorization)
1. If valid, request proceeds. Else, request is rejected with 401 Unauthorized.

