// 代码生成时间: 2025-09-04 00:18:45
package main

import (
    "flag"
    "fmt"
    "image"
    "image/jpeg"
    "os"
    "path/filepath"
    "strings"

    "github.com/gorilla/mux"
)

// Define constants for image dimensions
const (
    defaultWidth  = 800
    defaultHeight = 600
)

// ImageResizer struct to hold the image processing logic
type ImageResizer struct {
    Width  int
    Height int
}

// NewImageResizer initializes a new ImageResizer with default dimensions
func NewImageResizer() *ImageResizer {
    return &ImageResizer{
        Width:  defaultWidth,
        Height: defaultHeight,
    }
}

// Resize resizes the image to the specified dimensions
func (r *ImageResizer) Resize(img image.Image) image.Image {
    // Create a new image to draw on
    resImg := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
    // Draw the original image on the new image, resizing it as needed
    resImg.Draw(resImg.Bounds(), img, image.Pt(0, 0), draw.Src)
    return resImg
}

// ProcessFolder processes all JPEG images in the given directory
func ProcessFolder(folderPath string, width, height int) error {
    filelist, err := os.ReadDir(folderPath)
    if err != nil {
        return err
    }

    for _, file := range filelist {
        if !file.IsDir() {
            if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") {
                filePath := filepath.Join(folderPath, file.Name())
                processImage(filePath, width, height)
            }
        }
    }
    return nil
}

// processImage processes a single image file
func processImage(filePath string, width, height int) {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Printf("Error opening file: %s
", filePath)
        return
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        fmt.Printf("Error decoding image: %s
", filePath)
        return
    }

    resizer := ImageResizer{Width: width, Height: height}
    resizedImg := resizer.Resize(img)

    outputFilePath := filePath + "_resized.jpg"
    outFile, err := os.Create(outputFilePath)
    if err != nil {
        fmt.Printf("Error creating output file: %s
", outputFilePath)
        return
    }
    defer outFile.Close()

    err = jpeg.Encode(outFile, resizedImg, nil)
    if err != nil {
        fmt.Printf("Error encoding image: %s
", outputFilePath)
        return
    }
    fmt.Printf("Resized image saved: %s
", outputFilePath)
}

func main() {
    var folderPath string
    var width, height int
    flag.IntVar(&width, "width", defaultWidth, "Width of the resized images")
    flag.IntVar(&height, "height", defaultHeight, "Height of the resized images")
    flag.StringVar(&folderPath, "path", "./images", "Path to the folder containing images")
    flag.Parse()

    // Initialize the Gorilla Mux router
    router := mux.NewRouter()
    // Define a route to resize images
    router.HandleFunc("/resize", func(w http.ResponseWriter, r *http.Request) {
        // TODO: Implement resizing logic based on request parameters
        // For now, just return a message
        fmt.Fprintf(w, "Image resizing endpoint hit")
    })

    // Start the web server
    fmt.Println("Starting image resizer server on port 8080")
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        fmt.Println("Error starting server: ", err)
    }

    // Process the folder in the background
    go func() {
        err := ProcessFolder(folderPath, width, height)
        if err != nil {
            fmt.Println("Error processing folder: ", err)
        }
    }()
}