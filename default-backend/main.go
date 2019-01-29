package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// FormatHeader name of the header used to extract the format
	FormatHeader = "X-Format"

	// CodeHeader name of the header used as source of the HTTP statu code to return
	CodeHeader = "X-Code"

	// ContentType name of the header that defines the format of the reply
	ContentType = "Content-Type"

	// OriginalURI name of the header with the original URL from NGINX
	OriginalURI = "X-Original-URI"

	// Namespace name of the header that contains information about the Ingress namespace
	Namespace = "X-Namespace"

	// IngressName name of the header that contains the matched Ingress
	IngressName = "X-Ingress-Name"

	// ServiceName name of the header that contains the matched Service in the Ingress
	ServiceName = "X-Service-Name"

	// ServicePort name of the header that contains the matched Service port in the Ingress
	ServicePort = "X-Service-Port"

	// defaultHTMLTemplate
	defaultHTMLTemplate = "<html><head></head><body><h1>{{ .StatusCode }}</h1><p>{{ .Message }}<p></body></html>"
)

var box = packr.New("templates", "./templates")

var (
	logdest = flag.String("logdest", "/dev/null", "log messages destination")
	addr    = flag.String("addr", ":8080", "default listening address")
)

var T = template.Must(template.New("_").Parse("{{ template \"html\" }}"))

type Error struct {
	StatusCode  int    `json:"statusCode"`
	Message     string `json:"message"`
	URI         string `json:"-"`
	Namespace   string `json:"-"`
	IngressName string `json:"-"`
	ServiceName string `json:"-"`
	ServicePort int    `json:"-"`
}

func main() {
	flag.Parse()

	log.Printf("Starting default-backend server on %s...", *addr)

	if len(*logdest) > 0 && *logdest != "/dev/null" {
		f, err := os.OpenFile(*logdest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)
		defer f.Close()
	} else {
		log.SetOutput(ioutil.Discard)
	}

	box.Walk(func(path string, f file.File) error {
		if strings.HasSuffix(path, ".tmpl") {
			name := path[0 : len(path)-5]
			if _, err := T.New(name).Parse(f.String()); err != nil {
				panic(err)
			}
		}
		return nil
	})
	if t := T.Lookup("html"); t == nil {
		T.New("html").Parse(defaultHTMLTemplate)
	}

	http.HandleFunc("/", errorHandler())

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(*addr, nil)
}

func NewError(r *http.Request) *Error {
	errCode := r.Header.Get(CodeHeader)
	code, err := strconv.Atoi(errCode)
	if err != nil || len(errCode) == 0 {
		code = 404
		log.Printf("unexpected error reading return code: %v. Using %v", err, code)
	}

	servicePort, err := strconv.Atoi(r.Header.Get(ServicePort))
	if err != nil {
		log.Printf("unexpected service port %s", r.Header.Get(ServicePort))
	}

	return &Error{
		StatusCode:  code,
		Message:     http.StatusText(code),
		URI:         r.Header.Get(OriginalURI),
		Namespace:   r.Header.Get(Namespace),
		IngressName: r.Header.Get(IngressName),
		ServiceName: r.Header.Get(ServiceName),
		ServicePort: servicePort,
	}
}

func errorHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		e := NewError(r)
		ext := "html"

		format := r.Header.Get(FormatHeader)
		if format == "" {
			format = "text/html"
			log.Printf("format not specified. Using %v", format)
		}

		cext, err := mime.ExtensionsByType(format)
		if err != nil {
			log.Printf("unexpected error reading media type extension: %v. Using %v", err, ext)
		} else if len(cext) == 0 {
			log.Printf("couldn't get media type extension. Using %v", ext)
		} else {
			ext = cext[0]
		}
		if strings.HasPrefix(ext, ".") {
			ext = ext[1:]
		}
		w.Header().Set(ContentType, "text/html")
		w.WriteHeader(e.StatusCode)

		if ext == "json" {
			w.Header().Set(ContentType, "application/json")
			resp, _ := json.Marshal(e)
			w.Write(resp)
		} else {
			t := T.Lookup(ext)
			if t == nil {
				t = T.Lookup("html")
			}

			if err := t.Execute(w, e); err != nil {
				log.Printf("%s", err)
			}
		}

		duration := time.Now().Sub(start).Seconds()

		proto := strconv.Itoa(r.ProtoMajor)
		proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))

		requestCount.WithLabelValues(proto).Inc()
		requestDuration.WithLabelValues(proto).Observe(duration)
	}
}
