package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"xdp-blocker/controllers"

	"github.com/dropbox/goebpf"
	"github.com/gin-gonic/gin"
)

var iface = flag.String("iface", "", "Interface to bind XDP program to")

func main() {
	flag.Parse()
	if *iface == "" {
		log.Fatal("-iface is required.")
	}

	// Load XDP Into App
	bpf := goebpf.NewDefaultEbpfSystem()
	err := bpf.LoadElf("bpf/xdp_drop.elf")
	if err != nil {
		log.Fatalf("LoadELF() failed: %s", err)
	}
	matches := bpf.GetMapByName("matches")
	if matches == nil {
		log.Fatalf("eBPF map 'matches' not found\n")
	}
	blacklist := bpf.GetMapByName("blacklist")
	if blacklist == nil {
		log.Fatalf("eBPF map 'blacklist' not found\n")
	}
	xdp := bpf.GetProgramByName("firewall")
	if xdp == nil {
		log.Fatalln("Program 'firewall' not found in Program")
	}
	err = xdp.Load()
	if err != nil {
		fmt.Printf("xdp.Attach(): %v", err)
	}
	err = xdp.Attach(*iface)
	if err != nil {
		log.Fatalf("Error attaching to Interface: %s", err)
	}
	fmt.Println("XDP Program Loaded successfuly into the Kernel.")

	// Starting Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", controllers.Welcome)
	v1 := router.Group("/v1")
	{
		v1.GET("/", controllers.Welcome)
		v1.POST("/add-ip", controllers.Block(blacklist))
		v1.POST("/remove-ip", controllers.UnBlock(blacklist))
	}

	defer xdp.Detach()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	log.Println("\nDetaching program and exit")
}

// THIS DOESN"T WORK
// TODO:
// Loading XDP logic in a seperate Function rather than the main Function

// func LoadXDP() (*goebpf.Map, error) {
// 	bpf := goebpf.NewDefaultEbpfSystem()
// 	// Load .ELF files compiled by clang/llvm
// 	err := bpf.LoadElf("bpf/xdp_drop.elf")
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get eBPF maps
// 	matches := bpf.GetMapByName("matches")
// 	if matches == nil {
// 		return nil, err
// 	}

// 	blacklist := bpf.GetMapByName("blacklist")
// 	if blacklist == nil {
// 		return nil, err
// 	}

// 	xdp := bpf.GetProgramByName("firewall")
// 	if xdp == nil {
// 		return nil, err
// 	}
// 	err = xdp.Load()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Attach to the interface
// 	err = xdp.Attach("enp0s3")
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer xdp.Detach()

// 	// Add CTRL-C handler with channels
// 	log.Println("XDP Program Loaded successfuly into the Kernel.")
// 	return &blacklist, nil
// }
