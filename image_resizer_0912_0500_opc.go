// 代码生成时间: 2025-09-12 05:00:28
package main

import (
    "fmt"
    "image"
    "image/jpeg"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/gorilla/mux"
)

// ImageResizer struct defines the parameters needed for resizing images.
type ImageResizer struct {
    Width, Height int
}

// NewImageResizer creates a new ImageResizer with specified width and height.
func NewImageResizer(width, height int) *ImageResizer {
    return &ImageResizer{
        Width:  width,
        Height: height,
    }
}

// ResizeImage resizes the given image to the specified dimensions.
func (r *ImageResizer) ResizeImage(img image.Image) image.Image {
    // Create a new image with the desired dimensions.
    resImg := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
    
    // Set the new image's bounds based on the aspect ratio of the original image.
    w, h := img.Bounds().Dx(), img.Bounds().Dy()
    srcRatio := w / float64(h)
    dstRatio := float64(r.Width) / float64(r.Height)
    
    if srcRatio > dstRatio {
        // Fit the image to the width.
        newHeight := int(float64(r.Width) / srcRatio)
        resImg = image.NewRGBA(image.Rect(0, 0, r.Width, newHeight))
    } else {
        // Fit the image to the height.
        newWidth := int(float64(r.Height) * srcRatio)
        resImg = image.NewRGBA(image.Rect(0, 0, newWidth, r.Height))
    }
    
    // Draw the original image into the new image.
    resImg.Draw(resImg.Bounds(), img)
    return resImg
}

// HandleResize handles the HTTP request for resizing images.
func HandleResize(resizer *ImageResizer, w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    directory := vars["directory"]
    widthStr := vars["width"]
    heightStr := vars["height"]
    width, _ := strconv.Atoi(widthStr)
    height, _ := strconv.Atoi(heightStr)
    
    // Create an ImageResizer with the given width and height.
    resizer.Width = width
    resizer.Height = height
    
    // Open the directory and list all files.
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Loop through the files and resize each image.
    for _, file := range files {
        filePath := filepath.Join(directory, file.Name())
        imgFile, err := os.Open(filePath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer imgFile.Close()
        
        img, _, err := image.Decode(imgFile)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Resize the image.
        resizedImg := resizer.ResizeImage(img)
        
        // Save the resized image.
        outPath := filepath.Join(directory, "resized_" + file.Name())
        outFile, err := os.Create(outPath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer outFile.Close()
        
        // Encode the resized image to JPEG.
        err = jpeg.Encode(outFile, resizedImg, nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    fmt.Fprintf(w, "All images have been resized successfully.")
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/resize/{directory:[^/]+}/{width:[0-9]+}/{height:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
        handler := NewImageResizer(0, 0)
        HandleResize(handler, w, r)
    }, nil)
    
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", router)
}
