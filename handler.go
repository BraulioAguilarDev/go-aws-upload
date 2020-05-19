package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// JSON func
func JSON(w http.ResponseWriter, status int, res *Response) {
	w.WriteHeader(status)
	data, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Marshal Error: %v", err.Error())
	}

	w.Write([]byte(data))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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

func imgixHandler(w http.ResponseWriter, r *http.Request) {
	var hvalue string
	var wvalue string
	var pathCustom string

	params := mux.Vars(r)
	myImage := params["url"]

	if myImage == "" {
		w.WriteHeader(400)
		return
	}

	// Get original
	path := imgIXClient.Path("/" + myImage)

	paths := map[string]interface{}{
		"original": path,
	}

	// Get custom img
	queries := r.URL.Query()
	hvalue = queries.Get("h")
	wvalue = queries.Get("w")

	if hvalue != "" && wvalue != "" {
		pathCustom = imgIXClient.PathWithParams("/"+myImage, url.Values{
			"h": []string{hvalue},
			"w": []string{wvalue},
		})

		paths["custom"] = pathCustom
	}

	res := &Response{
		Data: paths,
	}

	JSON(w, 200, res)
}
