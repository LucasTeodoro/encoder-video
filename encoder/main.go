package main

import log "github.com/sirupsen/logrus"
import v "encoder/domain"

func main() {
	var video v.Video
	doneUpload := make(chan bool)
	data := []byte("{\"uuid\": \"batata123\", \"path\": \"batata.mp4\", \"status\": \"1\"}")
	video.Unmarshal(data)
	video.Download("batata.mp4", "/tmp")
	video.Fragment("/tmp")
	video.Encode("/tmp")

	go v.ProcessUpload(video, "/tmp", "codeeducationtest", doneUpload)
	<-doneUpload
	video.Finish("/tmp")
	log.Info(video.Path)
}