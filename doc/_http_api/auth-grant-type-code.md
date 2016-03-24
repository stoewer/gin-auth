---
title: "Authenticate: grant type code"
requests:
  - title: 1. Start authorization
    method: GET
    url: https://<host>/oauth/authorize
    header:
      Content-Type: application/json
      Accept: foo
    parameter:
      response_type: Must be set to 'code'
      client_id: The ID of a registered client
      redirect_uri: The URL to redirect to after authentication and approval
      scope: A comma separated list of scopes
      state: A random string to protect against CSRF (optional)
    response:
      - status_code: 302 (moved temporarily)
        header:
          Location: https://<host>/oauth/login_page
        parameter:
          request_id: A randomly created id associated with this grant request
        body: |
            ```json
            {
                "foo": "bar"
            }
            ```
    error:
      - status_code: 400 (bad request)
        conditions:
          - The client id is unknown
          - The redirect URL does not math exactly one registered URL for the client
          - The redirect URL does not use https
          - One of the given scopes is not registered
        body: An error page is shown

---

This authorization method should be used by server side web applications.