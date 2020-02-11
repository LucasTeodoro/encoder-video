package domain

import (
	"context"
	"encoding/json"
	"cloud.google.com/go/storage"
	"ioutil"
)

type Video struct {
	Uuid string `json: "uuid"`
	Path string `json: "path"`
	Status string `json: "status"`
}

func (video *Video) Unmarshal(payload []byte) Video {
	err := json.Unmarshal(payload, &video)
	if err != nil {
		panic(err)
	}

	return *video
}

func (video *Video) Download(bucketName string, storagePath string) (Video, error) {
	ctx := context.Background()

	client, err := storage.NewCliente(ctx)

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())

		return *video, err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(video.Path)

	reader, err := obj.NewReader(ctx)

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())

		return *video, err
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())

		return *video, err
	}

	file, err := os.Create(storagePath + "/" + video.Uuid + ".mp4")

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())

		return *video, err
	}
	
	_, err = file.Write(body)

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())

		return *video, err
	}

	defer file.Close()

	return *video, nil
}