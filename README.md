# GO Auth

This an example of API with authentication in Go. It gives some separation of concerns by organizing the layers into different packages, isolating the capabilities on their own files, and abstracting the dependencies using design patterns. All dependencies are interpreted in the `main.go` file. 

Libraries used:
* Gin for the API
* JWT for generating authorization token
* MongoDB Client for database connection.

The application is capable of creating, authorizing, and updating users. All the operations are logged in the database and can be also retrieved in the API.


You can run everything via docker compose.

## Prerequisites

- Docker installed on your machine
  

## Setting Up the Application

Firstly, you need to create a `.env` file in the root folder containing:

```
DB_URI=mongodb://appmongo1:27017
DB_REPLICA_SET=myReplicaSet
JWT_SECRET=<YOUR_TOKEN_SECRET_HERE>
```

You can change the details of `DB_URI` and `DB_REPLICA_SET` as needed.

Then, build the API and MongoDB cluster by running then via docker compose:

```bash
docker compose up -d
```

## Using the application

### Creating an user

```
POST http://localhost:8080/users/
{
    "email": "test@test.com",
    "password": "nicepassword",
    "confirmPassword": "nicepassword",
    "name": "wow nice name"
}
```

### Authenticating an user

```
POST http://localhost:8080/authenticate/
{
    "email": "test@test.com",
    "password": "nicepassword"
}
```

This should return the `token` and `id` that you can use to access the API.

### Authorizing an user

Add an HTTP `Authorization` header with the format `Bearer {token}`

### Fetching logged-in user info

```
GET http://localhost:8080/me
HEADER Autorization: Bearer {token}
```

Returns the details of the logged-in user.

### Fetching user logs

```
GET http://localhost:8080/users/{id}/logs
HEADER Authorization: Bearer {token}
```

Returns the logs for the specified user.
