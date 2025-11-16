package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(MediaType string) string {
	base := make([]byte, 32)
	_, err := rand.Read(base)

	if err != nil {
		panic("failed to generate random bytes")
	}

	id := base64.RawStdEncoding.EncodeToString(base)
	ext := mediaTypeToExt(MediaType)
	return fmt.Sprintf("%s%s", id, ext)
}

func (cfg *apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg *apiConfig) getAssetUrl(assetpath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetpath)
}

func (cfg *apiConfig) getObjectURL(Key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, Key)
}

func mediaTypeToExt(MediaType string) string {
	parts := strings.Split(MediaType, "/")

	if len(parts) != 2 {
		return ".bin"
	}
	return fmt.Sprintf(".%s", parts[1])
}

func getVideoAspectRatio(filePath string) (string, error) {

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	buffer := bytes.Buffer{}

	cmd.Stdout = &buffer
	cmd.Run()

	type parameters struct {
		Streams []struct {
			Width  float64 `json:"width"`
			Height float64 `json:"height"`
		} `json:"streams"`
	}

	params := parameters{}

	err := json.Unmarshal(buffer.Bytes(), &params)
	// fmt.Println(params)

	if err != nil {
		return "", err
	}

	videoRatio := params.Streams[0].Width / params.Streams[0].Height
	const portraitRatio = 9.0 / 16.0
	const landscapeRatio = 16.0 / 9.0
	const tolarence = 0.01
	// println(q)

	switch {

	case math.Abs(videoRatio-portraitRatio) < tolarence:
		return "9:16", nil

	case math.Abs(videoRatio-landscapeRatio) < tolarence:
		return "16:9", nil

	default:
		return "other", nil

	}
}

func getKeyPath(aspectRatio, assetPath string) string {
	prefix := ""

	switch aspectRatio {
	case "9:16":
		prefix = "portrait"
	case "16:9":
		prefix = "landscape"
	default:
		prefix = "other"
	}

	// fmt.Println(prefix)

	return fmt.Sprintf("%s/%s", prefix, assetPath)
}
