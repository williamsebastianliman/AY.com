package config

import "os"

type Config struct {
  HTTPPort  string
  AuthAddr  string
  UserAddr  string
  MediaAddr string
}

func Load() *Config {
  get := func(key, def string) string {
    if v := os.Getenv(key); v != "" {
      return v
    }
    return def
  }
  return &Config{
    HTTPPort:  get("HTTP_PORT", "8080"),
    AuthAddr:  get("AUTH_SVC_ADDR", "auth:50051"),
    UserAddr:  get("USER_SVC_ADDR", "user:50051"),
    MediaAddr: get("MEDIA_SVC_ADDR", "media:50051"),
  }
}