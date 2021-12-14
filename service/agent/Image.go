package agent

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
	"log"
	"os"
)

func ImagePull(imageName string) error {
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, out)
	return err
}

func ImageList() ([]types.ImageSummary, error) {
	return cli.ImageList(ctx, types.ImageListOptions{})
}

func ImagePrune() error {
	resp, err := cli.ImagesPrune(ctx, filters.NewArgs())
	if err != nil {
		return err
	}
	log.Println("ImagePrune.imageName length ", resp)

	return err
}

func ImageRemove(imageId string) error {
	items, err := cli.ImageRemove(ctx, imageId, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	log.Println("ImageRemove.imageId length ", len(items))
	for _, item := range items {
		log.Println("ImageRemove.imageId:", imageId, item)
	}

	return err
}
