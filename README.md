# Frames: Image Overlay Service

-----

Frames is a lightweight Go service that allows you to overlay a "frame" image onto a "main" image. It's designed to be simple, efficient, and capable of directly handling Google Drive shared image links by converting them into direct download URLs.

## Table of Contents

  * [Features](https://www.google.com/search?q=%23features)
  * [How It Works](https://www.google.com/search?q=%23how-it-works)
  * [API Endpoint](https://www.google.com/search?q=%23api-endpoint)
  * [Setup and Installation](https://www.google.com/search?q=%23setup-and-installation)
      * [Prerequisites](https://www.google.com/search?q=%23prerequisites)
      * [Clone the Repository](https://www.google.com/search?q=%23clone-the-repository)
      * [Run the Server](https://www.google.com/search?q=%23run-the-server)
  * [Usage Examples](https://www.google.com/search?q=%23usage-examples)
      * [Example with Regular URLs](https://www.google.com/search?q=%23example-with-regular-urls)
      * [Example with Google Drive Links](https://www.google.com/search?q=%23example-with-google-drive-links)
  * [Error Handling](https://www.google.com/search?q=%23error-handling)
  * [Contributing](https://www.google.com/search?q=%23contributing)
  * [License](https://www.google.com/search?q=%23license)

-----

## Features

  * **Image Overlay:** Seamlessly overlays a specified "frame" image onto a "main" image.
  * **Automatic Resizing & Cropping:** The main image is automatically resized and cropped to precisely fit the dimensions of the frame, ensuring a perfect fit.
  * **Direct URL Support:** Accepts direct image URLs for both the frame and the main image.
  * **Google Drive Link Conversion:** Smartly converts Google Drive shared links (e.g., `drive.google.com/file/d/...`) into direct download URLs, making it easy to use images stored in Google Drive.
  * **PNG Output:** Returns the processed image in PNG format.
  * **Simple HTTP Service:** Exposes its functionality via a straightforward RESTful HTTP endpoint.
  * **Robust Error Handling:** Provides clear error messages for invalid parameters, download failures, or image processing issues.

-----

## How It Works

The service operates as an HTTP server that listens for requests on a specific endpoint.

1.  **Request Reception:** When a `GET` request is received at `/imageprocessing/overlay`, it expects two query parameters: `frame_url` and `image_url`.
2.  **Parameter Validation:** It first validates that both URLs are provided and are well-formed.
3.  **Image Download:** For each URL:
      * It checks if the URL is a Google Drive shared link. If so, it extracts the file ID and reconstructs the URL into a direct download link (e.g., `https://drive.google.com/uc?export=download&id=...`).
      * It then downloads the image content from the (potentially modified) URL.
4.  **Image Processing:**
      * The `main_image` is resized and cropped to exactly match the dimensions of the `frame_image`. This ensures the main image fits neatly within the frame's boundaries without distortion or overflow.
      * A new blank image is created.
      * The processed `main_image` is drawn onto this new image.
      * Finally, the `frame_image` (resized to match the new image's dimensions if necessary) is drawn on top using `draw.Over` composition, effectively overlaying it.
5.  **Response:** The resulting composite image is encoded as a PNG and sent back as the HTTP response.

-----

## API Endpoint

The server runs on port `8080`.

**Endpoint:** `/imageprocessing/overlay`
**Method:** `GET`

**Query Parameters:**

  * `frame_url` (string, required): The URL of the image to be used as the overlay frame.
  * `image_url` (string, required): The URL of the main image onto which the frame will be overlaid.

**Response:**

  * **Success (200 OK):** The composite image in PNG format.
  * **Error (400 Bad Request):** If `frame_url` or `image_url` are missing or invalid.
  * **Error (500 Internal Server Error):** If there are issues downloading, decoding, or processing the images.

-----

## Setup and Installation

### Prerequisites

  * **Go (Golang)**: Version 1.15 or newer.
  * **Git**: For cloning the repository.

### Clone the Repository

First, clone the project to your local machine:

```bash
git clone https://github.com/sidsrbh/frames.git
cd frames
```

### Run the Server

1.  **Download Go Modules:**
    This command will fetch the necessary Go modules, including `github.com/disintegration/imaging` and standard Go libraries.

    ```bash
    go mod tidy
    ```

2.  **Run the application:**

    ```bash
    go run main.go imageprocessing/imageprocessing.go
    ```

    Alternatively, you can build an executable and run it:

    ```bash
    go build -o frames-server .
    ./frames-server
    ```

    The server will start and listen on port `8080`. You'll see a log message: `Starting server on port :8080...`

-----

## Usage Examples

Once the server is running, you can send `GET` requests to the `/imageprocessing/overlay` endpoint.

### Example with Regular URLs

You can use any publicly accessible image URLs.
Let's assume:

  * `frame_url`: `https://example.com/path/to/my_frame.png`
  * `image_url`: `https://example.com/path/to/my_main_image.jpg`

You can test this using `curl` or by pasting the URL directly into your web browser (for viewing the image).

```bash
curl -o output.png "http://localhost:8080/imageprocessing/overlay?frame_url=https://i.imgur.com/your_frame_image.png&image_url=https://i.imgur.com/your_main_image.jpg"
```

*(Replace `https://i.imgur.com/...` with actual image URLs you want to use for testing)*

### Example with Google Drive Links

**Important:** For Google Drive links to work, they must be publicly shared (or "Anyone with the link" access).

Let's assume:

  * `frame_url`: A Google Drive link to your frame image (e.g., `https://drive.google.com/file/d/SOME_FRAME_FILE_ID/view?usp=sharing`)
  * `image_url`: A Google Drive link to your main image (e.g., `https://drive.google.com/file/d/SOME_IMAGE_FILE_ID/view?usp=sharing`)

<!-- end list -->

```bash
curl -o output_gdrive.png "http://localhost:8080/imageprocessing/overlay?frame_url=https://drive.google.com/file/d/YOUR_FRAME_FILE_ID/view?usp=sharing&image_url=https://drive.google.com/file/d/YOUR_IMAGE_FILE_ID/view?usp=sharing"
```

*(Replace `YOUR_FRAME_FILE_ID` and `YOUR_IMAGE_FILE_ID` with actual Google Drive file IDs)*

The server will automatically detect the Google Drive URLs and convert them to direct download links before processing.

-----

## Error Handling

The service provides informative error messages:

  * **`frame_url and image_url are required`**: If one or both query parameters are missing.
  * **`invalid frame_url` / `invalid image_url`**: If the provided URL string is not a valid URL format.
  * **`failed to download image: non-200 response code`**: If the image could not be downloaded (e.g., 404 Not Found, 403 Forbidden).
  * **`failed to download image: ...`**: If there's a network error during download.
  * **`failed to decode image`**: If the downloaded file is not a valid image format.
  * **`failed to encode the resulting image`**: If there's an issue generating the final PNG output.

-----

## Contributing

Contributions are very welcome\! If you have suggestions, bug reports, or want to contribute to the code, please feel free to:

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature-name`).
3.  Make your changes.
4.  Commit your changes (`git commit -m 'feat: Add new feature'`).
5.  Push to the branch (`git push origin feature/your-feature-name`).
6.  Open a Pull Request.

-----

## License

This project is open-sourced under the [MIT License](https://www.google.com/search?q=LICENSE).

-----