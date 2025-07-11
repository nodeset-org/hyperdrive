package service

import (
	"fmt"
	"strings"

	"github.com/distribution/reference"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/client"
	"github.com/nodeset-org/hyperdrive/hyperdrive-cli/utils"
	"github.com/urfave/cli/v2"
)

// Settings
const (
	clientDataVolumeName string = "/ethclient"
	dataFolderVolumeName string = "/.hyperdrive/data"

	PruneFreeSpaceRequired uint64 = 50 * 1024 * 1024 * 1024
)

// Get the compose file paths for a CLI context
func getComposeFiles(c *cli.Context) []string {
	return c.StringSlice(utils.ComposeFileFlag.Name)
}

// Handle a network change by terminating the service, deleting everything, and starting over
func changeNetworks(c *cli.Context) error {
	// Create a new Hyperdrive client - important to ensure the config is loaded from disk and isn't the stale old one
	hd, err := client.NewHyperdriveClientFromCtx(c)
	if err != nil {
		return err
	}
	composeFiles := getComposeFiles(c)

	// Purge the data folder
	fmt.Print("Purging data folder... ")
	err = hd.PurgeData(composeFiles, false)
	if err != nil {
		return fmt.Errorf("error purging data folder: %w", err)
	}

	// Terminate the current setup
	fmt.Print("Removing old installation... ")
	err = hd.DownService(composeFiles, true)
	if err != nil {
		return fmt.Errorf("error terminating old installation: %w", err)
	}
	fmt.Println("done")

	// Start the service
	fmt.Print("Starting Hyperdrive... ")
	err = hd.StartService(getComposeFiles(c))
	if err != nil {
		return fmt.Errorf("error starting service: %w", err)
	}
	fmt.Println("done")

	return nil
}

// A parsed breakdown of a Docker image string's components
type DockerImageInfo struct {
	// The domain for the repository that the image is pulled from, e.g. "docker.io" or "ghcr.io"
	Domain string

	// The vendor (owning organization the image is published under). Since Docker doesn't really have
	// the concept of a vendor, this is everything in the repository name before the first '/'.
	Vendor string

	// The image name under the vendor. This is everything after the first '/' in the
	// repository name.
	Image string

	// The tag of the image, AKA the "version" of it pulled from the repository for this image.
	Tag string
}

// Extract the image origin details from a Docker image string
func getDockerImageInfo(fullImageName string) (DockerImageInfo, error) {
	// Return the empty string if the image didn't exist (probably because this is the first time starting it up)
	if fullImageName == "" {
		return DockerImageInfo{}, nil
	}

	// Parse the image string
	namedRef, err := reference.ParseNormalizedNamed(fullImageName)
	if err != nil {
		return DockerImageInfo{}, fmt.Errorf("error parsing Docker image string [%s]: %w", fullImageName, err)
	}

	// Get the vendor and image name
	repo := reference.Path(namedRef)
	vendor, image, foundSlash := strings.Cut(repo, "/")
	if !foundSlash {
		return DockerImageInfo{}, fmt.Errorf("Docker image string [%s] does not contain a vendor/image format", fullImageName)
	}

	// Get the tag of the image
	tag := ""
	if tagged, ok := namedRef.(reference.Tagged); ok {
		tag = tagged.Tag()
	}

	// Create the DockerImageInfo struct
	imageInfo := DockerImageInfo{
		Domain: reference.Domain(namedRef),
		Vendor: vendor,
		Image:  image,
		Tag:    tag,
	}
	return imageInfo, nil
}
