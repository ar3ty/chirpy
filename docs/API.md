## API
### Users group
#### Create user
```
POST http://localhost:8080/api/users
```
Accepts JSON in body:
```
{
    "password": "password",
    "email": "email"
}
```
Gets email and password, saves user into database, returns user info about created user in body with JSON:
```
{
    "id": "5a47789c-a617-444a-8a80-b50359247804"
    "created_at": "2021-07-01T00:00:00Z",
    "updated_at": "2021-07-01T00:00:00Z",
    "email": "mail@example.com",
    "is_chirpy_red": "false",
}
```
#### Update user
```
PUT http://localhost:8080/api/users
```
Accepts JWT in header "Authorization: Bearer KEY".
Accepts JSON in body:
```
{
    "password": "password",
    "email": "email"
}
```
Gets email and password, if authorized, saves updated data into database, returns user info about created user in body with JSON:
```
{
    "id": "5a47789c-a617-444a-8a80-b50359247804"
    "created_at": "2021-07-01T00:00:00Z",
    "updated_at": "2021-07-01T00:00:00Z",
    "email": "mail@example.com",
    "is_chirpy_red": "false",
}
```
### Chirps group
#### Create chirp
```
POST http://localhost:8080/api/chirps
```
Accepts JWT in header "Authorization: Bearer KEY".
Accepts JSON in body:
```
{
    "body": "text"
}
```
If authorized and chirp is valid, creates new chirp record and responds info about created unit in JSON:
```
{
  "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
  "created_at": "2021-01-01T00:00:00Z",
  "updated_at": "2021-01-01T00:00:00Z",
  "body": "Hello, world!",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```
#### Get chirps
```
GET http://localhost:8080/api/chirps
```
Optional query parameters:

`author_id` to return only the chirps for that author;

`sort` to sort the response by creation date in specified order. May be *asc* or *desc*, ascending or descending, respectively. Default blank value is descending.

Returns array of chirps in JSON:
```
[
  {
    "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
    "created_at": "2021-01-01T00:00:00Z",
    "updated_at": "2021-01-01T00:00:00Z",
    "body": "Yo fam this feast is lit ong",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  },
  {
    "id": "f0f87ec2-a8b5-48cc-b66a-a85ce7c7b862",
    "created_at": "2022-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "body": "What's good king?",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
]
```
#### Get chirp
```
GET http://localhost:8080/api/chirps/{chirpID}
```
Respond with requested chirp JSON:
```
{
  "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
  "created_at": "2021-01-01T00:00:00Z",
  "updated_at": "2021-01-01T00:00:00Z",
  "body": "Hello, world!",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```
#### Delete chirp
```
DELETE http://localhost:8080/api/chirps/{chirpID}
```
Accepts JWT in header "Authorization: Bearer KEY".

If user have sufficient authority rights, deletes specified chirp.

If successful, responds with code 204 and withot body.
### Auth group
#### Login
```
POST http://localhost:8080/api/login
```
Accepts JSON in body:
```
{
    "password": "password",
    "email": "email"
}
```
Gets user by email, makes JWT and refresh token, returns user info and tokens in body with JSON:
```
{
    "id": "5a47789c-a617-444a-8a80-b50359247804"
    "created_at": "2021-07-01T00:00:00Z",
    "updated_at": "2021-07-01T00:00:00Z",
    "email": "mail@example.com",
    "is_chirpy_red": "false",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "refresh_token": "56aa826d22baab4b5ec2cea41a59ecbba03e542aedbb31d9b80326ac8ffcfa2a"
}
```
#### Refresh token
```
POST http://localhost:8080/api/refresh
```
Accepts refresh token in header "Authorization: Bearer KEY".
Makes new JWT and returns it with JSON in body:
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```
#### Revoke refresh_token
```
POST http://localhost:8080/api/revoke
```
Accepts refresh token in header "Authorization: Bearer KEY".
Revokes the token in the database.
#### Upgrade subscription (webhook)
```
POST http://localhost:8080/api/polka/webhooks
```
Accepts Polka (hypothetic payment system) APIkey in header "Authorization: Polka KEY"
Accepts JSON in http-body:
```
{
    "event": "user.upgraded",
    "data": {
        "user_id: "3311741c-680c-4546-99f3-fc9efac2036c"
    }
}
```
Authorized and valid request is granted and upgrades the user subscription status in database.
### Readiness
```
GET http://localhost:8080/api/healthz
```
Return "OK" as a plain text. Simple check for server listening.
## Administration
### Reset 
```
POST http://localhost:8080/admin/reset
```
Resets database. Nullifies number of visits (hits).
For development usage.
Only allowed with 'dev' platform configuration
### Metrics
```
GET http://localhost:8080/admin/metrics
```
Returns information about number of visits of main page
## Resource
### Blank main page
```
GET http://localhost:8080/app/
```
Endpoint for serving a page from a filesystem. Responds as a fileserver for index.html in main directory
### Assets
```
GET http://localhost:8080/app/assets/logo.png
```
Only one asset is provided ^^