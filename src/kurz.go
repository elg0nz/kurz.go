package main

import (
	"github.com/joho/godotenv"
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fs111/simpleconfig"
	"github.com/op/go-logging"
	godis "github.com/simonz05/godis/redis"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	// special key in redis, that is our global counter
	COUNTER = "__counter__"
	HTTP    = "http"
)

var (
	redis            *godis.Client
	config           *simpleconfig.Config
	filenotfound     string
	allowedHostNames []string
)

var log = logging.MustGetLogger("kurz")

type KurzUrl struct {
	Key          string
	ShortUrl     string
	LongUrl      string
	CreationDate int64
	Clicks       int64
}

// Converts the KurzUrl to JSON.
func (k KurzUrl) Json() []byte {
	b, _ := json.Marshal(k)
	return b
}

// Creates a new KurzUrl instance. The Given key, shorturl and longurl will
// be used. Clicks will be set to 0 and CreationDate to time.Nanoseconds()
func NewKurzUrl(key, shorturl, longurl string) *KurzUrl {
	kurl := new(KurzUrl)
	kurl.CreationDate = time.Now().UnixNano()
	kurl.Key = key
	kurl.LongUrl = longurl
	kurl.ShortUrl = shorturl
	kurl.Clicks = 0
	return kurl
}

// stores a new KurzUrl for the given key, shorturl and longurl. Existing
// ones with the same url will be overwritten
func store(key, shorturl, longurl string) *KurzUrl {
	log.Info("Storing key: %s, short: %s, long: %s", key, shorturl, longurl)
	kurl := NewKurzUrl(key, shorturl, longurl)
	go redis.Hset(kurl.Key, "LongUrl", kurl.LongUrl)
	go redis.Hset(kurl.Key, "ShortUrl", kurl.ShortUrl)
	go redis.Hset(kurl.Key, "CreationDate", kurl.CreationDate)
	go redis.Hset(kurl.Key, "Clicks", kurl.Clicks)
	return kurl
}

// loads a KurzUrl instance for the given key. If the key is
// not found, os.Error is returned.
func load(key string) (*KurzUrl, error) {
	log.Info("Loading key: %s", key)
	if ok, _ := redis.Hexists(key, "ShortUrl"); ok {
		kurl := new(KurzUrl)
		kurl.Key = key
		reply, _ := redis.Hmget(key, "LongUrl", "ShortUrl", "CreationDate", "Clicks")
		kurl.LongUrl, kurl.ShortUrl, kurl.CreationDate, kurl.Clicks =
			reply.Elems[0].Elem.String(), reply.Elems[1].Elem.String(),
			reply.Elems[2].Elem.Int64(), reply.Elems[3].Elem.Int64()
		return kurl, nil
	}
	return nil, errors.New("unknown key: " + key)
}

func fileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// function to display the info about a KurzUrl given by it's Key
func info(w http.ResponseWriter, r *http.Request) {
	short := mux.Vars(r)["short"]
	log.Info("Display info for: %s", short)
	if strings.HasSuffix(short, "+") {
		short = strings.Replace(short, "+", "", 1)
	}

	kurl, err := load(short)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(kurl.Json())
		io.WriteString(w, "\n")
	} else {
		http.Redirect(w, r, filenotfound, http.StatusNotFound)
	}
}

// function to resolve a shorturl and redirect
func resolve(w http.ResponseWriter, r *http.Request) {
	short := mux.Vars(r)["short"]
	log.Info("Resolving: %s", short)
	kurl, err := load(short)
	if err == nil {
		go redis.Hincrby(kurl.Key, "Clicks", 1)
		http.Redirect(w, r, kurl.LongUrl, http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, filenotfound, http.StatusMovedPermanently)
	}
}

// Determines if the string rawurl is a valid URL to be stored.
func isValidUrl(rawurl string) (u *url.URL, err error) {
	if len(rawurl) == 0 {
		return nil, errors.New("empty url")
	}
	// XXX this needs some love...
	if !strings.HasPrefix(rawurl, HTTP) {
		rawurl = fmt.Sprintf("%s://%s", HTTP, rawurl)
	}
	return url.Parse(rawurl)
}

// function to shorten and store a url
func shorten(w http.ResponseWriter, r *http.Request) {
	host := os.Getenv("HOSTNAME")
	leUrl := r.FormValue("url")
	log.Info("Storing url: %s with host %s", string(leUrl), host)
	theUrl, err := isValidUrl(string(leUrl))
	if err == nil {
		if !check_hostname(allowedHostNames, theUrl) {
			log.Debug("Host name not allowed.")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Hostname not allowed"))
			return
		}

		ctr, _ := redis.Incr(COUNTER)
		encoded := Encode(ctr)
		location := fmt.Sprintf("%s://%s/%s",
			os.Getenv("PROTO"), host, encoded)
		store(encoded, location, theUrl.String())

		home := r.FormValue("home")
		if home != "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			// redirect to the info page
			http.Redirect(w, r, location+"+", http.StatusMovedPermanently)
		}
	} else {
		http.Redirect(w, r, filenotfound, http.StatusNotFound)
	}
}

//Returns a json array with information about the last shortened urls. If data
// is a valid integer, that's the amount of data it will return, otherwise
// a maximum of 10 entries will be returned.
func latest(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)["data"]
	howmany, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		howmany = 10
	}
	c, _ := redis.Get(COUNTER)

	last := c.Int64()
	upTo := (last - howmany)

	w.Header().Set("Content-Type", "application/json")

	var kurls = []*KurzUrl{}

	for i := last; i > upTo && i > 0; i -= 1 {
		kurl, err := load(Encode(i))
		if err == nil {
			kurls = append(kurls, kurl)
		}
	}
	s, _ := json.Marshal(kurls)
	w.Write(s)
}

func static(w http.ResponseWriter, r *http.Request) {
	fname := mux.Vars(r)["fileName"]
	// empty means, we want to serve the index file. Due to a bug in http.serveFile
	// the file cannot be called index.html, anything else is fine.
	if fname == "" {
		fname = "index.htm"
	}
	log.Info("Serving static: %s", fname)

	staticDir := os.Getenv("STATIC_DIR")
	staticFile := path.Join(staticDir, fname)
	if fileExists(staticFile) {
		http.ServeFile(w, r, staticFile)
	} else {
		log.Info("Static file %s does not exist.", staticFile)
	}
}

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	flag.Parse()
	path := flag.Arg(0)

	var format = logging.MustStringFormatter("%{level} %{message}")
	logging.SetFormatter(format)

	config, _ = simpleconfig.NewConfig(path)

	var loglevel logging.Level
	switch os.Getenv("LOGLEVEL") {
	case "debug":
		loglevel = logging.DEBUG
	case "critical":
		loglevel = logging.CRITICAL
	case "error":
		loglevel = logging.ERROR
	case "info":
		loglevel = logging.INFO
	case "notice":
		loglevel = logging.NOTICE
	case "warning":
		loglevel = logging.WARNING
  }

	logging.SetLevel(loglevel, "info")

	log.Debug("Starting kurz")

	filenotfound = os.Getenv("FILE_NOT_FOUND_URL")

  // TODO: Restore this feature
  allowedHostNames = []string{}

	log.Debug("Will short URLs from following sites: %q", allowedHostNames)

	redis = redis_instance()

	router := mux.NewRouter()
	router.HandleFunc("/shorten/{url:(.*$)}", shorten)

	router.HandleFunc("/{short:([a-zA-Z0-9]+$)}", resolve)
	router.HandleFunc("/{short:([a-zA-Z0-9]+)\\+$}", info)
	router.HandleFunc("/info/{short:[a-zA-Z0-9]+}", info)
	router.HandleFunc("/latest/{data:[0-9]+}", latest)

	router.HandleFunc("/{fileName:(.*$)}", static)

	listen := os.Getenv("LISTEN_ADDRESS")
	port := os.Getenv("LISTEN_PORT")
	log.Info("Listening %s:%s", listen, port)
	s := &http.Server{
		Addr:    listen + ":" + port,
		Handler: router,
	}
	s.ListenAndServe()
}
