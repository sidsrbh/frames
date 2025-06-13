package imageprocessing

import (
	"errors"
	"image"
	"image/draw"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/disintegration/imaging"
)

// OverlayHandler handles the HTTP request to overlay a frame on an image.
func OverlayHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the query parameters
	frameURL, imageURL, err := parseAndValidateParams(r)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Download the frame and image
	frame, err := downloadImage(frameURL)
	if err != nil {
		log.Printf("Error downloading frame image: %v", err)
		http.Error(w, "failed to download frame image", http.StatusInternalServerError)
		return
	}

	mainImage, err := downloadImage(imageURL)
	if err != nil {
		log.Printf("Error downloading main image: %v", err)
		http.Error(w, "failed to download main image", http.StatusInternalServerError)
		return
	}

	// Resize and crop the main image to fit exactly within the frame
	finalImage := resizeAndCropImageToFitFrame(mainImage, frame.Bounds())

	// Resize the frame to match the final image dimensions (in case frame was not the original size)
	resizedFrame := imaging.Resize(frame, finalImage.Bounds().Dx(), finalImage.Bounds().Dy(), imaging.Lanczos)

	// Create a new image to hold the result
	result := image.NewRGBA(finalImage.Bounds())
	draw.Draw(result, finalImage.Bounds(), finalImage, image.Point{}, draw.Src)
	draw.Draw(result, resizedFrame.Bounds(), resizedFrame, image.Point{}, draw.Over)

	// Encode and return the result
	w.Header().Set("Content-Type", "image/png")
	if err := imaging.Encode(w, result, imaging.PNG); err != nil {
		log.Printf("Error encoding the result image: %v", err)
		http.Error(w, "failed to encode the resulting image", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully processed the image overlay request.")
}

// parseAndValidateParams extracts and validates the frame_url and image_url parameters.
func parseAndValidateParams(r *http.Request) (string, string, error) {
	frameURL := r.URL.Query().Get("frame_url")
	imageURL := r.URL.Query().Get("image_url")

	if frameURL == "" || imageURL == "" {
		return "", "", errors.New("frame_url and image_url are required")
	}

	// Validate that the URLs are properly formed
	if _, err := url.ParseRequestURI(frameURL); err != nil {
		return "", "", errors.New("invalid frame_url")
	}
	if _, err := url.ParseRequestURI(imageURL); err != nil {
		return "", "", errors.New("invalid image_url")
	}

	return frameURL, imageURL, nil
}

// downloadImage downloads an image from a given URL and returns it as an image.Image.
func downloadImage(urlStr string) (image.Image, error) {
	// Check if the URL is a Google Drive link and convert it to a direct download link
	if strings.Contains(urlStr, "drive.google.com") {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}

		// Extract the file ID from the Google Drive URL
		pathParts := strings.Split(parsedURL.Path, "/")
		if len(pathParts) > 3 {
			fileID := pathParts[3]
			urlStr = "https://drive.google.com/uc?export=download&id=" + fileID
		} else {
			return nil, errors.New("invalid Google Drive URL format")
		}
	}

	// Download the image from the provided URL
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download image: non-200 response code")
	}

	// Decode the image
	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// resizeAndCropImageToFitFrame resizes and crops the main image to fit exactly within the frame's dimensions.
func resizeAndCropImageToFitFrame(mainImage image.Image, frameBounds image.Rectangle) image.Image {
	// Resize the main image to fit within the frame bounds, preserving aspect ratio
	resizedImage := imaging.Fit(mainImage, frameBounds.Dx(), frameBounds.Dy(), imaging.Lanczos)

	// Calculate the new bounds after resizing
	resizedBounds := resizedImage.Bounds()

	// If resized image dimensions do not match the frame dimensions exactly, crop the excess parts
	if resizedBounds.Dx() > frameBounds.Dx() || resizedBounds.Dy() > frameBounds.Dy() {
		cropRect := image.Rect(0, 0, frameBounds.Dx(), frameBounds.Dy())
		return imaging.Crop(resizedImage, cropRect)
	}

	// If no cropping is needed, return the resized image
	return resizedImage
}
