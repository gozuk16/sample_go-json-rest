package main

import (
    "github.com/gozuk16/go-json-rest/rest"
    "log"
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
    "net/http"
    "time"
)

type Health struct {
	Name string
	Version string
	StartTime string
	Now string
}

const datetimeFormat string = "2006/01/02 15:04:05.000"

var version string
var startTime string

func main() {
	startTime = time.Now().Format(datetimeFormat)

	api := rest.NewApi()
	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/hello", hello),
		rest.Get("/health", getHealth),
		rest.Get("/status", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(statusMw.GetStatus())
		}),
		rest.Get("/dir", getDirs),
		rest.Get("/redirect", redirect),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

func hello(w rest.ResponseWriter, r *rest.Request) {
        w.WriteJson(map[string]string{"Body": "Hello World!"})
}

func getHealth(w rest.ResponseWriter, r *rest.Request) {
	now := time.Now().Format(datetimeFormat)
	h := Health{Name: "sample_go-json-rest", Version: version, StartTime: startTime, Now: now}
        w.WriteJson(h)
}

func getDirs(w rest.ResponseWriter, r *rest.Request) {
	var result string
	out, err := exec.Command("ls", "-a").Output()
	if err != nil {
		result = ("Command Exec Error.")
	} else {
		result = fmt.Sprintf("%s", out)
	}
	w.WriteJson(map[string]string{"ls":result})
}

func redirect(w rest.ResponseWriter, r *rest.Request) {
	f, err := os.Open("example.json")
	if err != nil{
		fmt.Println("error")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)

        w.Write(b)
}
