package config

import (
	"flag"
	"fmt"
	"github.com/goccy/go-json"
	"image-resizing-shared/pkg/utils"
	"log"
	"os"
	"path/filepath"
)

type CompressionSettings struct {
	Lossless bool    `json:"lossless"`
	Quality  float32 `json:"quality"`
}

type Config struct {
	WebServerPort   string              `json:"web_server_port"`
	GrpcServerPort  string              `json:"grpc_server_port"`
	RestUploadToken string              `json:"rest_upload_token"`
	GrpcToken       string              `json:"grpc_token"`
	DBPath          string              `json:"db_path"`
	Compression     CompressionSettings `json:"image_compression"`
}

func LoadConfig(configPath string) *Config {
	flag.Parse()

	cfg := &Config{
		WebServerPort:   utils.GetEnv("WEB_SERVER_PORT", "5689"),
		GrpcServerPort:  utils.GetEnv("GRPC_SERVER_PORT", "50066"),
		RestUploadToken: os.Getenv("REST_UPLOAD_TOKEN"),
		GrpcToken:       os.Getenv("GRPC_TOKEN"),
		DBPath:          utils.GetEnv("DB_PATH", ""),
		Compression: CompressionSettings{
			Lossless: false,
			Quality:  80,
		},
	}

	if configPath == "" {
		return cfg
	}

	execPath, err := os.Executable()
	if err != nil {
		log.Printf("config: failed to resolve binary path: %v", err)
		return cfg
	}
	fullConfigPath := filepath.Join(filepath.Dir(execPath), configPath)

	fmt.Println(fullConfigPath)

	data, err := os.ReadFile(fullConfigPath)
	if err != nil {
		log.Printf("config: failed to read config file %s: %v", fullConfigPath, err)
		return cfg
	}

	var fileCfg Config

	if err := json.Unmarshal(data, &fileCfg); err != nil {
		log.Printf("config: invalid JSON: %v", err)
		return cfg
	}

	if fileCfg.WebServerPort != "" {
		cfg.WebServerPort = fileCfg.WebServerPort
	}
	if fileCfg.GrpcServerPort != "" {
		cfg.GrpcServerPort = fileCfg.GrpcServerPort
	}
	if fileCfg.RestUploadToken != "" {
		cfg.RestUploadToken = fileCfg.RestUploadToken
	}
	if fileCfg.GrpcToken != "" {
		cfg.GrpcToken = fileCfg.GrpcToken
	}
	if fileCfg.DBPath != "" {
		cfg.DBPath = fileCfg.DBPath
	}
	if fileCfg.Compression.Quality != 0 {
		cfg.Compression.Quality = fileCfg.Compression.Quality
	}
	cfg.Compression.Lossless = fileCfg.Compression.Lossless

	return cfg
}
