package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
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

func mediaTypeToExt(MediaType string) string {
	parts := strings.Split(MediaType, "/")

	if len(parts) != 2 {
		return ".bin"
	}
	return fmt.Sprintf(".%s", parts[1])
}
