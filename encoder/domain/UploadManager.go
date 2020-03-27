package domain

import(
	"cloud.google.com/go/storage"
	"context"
	"fmt"
)

func ProcessUpload(f Video, storagePath string, bucketName string, doneUpload chan bool) {
	concurrency := 5

	in := make(chan int, runtime.NumCPU())
	ret := make(chan error)
	paths := f.GetVideoPaths()
	uploadClient, ctx := getClientUpload()
	
	if f.Status != "error" {
		fmt.Println("Uploading ", f.Uuid, "....")
		for x:= 0; x < concurrency; x++ {
			go uploadWorker(f, storagePath, bucketName, in, ret, paths, uploadClient, ctx)
		}

		go func () {
			for x:=0; x < len(paths);x++ {
				in <- x
			}
			close(in)
		}()

		for err := range ret {
			if err != nil {
				fmt.Println(err.Error())
				doneUpload <- true
				break
			}
		}
	}
}

func uploadWorker(video Video, storagePath string, bucketName string, in chan int, returnChan chan error, paths []string, uploadClient *storage.Client, ctx context.Context) {
	for x:= range in {
		fmt.Println("Object: ", x)
		err := video.UploadObject(paths[x], storagePath, bucketName, uploadClient, ctx)
		returnChan <- err
	}

	returnChan <- fmt.Error("Upload completed")
}

func getClientUpload() (*storage.Client, context.Context) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		fmt.Println("Failed to create client:", err.Error());
	}

	return client, ctx
}