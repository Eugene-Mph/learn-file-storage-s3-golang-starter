package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	const maxMemory = 10 << 20 // 10MB

	err = r.ParseMultipartForm(maxMemory)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse request", err)
		return
	}

	file, header, err := r.FormFile("thumbnail")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse from file", err)
		return
	}

	defer file.Close()

	contentType := header.Header.Get("Content-Type")

	if contentType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for file", err)
		return
	}

	// log.Printf("content type: %s | file extension: %s", contentType, fileExtension)
	assetPath := getAssetPath(videoID, contentType)
	assetDiskPath := cfg.getAssetDiskPath(assetPath)

	NewFile, err := os.Create(assetDiskPath)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create img file on server", err)
		return
	}

	_, err = io.Copy(NewFile, file)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Eror saving file data", err)
		return
	}

	dbVedio, err := cfg.db.GetVideo(videoID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find vedio", err)
		return
	}

	if dbVedio.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to update this vedio", err)
		return
	}
	// log.Printf("pre update: %v\n", dbVedio)
	assetPathURL := cfg.getAssetUrl(assetPath)

	dbVedio.ThumbnailURL = &assetPathURL
	log.Printf("file url: %s", assetPathURL)

	err = cfg.db.UpdateVideo(dbVedio)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update vedio data", err)
		return
	}
	// log.Printf("Post update: %v", dbVedio)

	respondWithJSON(w, http.StatusOK, dbVedio)
}
