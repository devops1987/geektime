package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"net/http/pprof"

	"geektime/httpserver/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// var startTime = time.Now()

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	log.Println("Starting http server...")
	/*
		http.HandleFunc("/", rootHandler)
		http.HandleFunc("/healthz", healthz)
		http.Handle("/metrics", promhttp.Handler())
	*/
	metrics.Register()

	// err := http.ListenAndServe(":80", nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// err := http.ListenAndServe(":80", mux)
	srv := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Startup Failed: %+v", err)
		}
	}()

	fmt.Println("awaiting signal")
	<-sigs
	// time.Sleep(5 * time.Second)
	fmt.Println("exiting")

	if err := srv.Shutdown(nil); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering health check...")
	/* duration := time.Now().Sub(startTime)
	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
	*/
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering root handler...")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))

	// 将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		// fmt.Println(k, v[0])
		w.Header().Set(k, v[0])
		fmt.Printf("Set request header to response header! key: %v; value: %v\n", k, v[0])
	}

	// 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	if env := os.Getenv("VERSION"); env != "" {
		w.Header().Set("VERSION", env)
	} else {
		w.Header().Set("VERSION", "null")
	}

	// 将 header 作为响应内容
	io.WriteString(w, "=================== Details of the http response header ============\n")
	for k, v := range w.Header() {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v[0]))
	}

	// Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	log.Printf("Request info, client ip: %v, response code: %v", strings.Split(r.RemoteAddr, ":")[0], http.StatusOK)
}
