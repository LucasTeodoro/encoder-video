package main

import log "github.com/sirupsen/logrus"
import v "encoder/domain"

func main() {
	var video v.Video
	log.Info("Olá ")
	data := []byte("{\"uuid\": \"batata123\", \"path\": \"batata.mp4\", \"status\": \"1\"}")
	video.Unmarshal(data)
	video.Download("batata.mp4", "/tmp")
}