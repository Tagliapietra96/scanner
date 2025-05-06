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
- [APIs and Data Structures](#apis-and-data-structures)
- [How it Works](#how-it-works)
- [Use Cases](#use-cases)
- [Support the Project](#support-the-project)
- [License](#license)

## 🌟 Introduction

The `scanner` package provides a delightful way to traverse directory structures recursively with powerful filtering options. It leverages Go's concurrency model to scan directories significantly faster than standard library functions like `filepath.Walk` and `filepath.WalkDir`. 

Think of it as your personal file system explorer with superpowers - it's like having a tiny, efficient robot that zooms through your directories while you sip coffee! ☕

## ✨ Features

- 🚄 **Blazing Fast Performance**: Concurrent design for maximum speed
- 🔄 **Flexible API**: Both synchronous and asynchronous scanning options
- 🧩 **Rich Filtering System**: Built-in filters for common use cases
- 📏 **Configurable Depth**: Control how deep you want to go in the directory tree
- 🔌 **Platform-Aware**: Special handling for hidden files on different operating systems
- 💪 **Robust Error Handling**: Graceful recovery from permission errors and other issues
- 🧠 **Smart Resource Management**: Optimized for CPU utilization

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
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Let's scan the current directory with no depth limit (-1 means scan everything)
    // and no filters (nil means include everything)
    results, err := scanner.ScanSync(".", -1, nil)
    
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
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Let's find all directories (excluding files) within 2 levels of depth
    directories, err := scanner.ScanSync("/path/to/scan", 2, scanner.FilterDir)
    
    if err != nil {
        log.Fatalf("Directory scan failed: %v", err)
    }
    
    fmt.Println("Directories found:")
    for _, dir := range directories {
        fmt.Println(dir)
    }
}
```

### Advanced Filtering Example

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Create custom filter: find non-hidden Go files larger than 1KB
    customFilter := func(path string, entry os.DirEntry) bool {
        // Skip hidden files
        if scanner.IsHidden(path) {
            return false
        }
        
        // Check if it's a Go file
        extFilter := scanner.FilterByExtension(".go")
        if !extFilter(path, entry) {
            return false
        }
        
        // Check file size (greater than 1KB)
        sizeFilter := scanner.FilterBySize(1024, ">")
        return sizeFilter(path, entry)
    }
    
    // Scan with our custom filter
    goFiles, err := scanner.ScanSync("./src", -1, customFilter)
    
    if err != nil {
        log.Fatalf("Scan failed with error: %v", err)
    }
    
    fmt.Printf("Found %d Go files larger than 1KB:\n", len(goFiles))
    for _, file := range goFiles {
        fmt.Println(file)
    }
}
```

### Asynchronous Scanning

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/Tagliapietra96/scanner"
)

func main() {
    // Create channels for results and errors
    resultChan := make(chan string)
    errorChan := make(chan error)
    
    // Start scanning asynchronously
    scanner.Scan("/path/to/scan", -1, scanner.FilterFile, resultChan, errorChan)
    
    // Process results and errors as they come in
    fileCount := 0
    
    // This loop will exit when both channels are closed by the scanner
    for {
        select {
        case path, ok := <-resultChan:
            if !ok {
                // Channel is closed, no more results
                fmt.Printf("Scan complete! Found %d files.\n", fileCount)
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

## 📚 APIs and Data Structures

### Core Functions

- **`Scan(root string, maxDepth int, filter func, resultChan, errorChan)`**: Asynchronously scans directories
- **`ScanSync(root string, maxDepth int, filter func) ([]string, error)`**: Synchronously scans directories

### Filter Functions

- **`FilterDir`**: Matches only directories
- **`FilterFile`**: Matches only files (non-directories)
- **`FilterHidden`**: Matches hidden files/directories
- **`FilterRegular`**: Matches regular files
- **`FilterSymlink`**: Matches symbolic links
- **`FilterDevice`**: Matches device files
- **`FilterNamedPipe`**: Matches named pipes
- **`FilterSocket`**: Matches socket files
- **`FilterCharDev`**: Matches character devices
- **`FilterByExtension(ext)`**: Returns filter matching files with specified extension
- **`FilterBySize(size, operator)`**: Returns filter matching files based on size comparisons

### Platform-Specific Functions

- **`IsHidden(path)`**: Cross-platform detection of hidden files/directories

## ⚙️ How it Works

The scanner package employs a beautifully orchestrated concurrent design to traverse directories efficiently:

1. The `scan` function is the core workhorse that recursively traverses directories
2. It uses goroutines for concurrent scanning, with a semaphore to limit concurrency
3. Each directory entry is evaluated against filter functions
4. Matching entries are sent to a result channel
5. Errors encountered are sent to an error channel
6. The function respects the specified maximum depth

The package optimizes CPU utilization by limiting the number of concurrent operations based on the available CPU cores. This prevents overwhelming the system while maximizing throughput.

Cross-platform support is achieved through build tags that provide platform-specific implementations for functions like `IsHidden`.

## 🎯 Use Cases

- **Web Servers**: Quickly scan for static assets
- **Build Tools**: Find source files for compilation
- **Backup Systems**: Efficiently list files for backup
- **Search Utilities**: Build fast file search applications
- **Content Management**: Catalog media files
- **DevOps Tools**: Find configuration files across directories
- **Package Managers**: Scan for dependencies
- **Data Analysis**: Process collections of data files
- **Testing**: Scan for test files to execute

## 💝 Support the Project

Love this package? Here's how you can support its development:

- ⭐ Star the repository on GitHub
- 🐛 Report issues and contribute bug fixes
- 🌟 Contribute enhancements and new features
- 📚 Improve documentation
- 🌐 Share with your network and fellow developers

## 📜 License

This package is available under the [MIT License](https://github.com/Tagliapietra96/scanner/blob/main/LICENSE). Feel free to use it in your projects, modify it, and share the love! Just remember to include the original license text.

Happy scanning! May your traversals be swift and your filters precise! 🎉

