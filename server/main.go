package main

import (
	"flag"
	gogpt "github.com/sashabaranov/go-gpt3"
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const HELP = ` Flags:
--port, -p
        Set the application port
--cert-full, -cf
        Path to the full chain SSL certificate
--cert-priv, -cp
        Path to the private SSL certificate
--help, -h
        Print this message
`

var (
	PORT    string
	API_KEY string
)

var GptClient *gogpt.Client

func init() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&PORT, "port", "80", "Application Port")
	flags.StringVar(&PORT, "p", "80", "Application Port")
	flags.StringVar(&API_KEY, "key", "80", "OpenAI API key")
	flags.StringVar(&API_KEY, "k", "80", "OpenAI API key")
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	GptClient = gogpt.NewClient(API_KEY)
}

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	lim, exists := i.ips[ip]
	if !exists {
		lim = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = lim
	}
	return lim
}

func (i *IPRateLimiter) CheckIP(ip string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	lim := i.getLimiter(ip)
	return lim.Allow()
}

func IPLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if !limiter.CheckIP(ip) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests),
				http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func FileServerFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") && r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "max-age=432000")
		next.ServeHTTP(w, r)
	})
}

func CheckMethod(t string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t != r.Method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var limiter = NewIPRateLimiter(0.075, 7)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("../contnent"))
	mux.Handle("/", FileServerFilter(fs))
	mux.Handle("/req", CheckMethod("POST", IPLimit(http.HandlerFunc(CompletionRequest))))
	s := &http.Server{
		Addr:           ":" + PORT,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening on port: ", PORT)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}
