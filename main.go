package main

import (
	"os"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

func basicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, pass, ok := r.BasicAuth()
        if !ok || !checkCredentials(user, pass) {
            w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    }
}

func checkCredentials(user, pass string) bool {
    expectedUser := os.Getenv("BASIC_AUTH_USERNAME")
    expectedPass := os.Getenv("BASIC_AUTH_PASSWORD")
    return user == expectedUser && pass == expectedPass
}

func insertEmailHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprint(w, "Method Not Allowed")
        return
    }

    var email bson.M
    if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, "Bad Request")
        return
    }

    _, err := collection.InsertOne(context.Background(), email)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error inserting document: %v", err)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, "Email inserted")
}

func GetAllMails(collection *mongo.Collection) (string, error) {
    var results []bson.M

    findOptions := options.Find().SetProjection(bson.D{{"_id", 0}})

    cur, err := collection.Find(context.Background(), bson.D{{}}, findOptions)
    if err != nil {
        return "", err
    }
    defer cur.Close(context.Background())

    for cur.Next(context.Background()) {
        var doc bson.M
        if err := cur.Decode(&doc); err != nil {
            return "", err
        }
        results = append(results, doc)
    }

    if err := cur.Err(); err != nil {
        return "", err
    }

    jsonData, err := json.Marshal(results)
    if err != nil {
        return "", err
    }

    return string(jsonData), nil
}

func getAllMailsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprint(w, "Method Not Allowed")
        return
    }

    jsonMails, err := GetAllMails(collection)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error fetching mails: %v", err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, jsonMails)
}

var collection *mongo.Collection

func initializeDB() (*mongo.Collection, error) {
    ctx := context.Background()

    username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
    password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
    database := os.Getenv("MONGO_INITDB_DATABASE")
    mongoURI := fmt.Sprintf("mongodb://%s:%s@mongo:27017/%s?authSource=admin", username, password, database)

	clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }

    collection = client.Database("gosubscriber").Collection("emails")
    return collection, nil
}

func return404(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "Not Found")
}

func main() {
    _, err := initializeDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    http.HandleFunc("/", return404)
    http.HandleFunc("/subscribe", insertEmailHandler)
    http.HandleFunc("/emails", basicAuthMiddleware(getAllMailsHandler))

    fmt.Println("Server is listening on http://0.0.0.0:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
