package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type GetImageResponse struct {
	ImageURL   string
	Detections *DetectionResult
}

// HealthCheck just returns true if the service is up.
func HealthCheck(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func GetImage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get the bytes for a random image
	imgSrvc := NewImageService()
	options := []string{"person", "bicycle", "car", "motorcycle", "airplane", "bus", "train", "truck", "boat", "traffic light", "fire hydrant", "stop_sign", "parking meter", "bench", "bird", "cat", "dog", "horse", "sheep", "cow", "elephant", "bear", "zebra", "giraffe", "backpack", "umbrella", "handbag", "tie", "suitcase", "frisbee", "skis", "snowboard", "sports ball", "kite", "baseball bat", "baseball glove", "skateboard", "surfboard", "tennis racket", "bottle", "wine glass", "cup", "fork", "knife", "spoon", "bowl", "banana", "apple", "sandwich", "orange", "broccoli", "carrot", "hot dog", "pizza", "donot", "cake", "chair", "couch", "potted plant", "bed", "dining table", "toilet", "tv", "laptop", "mouse", "remote", "keyboard", "cell phone", "microwave", "oven", "toaster", "sink", "refrigerator", "book", "clock", "vase", "scissors", "teddy bear", "hair dryer", "toothbrush"}
	searchTerm := fmt.Sprintf("%s %s",
		options[rand.Intn(len(options))],
		options[rand.Intn(len(options))])
	imgURL, randomImage, err := imgSrvc.RandomImage(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	detectSrvc := NewDetectService()
	results, err := detectSrvc.DetectObjects(randomImage)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Shoot the response back to the front-end
	output, err := json.Marshal(GetImageResponse{
		Detections: results,
		ImageURL:   imgURL,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func main() {

	router := httprouter.New()

	// Health Check
	router.GET("/hc", HealthCheck)
	router.GET("/", GetImage)

	log.Fatal(http.ListenAndServe(":8080", router))
}
