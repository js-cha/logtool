package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/js-cha/logtool/internal/logtool"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("Duration: %s\n", time.Since(start))
	}()
	path := flag.String("path", "./logs", "the file path or directory")
	mode := flag.String("mode", "concurrent", "sync or concurrent")
	workers := flag.Int("workers", 4, "number of workers")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("\nshutdown signal received")
		cancel()
	}()

	fileInfo, err := os.Stat(*path)
	if err != nil {
		log.Fatalf("processing failed: %v", err)
	}

	if fileInfo.IsDir() {
		fmt.Println("Mode:", *mode)
		switch *mode {
		case "concurrent":
			if *workers < 1 {
				log.Fatal("worker count needs to be greater than zero")
			}
			fmt.Println("Number of workers: ", *workers)
			counts, err := logtool.ProcessDirConcurrent(ctx, *path, *workers)
			if err != nil {
				fmt.Printf("processing failed: %v\n", err)
			}
			logtool.PrintLogLevels(counts)
		case "sync":
			counts, err := logtool.ProcessDirSync(ctx, *path)
			if err != nil {
				fmt.Printf("processing failed: %v\n", err)
			}
			logtool.PrintLogLevels(counts)
		default:
			log.Fatal("invalid mode flag, use: 'concurrent' or 'sync'")
		}
	} else {
		counts, err := logtool.ProcessLogFile(ctx, *path)
		if err != nil {
			fmt.Printf("processing failed: %v\n", err)
		}
		logtool.PrintLogLevels(counts)
	}
}
