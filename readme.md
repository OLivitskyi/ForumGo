
To run project:
```
make run
```
   or
```
go run cmd/forum/main.go

```

Project runs on:
```
http://localhost:8080/
```


Project Overview:
A web forum with the following features:

1. **Communication:**
   - Enable users to communicate through posts and comments.
   - Allow registered users to associate categories with their posts.
   - Posts and comments should be visible to all users, but only registered users can create them.

2. **SQLite Database:**
   - Utilize SQLite for data storage, managing users, posts, and comments.
   - Implement at least one SELECT, one CREATE, and one INSERT query.
   - Refer to the entity relationship diagram for structuring the database.

3. **Authentication:**
   - Enable user registration with email, username, and password.
   - Use cookies for login sessions with an expiration date.
   - Implement error handling for duplicate email and encrypted password storage (optional).

4. **Communication Features:**
   - Only registered users can create posts and comments.
   - Allow registered users to associate categories with their posts.
   - Posts and comments should be visible to all users, with non-registered users limited to viewing.

5. **Likes and Dislikes:**
   - Registered users can like or dislike posts and comments.
   - Display the number of likes and dislikes to all users.

6. **Filter Mechanism:**
   - Implement a filter for posts based on categories, created posts, and liked posts.
   - Consider categories as subforums, dedicated to specific topics.
   - Filtering options are available only for registered users, referring to the logged-in user.

Official project description:
	https://01.kood.tech/git/root/public/src/branch/master/subjects/forum

Project audit page:
	https://01.kood.tech/git/root/public/src/branch/master/subjects/forum/audit


Project tree:

```
Forum
├── Makefile
├── build
│   └── app
├── cmd
│   └── forum
│       └── main.go
├── configs
│   └── config.json
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── server
│   │       ├── config.go
│   │       ├── database.go
│   │       ├── handlers.go
│   │       └── server.go
│   ├── auth
│   │   └── users.go
│   ├── db
│   │   └── forum.db
│   └── model
│       ├── post.go
│       ├── reaction.go
│       └── user.go
├── readme.md
└── web
    ├── static
    │   ├── bootstrap.bundle.js
    │   ├── bootstrap.css
    │   ├── bootstrap.js
    │   └── main.css
    └── templates
        ├── createPost.html
        ├── home.html
        ├── login.html
        ├── main.html
        ├── notfound404.html
        ├── post.html
        └── register.html

```
