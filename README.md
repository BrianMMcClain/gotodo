## Setup Postgres Container

```
docker run -itd -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=gotodo -p 5432:5432 postgres
```

## Supported Methods

- `GET /` - Get all ToDo items
- `GET /:id` - Get a specific ToDo item
- `POST /` - Create a new ToDo item
- `POST /:id` - Update an existing ToDo item
- `DELETE /:id` - Delete an existing ToDo item