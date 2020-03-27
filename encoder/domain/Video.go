package domain

import (
	"context"
	"encoding/json"
	"cloud.google.com/go/storage"
	"io/ioutil"
	"os"
	"fmt"
	"os/exec"
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

	client, err := storage.NewClient(ctx)

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

func (video *Video) Fragment(storagePath string) Video {
	err := os.Mkdir(storagePath + "/" + video.Uuid, os.ModePerm)

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
	}

	source := storagePath + "/" + video.Uuid + ".mp4"
	target := storagePath + "/" + video.Uuid + ".frag"
	
	cmd :=  exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
	}

	printOutput(output)

	return *video
}

func (video *Video) Encode(storagePath string) Video {
	cmdArgs := []string{}

	cmdArgs = append(cmdArgs, storagePath + "/" + video.Uuid + ".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, storagePath + "/" + video.Uuid)
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/usr/local/bin")

	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
	}

	printOutput(output)

	return *video
}

func (video *Video) UploadObject(completPath string, storagePath string, bucketName string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(completPath, storagePath + "/")

	f, err := os.Open(completPath)
	if err != nil {
		fmt.Println("Error during the upload", err.Error())
		return err
	}
	defer f.Close()

	wc := client.Bucket(bucketName).Object(path[1]).NewWrite(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role:storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (video *Video) Finish(storagePath string) {
	err := os.Remove(storagePath + "/" + video.Uuid + ".mp4")
	if err != nil {
		fmt.Println("Error removing MP4: ", video.Uuid)
	}
	err := os.Remove(storagePath + "/" + video.Uuid + ".frag")
	if err != nil {
		fmt.Println("Error removing FRAG: ", video.Uuid + ".frag")
	}
	err := os.RemoveAll(storagePath + "/" + video.Uuid)
	if err != nil {
		fmt.Println("Error removing folder: ", video.Path)
	}

	fmt.Println("Files has been removed", video.Uuid)
}

func (video *Video) GetVideoPaths() []string {
	var paths []string
	filepath.Walk("/tmp/convite", func(path string, info os.FileInfo, err error) error {
		paths = append(paths, path)
		return nil
	})

	return paths
}

func printOutput(out []byte) {
	if(len(out) > 0) {
		fmt.Printf("Output: %s\n", string(out))
	}
}