package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
)

var (
	tagsToSkip   stringSliceFlag
	repositories stringSliceFlag
)

// stringSliceFlag is repeatable option from command line
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprint(*s)
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Images in the container registry
type Images []struct {
	Name string `json:"name"`
}

// Tags that image has
type Tags []Tag

// Tag is one of the Tags
type Tag struct {
	Digest    string   `json:"digest"`
	Tags      []string `json:"tags"`
	Timestamp struct {
		Datetime    string `json:"datetime"`
		Day         int    `json:"day"`
		Hour        int    `json:"hour"`
		Microsecond int    `json:"microsecond"`
		Minute      int    `json:"minute"`
		Month       int    `json:"month"`
		Second      int    `json:"second"`
		Year        int    `json:"year"`
	} `json:"timestamp"`
}

func gcloud(args ...string) ([]byte, error) {
	return exec.Command("gcloud", args...).CombinedOutput()
}

func getImages(repository string) []string {
	output, err := gcloud("container", "images", "list", "--format", "json", "--repository", repository)
	if err != nil {
		log.Fatalf("could not list images: %s - %v", output, err)
	}
	images := Images{}
	if err := json.Unmarshal(output, &images); err != nil {
		log.Fatalf("could not decode output: %v", err)
	}
	result := []string{}
	for _, image := range images {
		result = append(result, image.Name)
	}
	return result
}

func getImageTags(image string) Tags {
	output, err := gcloud("container", "images", "list-tags", image, "--format", "json")
	if err != nil {
		log.Fatalf("could not list tags for %s: %v", image, err)
	}
	tags := Tags{}
	if err := json.Unmarshal(output, &tags); err != nil {
		log.Fatalf("could not decode output: %v", err)
	}
	return tags
}

func removeImage(image string, tag Tag) {
	imageWithDigest := fmt.Sprintf("%s@%s", image, tag.Digest)
	fmt.Printf(" - deleting %s... created at %s...\n", imageWithDigest[:len(imageWithDigest)-32], tag.Timestamp.Datetime)
	exec.Command("gcloud", "container", "images", "delete", "-q", "--force-delete-tags", imageWithDigest).Run()
}

func okToRemove(tags, tagsToKeep []string) bool {
	for _, remoteTag := range tags {
		for _, tag := range tagsToKeep {
			if tag == remoteTag {
				return false
			}
		}
	}
	return true
}

func main() {
	flag.Var(&repositories, "r", "repositories to go over")
	flag.Var(&tagsToSkip, "t", "tags to skip (by default not 'latest')")
	flag.Parse()
	if len(repositories) == 0 {
		log.Fatalf("please specify atleast one repository")
	}
	if len(tagsToSkip) == 0 {
		tagsToSkip = append(tagsToSkip, "latest")
	}
	for _, repository := range repositories {
		for _, image := range getImages(repository) {
			fmt.Printf("image %s\n", image)
			for _, tag := range getImageTags(image) {
				if okToRemove(tagsToSkip, tag.Tags) {
					removeImage(image, tag)
				}
			}
		}
	}
}
