package main

import (
	//"github.com/gozuk16/go-json-rest/rest"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gozuk16/mylib"
	"gopkg.in/tylerb/graceful.v1"
)

type Health struct {
	Name      string
	Version   string
	StartTime string
	Now       string
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
		rest.Get("/mem", getMem),
		rest.Get("/redirect", redirect),
		rest.Get("/stop", func(w rest.ResponseWriter, req *rest.Request) {
			for cpt := 1; cpt <= 3; cpt++ {

				time.Sleep(time.Duration(1) * time.Second)

				w.WriteJson(map[string]string{
					"Message": fmt.Sprintf("%d seconds", cpt),
				})
				w.(http.ResponseWriter).Write([]byte("\n"))

				// Flush the buffer to client
				w.(http.Flusher).Flush()
			}
			os.Exit(0)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	server := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    ":8010",
			Handler: api.MakeHandler(),
		},
	}
	log.Fatal(server.ListenAndServe())
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
	w.WriteJson(map[string]string{"ls": result})
}

func getMem(w rest.ResponseWriter, r *rest.Request) {
	//w.WriteJson(map[string]string{"mem": mylib.Mem()})
	w.(http.ResponseWriter).Write(mylib.Mem())
}

func redirect(w rest.ResponseWriter, r *rest.Request) {
	f, err := os.Open("example.json")
	if err != nil {
		fmt.Println("error")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)

	w.(http.ResponseWriter).Write(b)
	//w.Write(b)
}
