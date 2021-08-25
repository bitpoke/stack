package main

import (
	"embed"
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

//go:embed templates
var box embed.FS

var (
	logdest = flag.String("logdest", "/dev/null", "log messages destination")
	addr    = flag.String("addr", ":8080", "default listening address")
)

var T *template.Template

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
	var err error
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

	T, err = template.ParseFS(box, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	if t := T.Lookup("html.tmpl"); t == nil {
		if _, err := T.New("html.tmpl").Parse(defaultHTMLTemplate); err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/", errorHandler())

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
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
		ext = strings.TrimPrefix(ext, ".")

		w.Header().Set(ContentType, "text/html")
		w.WriteHeader(e.StatusCode)

		if ext == "json" {
			w.Header().Set(ContentType, "application/json")
			resp, _ := json.Marshal(e)
			if _, err := w.Write(resp); err != nil {
				panic(err)
			}
		} else {
			t := T.Lookup(ext + ".tmpl")
			if t == nil {
				t = T.Lookup("html.tmpl")
			}

			if err := t.Execute(w, e); err != nil {
				log.Printf("%s", err)
			}
		}

		duration := time.Since(start).Seconds()

		proto := strconv.Itoa(r.ProtoMajor)
		proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))

		requestCount.WithLabelValues(proto).Inc()
		requestDuration.WithLabelValues(proto).Observe(duration)
	}
}
