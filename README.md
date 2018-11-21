
# Local start
load local DB
$ heroku local
or
$ cd cmd/website/
$ export DATABASE_URL=postgres://perelin@localhost/twissr_local?sslmode=disable
$ export TEMPLATE_FOLDER_PREFIX=../../
$ go build && ./website OR fresh
 
 # INstalle/Update depnedencies
$ govendor fetch PACKAGE

# Dev Log

https://stackoverflow.com/questions/21532113/golang-converting-string-to-int64
i, err := strconv.ParseInt(s, 10, 64)

https://medium.com/@motyar/golang-pretty-print-structure-8617379e29f4
fmt.Printf("%#v", p) //with name, value and type

"path": "github.com/ChimeraCoder/anaconda",
"revision": "afffdbac178c0e78746ea42d7c1af963cb7ef28e",