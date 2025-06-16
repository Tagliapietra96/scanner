![GitHub release](https://img.shields.io/github/v/release/Tagliapietra96/scanner)
![Build Status](https://github.com/Tagliapietra96/scanner/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/Tagliapietra96/scanner/path.svg)](https://pkg.go.dev/github.com/Tagliapietra96/scanner)
[![Go Report Card](https://goreportcard.com/badge/github.com/Tagliapietra96/scanner)](https://goreportcard.com/report/github.com/Tagliapietra96/scanner)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# 🔍 Scanner - Blazing Fast Directory Traversal for Go! 🚀

A high-performance, concurrent directory scanning package that will make your file system operations faster than a caffeinated squirrel on a sugar rush! 

## 📋 Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage Examples](#usage-examples)
- [Callback Functions - The Heart of Flexibility](#callback-functions---the-heart-of-flexibility)
- [APIs and Data Structures](#apis-and-data-structures)
- [How it Works](#how-it-works)
- [Use Cases](#use-cases)
- [Support the Project](#support-the-project)
- [License](#license)

## 🌟 Introduction

The `scanner` package provides a delightful way to traverse directory structures recursively with powerful callback-based filtering. It leverages Go's concurrency model to scan directories significantly faster than standard library functions like `filepath.Walk` and `filepath.WalkDir`. 

Think of it as your personal file system explorer with superpowers - it's like having a tiny, efficient robot that zooms through your directories while you sip coffee! ☕

## ✨ Features

- 🚄 **Blazing Fast Performance**: Concurrent design for maximum speed
- 🔄 **Flexible API**: Both synchronous and asynchronous scanning options
- 🎯 **Powerful Callback System**: Sophisticated filtering and processing with callback functions
- 💪 **Robust Error Handling**: Graceful recovery from permission errors and other issues
- 🧠 **Smart Resource Management**: Optimized for CPU utilization
- 🔌 **Cross-Platform**: Works seamlessly across different operating systems

## 📦 Installation

Installing the scanner package is easier than teaching a duck to swim! Just use `go get`:

```bash
go get github.com/Tagliapietra96/scanner
```

## 🚀 Usage Examples

### Basic Directory Scanning

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Let's scan the current directory with no filtering (nil means include everything)
    results, err := scanner.ScanSync(".", nil)
    
    // Always handle your errors, folks! Even the most perfect code can have a bad day.
    if err != nil {
        log.Fatalf("Oh no! Scanner got confused: %v", err)
    }
    
    // Print out what we found
    fmt.Println("Found these lovely files and directories:")
    for _, path := range results {
        fmt.Println(path)
    }
}
```

### Filtering Only Directories

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Create a callback that only includes directories
    dirFilter := func(path string, entry os.DirEntry) error {
        if !entry.IsDir() {
            return scanner.ErrDoNotInclude
        }
        return nil
    }
    
    directories, err := scanner.ScanSync("/path/to/scan", dirFilter)
    
    if err != nil {
        log.Fatalf("Directory scan failed: %v", err)
    }
    
    fmt.Println("Directories found:")
    for _, dir := range directories {
        fmt.Println(dir)
    }
}
```

### Advanced Filtering with Skip Logic

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Create sophisticated filter: find Go files, skip hidden directories entirely
    goFileFilter := func(path string, entry os.DirEntry) error {
        // Skip hidden directories entirely (don't even traverse them)
        if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") {
            return scanner.ErrSkip
        }
        
        // For files, only include Go files
        if !entry.IsDir() {
            if !strings.HasSuffix(entry.Name(), ".go") {
                return scanner.ErrDoNotInclude
            }
        }
        
        return nil
    }
    
    // Scan with our sophisticated filter
    goFiles, err := scanner.ScanSync("./src", goFileFilter)
    
    if err != nil {
        log.Fatalf("Scan failed with error: %v", err)
    }
    
    fmt.Printf("Found %d Go files (skipped hidden directories):\n", len(goFiles))
    for _, file := range goFiles {
        fmt.Println(file)
    }
}
```

### Processing Files During Scan

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    var totalSize int64
    var fileCount int
    
    // Callback that processes files and calculates statistics
    statsCallback := func(path string, entry os.DirEntry) error {
        if !entry.IsDir() {
            info, err := entry.Info()
            if err == nil {
                totalSize += info.Size()
                fileCount++
            }
        }
        return nil
    }
    
    // Scan and process files
    _, err := scanner.ScanSync(".", statsCallback)
    
    if err != nil {
        log.Fatalf("Scan failed: %v", err)
    }
    
    fmt.Printf("Statistics:\n")
    fmt.Printf("Files: %d\n", fileCount)
    fmt.Printf("Total size: %d bytes\n", totalSize)
}
```

### Asynchronous Scanning

```go
package main

import (
    "fmt"
    "os"
    "strings"
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Create channels for results and errors
    resultChan := make(chan string)
    errorChan := make(chan error)
    
    // Filter for text files only
    textFileFilter := func(path string, entry os.DirEntry) error {
        if entry.IsDir() {
            return nil // Include directories for traversal
        }
        if strings.HasSuffix(entry.Name(), ".txt") ||
           strings.HasSuffix(entry.Name(), ".md") {
            return nil // Include these files
        }
        return scanner.ErrDoNotInclude
    }
    
    // Start scanning asynchronously
    scanner.Scan("/path/to/scan", textFileFilter, resultChan, errorChan)
    
    // Process results and errors as they come in
    fileCount := 0
    
    // This loop will exit when both channels are closed by the scanner
    for {
        select {
        case path, ok := <-resultChan:
            if !ok {
                // Channel is closed, no more results
                fmt.Printf("Scan complete! Found %d text files.\n", fileCount)
                return
            }
            fileCount++
            fmt.Printf("Found: %s\n", path)
            
        case err, ok := <-errorChan:
            if !ok {
                // Error channel closed
                continue
            }
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        }
    }
}
```

## 🎯 Callback Functions - The Heart of Flexibility

The real power of the scanner package lies in its callback system! The callback function is your Swiss Army knife for directory traversal, offering incredible flexibility:

### Callback Function Signature
```go
func(path string, entry os.DirEntry) error
```

### Return Values and Their Meanings

- **`nil`**: Include this entry in results and continue traversal
- **`scanner.ErrDoNotInclude`**: Exclude from results but continue traversal (useful for directories you want to explore but not include)
- **`scanner.ErrSkip`**: Skip this entry entirely - if it's a directory, don't even traverse into it
- **Any other error**: Abort the entire scan and return this error

### Creative Uses for Callbacks

#### 1. Smart Filtering with State
```go
var largeFiles []string
callback := func(path string, entry os.DirEntry) error {
    if !entry.IsDir() {
        if info, err := entry.Info(); err == nil && info.Size() > 1048576 {
            largeFiles = append(largeFiles, path)
        }
    }
    return nil
}
```

#### 2. Real-time Processing
```go
callback := func(path string, entry os.DirEntry) error {
    if !entry.IsDir() {
        // Process file immediately during scan
        go processFile(path)
    }
    return nil
}
```

#### 3. Conditional Directory Skipping
```go
callback := func(path string, entry os.DirEntry) error {
    if entry.IsDir() && entry.Name() == "node_modules" {
        return scanner.ErrSkip // Skip node_modules entirely
    }
    return nil
}
```

## 📚 APIs and Data Structures

### Core Functions

- **`Scan(root string, callBack func(string, os.DirEntry) error, resultChan chan<- string, errorChan chan<- error)`**: Asynchronously scans directories with callback processing
- **`ScanSync(root string, callBack func(string, os.DirEntry) error) ([]string, error)`**: Synchronously scans directories with callback processing

### Error Constants

- **`ErrDoNotInclude`**: Returned by callbacks to exclude entries from results while continuing traversal
- **`ErrSkip`**: Returned by callbacks to skip entries entirely (prevents directory traversal)

## ⚙️ How it Works

The scanner package employs a beautifully orchestrated concurrent design to traverse directories efficiently:

1. The `Scan` function creates a goroutine pool with intelligent concurrency limits
2. Each directory is processed by a separate goroutine to maximize throughput
3. Your callback function is called for every single entry, giving you complete control
4. Results and errors flow through separate channels for clean separation of concerns
5. The semaphore system prevents overwhelming the system while maximizing performance

The package optimizes CPU utilization by limiting the number of concurrent operations based on the available CPU cores (NumCPU - 2). This prevents overwhelming the system while maximizing throughput.

## 🎯 Use Cases

- **Build Systems**: Find source files with complex filtering logic
- **Code Analysis Tools**: Process files while scanning with custom logic
- **Backup Systems**: Efficiently list files with size/date filtering
- **Search Utilities**: Build lightning-fast file search with real-time processing
- **Content Management**: Catalog media files with metadata extraction
- **DevOps Tools**: Find configuration files while processing them
- **Data Processing**: Stream-process large directory structures
- **Testing Frameworks**: Discover and process test files dynamically

## 💝 Support the Project

Love this package? Here's how you can support its development:

- ⭐ Star the repository on GitHub
- 🐛 Report issues and contribute bug fixes
- 🌟 Contribute enhancements and new features
- 📚 Improve documentation
- 🌐 Share with your network and fellow developers

## 📜 License

This package is available under the [MIT License](https://github.com/Tagliapietra96/scanner/blob/main/LICENSE). Feel free to use it in your projects, modify it, and share the love! Just remember to include the original license text.

Happy scanning! May your traversals be swift and your callbacks powerful! 🎉

