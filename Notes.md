Go module and go.mod contain - Collection of Go packages managed as single unit , provides dependency management
 and versioning of Project
 go.mod contains -  module name/path Go version , required Deps

 main belong cmd/api ? - main package represents an executable application
 cmd/api - clearly identifies the application entry point, separates app startup from reusable business logic


os.GetEnv = return a empty string
os.LookupEnv = Return both the value and a Boolean indicating whether the value exists

%w formatting verb used with fmt.Errorf to wrap an error

`json:"status"` = struct tag encoding package used this tag tp determine the JSON field name.

http.ResponseWriter = is an interface provided by net/http. it Allows handlers to construct HTTP response
\*http.Request - is pointer to the incoming HTTP request (it have info like method, headers etc)
&http.Server{} creates a pointer to an http,server value
why request use a pointer = so the functions can work with the request without copying its entire struct

any? is a alias for interface{} , it means the parameter can hold a value of any Go type

(\_ = json.NewEncoder(w).Encode(data))
json.Encoder(w) create a JSON Encoder that directly write to HTTP response write
.encode(data) coverts the Go Value to json and write it (returns error - discards the error)

configure timeouts? without it, slow or unresponsive connections can consume server for too long


context.Context carries request-scoped info , cancellation signals and deadlines through function calls

regexp.MustCompile - compiles regex and panics if the expression is valid

in struct where the values are without a name it is called struct embedding

Higher order function - when a function return another function


A single pgx.Conn is not safe for unrestricted  use - A connection pool manages multiple reusable
and allows concurrent requests to acquire connections safely

============================DATABASE========================================
Why does Connect() accept a context
-if application startup is cancelled, operations using the context can stop instead of continuing unnecessarily.
