package main

import (
	"context"
	"fmt"
	"neonexcore/pkg/graphql"
)

// Example models
type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Post struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	AuthorID uint  `json:"author_id"`
}

// Mock data
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
}

var posts = []Post{
	{ID: 1, Title: "First Post", Content: "Hello World", AuthorID: 1},
	{ID: 2, Title: "Second Post", Content: "GraphQL is awesome", AuthorID: 1},
	{ID: 3, Title: "Third Post", Content: "NeonexCore rocks", AuthorID: 2},
}

func main() {
	// Build schema using fluent API
	builder := graphql.NewBuilder()

	// Define Query type
	builder.Query(
		graphql.F("user", graphql.TypeObject,
			func(ctx context.Context, parent interface{}, args map[string]interface{}) (interface{}, error) {
				id := uint(args["id"].(float64))
				for _, u := range users {
					if u.ID == id {
						return u, nil
					}
				}
				return nil, fmt.Errorf("user not found")
			},
			graphql.WithDescription("Get a user by ID"),
			graphql.WithElementType("User"),
			graphql.WithArgs(
				graphql.Arg("id", graphql.TypeID, graphql.ArgRequired(), graphql.ArgDescription("User ID")),
			),
		),
		graphql.F("users", graphql.TypeList,
			func(ctx context.Context, parent interface{}, args map[string]interface{}) (interface{}, error) {
				return users, nil
			},
			graphql.WithDescription("Get all users"),
			graphql.WithElementType("User"),
		),
		graphql.F("post", graphql.TypeObject,
			func(ctx context.Context, parent interface{}, args map[string]interface{}) (interface{}, error) {
				id := uint(args["id"].(float64))
				for _, p := range posts {
					if p.ID == id {
						return p, nil
					}
				}
				return nil, fmt.Errorf("post not found")
			},
			graphql.WithDescription("Get a post by ID"),
			graphql.WithElementType("Post"),
			graphql.WithArgs(
				graphql.Arg("id", graphql.TypeID, graphql.ArgRequired()),
			),
		),
		graphql.F("posts", graphql.TypeList,
			func(ctx context.Context, parent interface{}, args map[string]interface{}) (interface{}, error) {
				return posts, nil
			},
			graphql.WithDescription("Get all posts"),
			graphql.WithElementType("Post"),
		),
	)

	// Define Mutation type
	builder.Mutation(
		graphql.F("createUser", graphql.TypeObject,
			func(ctx context.Context, parent interface{}, args map[string]interface{}) (interface{}, error) {
				input := args["input"].(map[string]interface{})
				newUser := User{
					ID:    uint(len(users) + 1),
					Name:  input["name"].(string),
					Email: input["email"].(string),
					Age:   int(input["age"].(float64)),
				}
				users = append(users, newUser)
				return newUser, nil
			},
			graphql.WithDescription("Create a new user"),
			graphql.WithElementType("User"),
			graphql.WithArgs(
				graphql.Arg("input", graphql.TypeObject, graphql.ArgRequired(), graphql.ArgElementType("CreateUserInput")),
			),
		),
		graphql.F("createPost", graphql.TypeObject,
			func(ctx context.Context, parent interface{}, args map[string]interface{}) (interface{}, error) {
				input := args["input"].(map[string]interface{})
				newPost := Post{
					ID:       uint(len(posts) + 1),
					Title:    input["title"].(string),
					Content:  input["content"].(string),
					AuthorID: uint(input["authorId"].(float64)),
				}
				posts = append(posts, newPost)
				return newPost, nil
			},
			graphql.WithDescription("Create a new post"),
			graphql.WithElementType("Post"),
			graphql.WithArgs(
				graphql.Arg("input", graphql.TypeObject, graphql.ArgRequired(), graphql.ArgElementType("CreatePostInput")),
			),
		),
	)

	// Define User type from struct
	builder.TypeFromStruct("User", User{}, "A user in the system")

	// Define Post type from struct
	builder.TypeFromStruct("Post", Post{}, "A blog post")

	// Define input types
	builder.Input("CreateUserInput",
		graphql.IF("name", graphql.TypeString, graphql.IFRequired(), graphql.IFDescription("User's name")),
		graphql.IF("email", graphql.TypeString, graphql.IFRequired(), graphql.IFDescription("User's email")),
		graphql.IF("age", graphql.TypeInt, graphql.IFRequired(), graphql.IFDescription("User's age")),
	)

	builder.Input("CreatePostInput",
		graphql.IF("title", graphql.TypeString, graphql.IFRequired()),
		graphql.IF("content", graphql.TypeString, graphql.IFRequired()),
		graphql.IF("authorId", graphql.TypeID, graphql.IFRequired()),
	)

	// Define enum
	builder.Enum("PostStatus",
		graphql.EV("DRAFT", graphql.EVDescription("Post is in draft")),
		graphql.EV("PUBLISHED", graphql.EVDescription("Post is published")),
		graphql.EV("ARCHIVED", graphql.EVDescription("Post is archived")),
	)

	// Build schema
	schema := builder.Build()

	// Print schema SDL
	fmt.Println("=== GraphQL Schema ===")
	fmt.Println(schema.String())
	fmt.Println()

	// Create executor
	executor := graphql.NewExecutor(schema)

	// Example queries
	fmt.Println("=== Example Query: Get all users ===")
	queryAllUsers := &graphql.Query{
		Query: "query { users }",
	}
	response := executor.Execute(context.Background(), queryAllUsers)
	responseJSON, _ := response.ToJSON()
	fmt.Println(string(responseJSON))
	fmt.Println()

	fmt.Println("=== Example Query: Get user by ID ===")
	queryUserByID := &graphql.Query{
		Query: "query { user(id: 1) }",
	}
	response = executor.Execute(context.Background(), queryUserByID)
	responseJSON, _ = response.ToJSON()
	fmt.Println(string(responseJSON))
	fmt.Println()

	fmt.Println("=== Example Mutation: Create user ===")
	mutationCreateUser := &graphql.Query{
		Query: "mutation { createUser(input: {name: \"Bob\", email: \"bob@example.com\", age: 35}) }",
	}
	response = executor.Execute(context.Background(), mutationCreateUser)
	responseJSON, _ = response.ToJSON()
	fmt.Println(string(responseJSON))
	fmt.Println()
}
