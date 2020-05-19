package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Response struct
type Response struct {
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/profiles", handler).Methods("POST")

	fmt.Printf("Running app in: %v\n", PORT)
	http.ListenAndServe(":"+PORT, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	res := &Response{}
	maxSize := int64(1024000)

	if err := r.ParseMultipartForm(maxSize); err != nil {
		res.Errors = fmt.Sprintf("Image too large. Max size: %v", maxSize)
		JSON(w, 400, res)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		res.Errors = err.Error()
		res.Message = "Could not get uploaded file"

		JSON(w, 400, res)
		return
	}
	defer file.Close()

	aws := &Amazon{
		Region:    AWS_S3_REGION,
		Bucket:    AWS_S3_BUCKET,
		AccessID:  AWS_ACCESS_KEY_ID,
		AccessKey: AWS_SECRET_ACCESS_KEY,
	}

	fileName, err := aws.UploadFileS3(file, fileHeader)
	if err != nil {
		res.Message = "Could not upload file"
		res.Errors = err.Error()
		JSON(w, 400, res)
		return
	}

	res.Data = fileName
	JSON(w, 200, res)
}

// JSON func
func JSON(w http.ResponseWriter, status int, res *Response) {
	w.WriteHeader(status)
	data, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Marshal Error: %v", err.Error())
	}

	w.Write([]byte(data))
}
