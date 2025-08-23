# Chirpy
## About

**Chirpy** is a small server for hypothetical blog aggregator, written in *golang*, with use of goose and sqlc. Exists for educational purposes, such as creating an REST API server and documentation practice.

### Requirements

It is required to have [Go](https://go.dev/doc/install) and [PostgreSQL](https://www.postgresql.org/download/) to run the server.

### Database configuration

1. Make sure you're on version 15+ of Postgres:

```
psql --version
```
2. (Linux only) Update postgres password:

```
sudo passwd postgres
```
Enter a password, and be sure you won't forget it. You can just use something easy like `postgres`.

3. Start the Postgres server in the background

    * Mac: `brew services start postgresql@15`
    * Linux: `sudo service postgresql start`

4. Enter the psql shell:

    * Mac: `psql postgres`
    * Linux: `sudo -u postgres psql`

5. Create a new database (`chirpy` is quite good):

```
CREATE DATABASE chirpy;
```

6. Connect to db:

```
\c chirpy
```

7. Set the user password (Linux only). I used `postgres` by default:

```
ALTER USER postgres PASSWORD 'postgres';
```

8. You can type `exit` to leave the **psql** shell.

### Run up migrations

To set up database for further api requests it is necessary to migrate up sql files in *sql/scheme* directory using *goose*

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Run`goose -version` to make sure if it installed correctly.

For running migration it is necessary to `cd` in directory where the migration are located and run

```
goose postgres "postgres://username:password@localhost:5432/database" up

"db_url": "postgres://username:password@localhost:5432/database?sslmode=disable"
```

`username` - is username you set previously

`password` - is password you set previously

`@localhost` - by default it supposed that you're running it locally. You may not, if you want so

`5432` - is default `port` for SQL databases, and for **PostgreSQL** too

`/database` - is your name of database, set previously

`?sslmode=disable` - query for app, it doesn't supposed to try to use SSL locally

### Environment

**.env** file consist of:

`DB_URL="postgres://username:password@localhost:5432/database?sslmode=disable"`

Check upper chapter for fields explanation. Addititional field is `?sslmode=disable` - query for app, it doesn't supposed to try to use SSL locally

`PLATFORM="dev"`

*dev* is set for database reset option in development environment, for wide testing possibilities.

`SECRET="key"`

*key* is nice long string serving as a secret token for JWT verification. You may create your own using `openssl rand -base64 64`. Copy from terminal and pass it here.

`POLKA_KEY="key"`

Here the *key* is 32 symbol long string, saved as API key for so called imaginery side service *POLKA*, sending request for updating users subscription. Made up for webhooks implementation.

## Usage

It is required to build and run the server. You may [find API documentation here](https://github.com/ar3ty/chirpy/tree/main/docs/API.md)