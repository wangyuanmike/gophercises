# Choose your own adventure

## Create structure Chapter and parse JSON
* Use this tool to generate a structure according to the json file automatically. https://mholt.github.io/json-to-go/
* Adjust structure name and field names
```
type Chapter struct {
    ...
}
```
* Create type story
```
type story map[string]Chapter
```
* Create function loadStory() to parse JSON into structure
```
func loadStory(io.Reader) story, error {
    - Use json.Decoder.Decode() instead of json.Unmarshal() to parse json, because json is stored in file 
}
```
## Create default template file
* Create a html file as the web page template and embed the fields of the story structure
## Create handler
* Create a global variable tpl for default template and initialize it in function init()
* Create type structure handler
```
type handler struct {
    t *template.Template
    s story     
}
```
* Create function NewHandler() as a handler factory. To make the documentation of NewHandler() more explicit, this function returns type http.Handler instead of handler.
```
func NewHandler(s story) http.Handler {
    return handler{tpl, s}  //tpl is the default template
}
```
* Implement method ServeHTTP() of interface http.Handler for struct handler
```
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    - Retrieve path from request URL (path indexes story chapter), esp. need to handle the root path, i.e. "/" -> "/intro"
    - Execute template with the path
}
```
## Create function main() to test run the first version
```
func main() {
    - Create two flags
    flag.Int("port", 8080, ...)
    flag.String("file", "gopher.json", ...)

    - Create json file reader f
    - Call loadStory(f) and load json into story isntance s 
    - Call NewHandler(s) to create a new handler instance h
    - Pass handler to http.ListenAndServe(addr, h)  //addr = :8080
}
``` 
## Create functional option type handlerOption 
```
type handlerOption func(h *handler)
```
* Functional option is a design pattern that could handle creating instance with opitons neatly. 
* Basically the factory function accepts [0, N] options, each of which is a function. It first create instance with default value(or we say default constructor), and then executes each functional option to manipulate corresponding property with custom value that is defined in the corresponding function. 
* Functional option is raised by [Robert Pike](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html) and [Dave Cheney](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).
## Create functional option for custom template
* Create the custom template file by coping the file from Jon's solution and adjust field names of struct Chapter that are embedded in the file
* Create the funtional option for custom template
```
func WithTemplate(h *handler) {
    h.t = parsed custom template
}
```
## Create functional option for custom path
* Add a new field pathFn to structure handler
```
type handler struct {
    t *template.Template
    s story
    pathFn func(*http.Request) string 
}
```
* Split the path retrieval logic from ServeHTTP() to a standalone function and adjust NewHandler() and ServeHTTP() accordingly
```
func defaultPathFn(r *http.Request) string {
    - Retrieve path from request URL (path indexes story chapter), esp. need to handle the root path, i.e. "/" -> "/intro"
}

func NewHandler(s story) http.Handler {
    return handler{tpl, s, defaultPathFn}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    - path := h.pathFn()
    - Execute template with path
}
```
* Create the functional option and corresponding function for custom path
```
func pathFn(r *http.Request) string {
    - Retrieve path from request URL and add "/story/" as the prefix of path
    - return path
}

func WithPathFn(h *handler) {
    h.pathFn = pathFn
}
```
## Add functional options to factory function NewHandler()
```
func NewHandler(s story, opts ...handlerOption) http.Handler {
    - create a handler instance with default value
    - range over all opts to apply custom value
}
```
## Create a Mux instance to route the request based on path
```
func main() {
    ...
    mux := http.NewServeMux()
    normalHandler := NewHandler(s)
    prettyHandler := NewHandler(s, WithTemplate(prettyTpl), WithPathFn(pathFn))
    mux.Handle("/", normalHandler)
    mux.Handle("/story", prettyHandler)

    log.Fatal(http.ListenAndServe(addr, mux))
}
```