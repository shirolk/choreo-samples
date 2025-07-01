/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Start periodic logging in background
	go startPeriodicLogging()

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/greeter/greet", greet)

	serverPort := 9090
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: serverMux,
	}
	go func() {
		log.Printf("Starting HTTP Greeter on port %d\n", serverPort)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP ListenAndServe error: %v", err)
		}
		log.Println("HTTP server stopped serving new requests.")
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh // Wait for shutdown signal

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Shutdown complete.")
}

func greet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Stranger"
	}
	log.Printf("INFO: Greeting request received for name: %s from %s", name, r.RemoteAddr)
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

func startPeriodicLogging() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	requestCount := 0
	systemErrors := []string{
		"database connection timeout",
		"memory usage high",
		"disk space low",
		"network latency spike",
		"cache miss rate elevated",
	}

	for {
		select {
		case <-ticker.C:
			requestCount++
			
			// Generate different types of logs
			switch requestCount % 4 {
			case 0:
				log.Printf("DEBUG: Periodic health check - uptime: %s, goroutines: %d", 
					time.Since(time.Now().Add(-time.Duration(requestCount*30)*time.Second)), 
					rand.Intn(50)+10)
			case 1:
				log.Printf("INFO: System metrics - requests processed: %d, memory usage: %d%%", 
					requestCount*rand.Intn(100)+50, rand.Intn(30)+60)
			case 2:
				log.Printf("WARN: Performance alert - response time: %dms (threshold: 500ms)", 
					rand.Intn(400)+600)
			case 3:
				errorType := systemErrors[rand.Intn(len(systemErrors))]
				log.Printf("ERROR: System issue detected - %s (correlation_id: %d)", 
					errorType, rand.Intn(10000)+1000)
			}
		}
	}
}
