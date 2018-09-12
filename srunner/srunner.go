package main

import (
    "fmt"
    "sdcc/p1"
)

const defaultPort = 9999

func main() {
    // Initialize the server.
    server := p1.New()
    if server == nil {
        fmt.Println("New() returned a nil server. Exiting...")
        return
    }

    // Start the server and continue listening for client connections in the background.
    if err := server.Start(defaultPort); err != nil {
        fmt.Printf("KeyValueServer could not be started: %s\n", err)
        return
    }

    fmt.Printf("Started KeyValueServer on port %d...\n", defaultPort)

    // Block forever.
    select {}
}
