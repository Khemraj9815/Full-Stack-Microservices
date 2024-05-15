git clone <repository_url>

node server.js

## Code Explanation

This code sets up a simple HTTP server using Node.js to handle CRUD operations for a collection of blog posts stored in a JSON file (`data.json`). Let's break down how the code works and explain each endpoint:

1. **`loadBlogPosts` function**: This function asynchronously loads the blog posts from the `data.json` file into the `blogPosts` array when the server starts.

2. **`parseBody` middleware**: This middleware function is used to parse the request body as JSON. It listens for the `data` and `end` events on the request stream to accumulate the request body, then parses it into JSON and attaches it to `req.body`. It passes control to the next middleware or route handler using the `next` callback.

3. **Server creation and request handling**:
    - The server is created using `http.createServer`, and it listens for incoming HTTP requests.
    - The server handles different types of HTTP requests (`POST`, `PUT`, `GET`) to different endpoints.

4. **`POST /users/add` endpoint**:
    - This endpoint allows clients to add a new blog post. It expects JSON data containing the `title` and optional `content` of the new post in the request body.
    - When a request is received, the `parseBody` middleware parses the request body.
    - The new blog post is created with a unique ID (`Date.now()`), the provided title, and optional content. It's then added to the `blogPosts` array.
    - The updated `blogPosts` array is written back to the `data.json` file.
    - The server responds with a status code `201` (Created) and the JSON representation of the newly added post.

5. **`PUT /users/update/:id` endpoint**:
    - This endpoint allows clients to update an existing blog post specified by its ID.
    - It expects JSON data containing the updated properties of the blog post in the request body.
    - When a request is received, the `parseBody` middleware parses the request body.
    - The server finds the blog post with the specified ID and updates its properties with the provided data.
    - The updated `blogPosts` array is written back to the `data.json` file.
    - The server responds with a status code `200` (OK) and the JSON representation of the updated post.

6. **GET request handling**:
    - GET requests to any other endpoint (except `/users/add` and `/users/update/:id`) are currently not implemented and will result in an empty response.

7. **Error Handling**:
    - Errors during request handling, file reading, or file writing are caught and logged to the console with an appropriate status code (`500` for Internal Server Error).
    - Unsupported methods (`GET`, `DELETE`, etc.) or routes result in a `405` (Method Not Allowed) response.

## Endpoints

- `POST /users/add`: Add a new blog post. Expects JSON data with `title` and optional `content`.
- `PUT /users/update/:id`: Update an existing blog post specified by its ID. Expects JSON data with updated properties.
