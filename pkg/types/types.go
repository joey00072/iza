package types

import (
	"time"
)

// PullImageOptions holds the options for pulling an image.
type ImageOptions struct {
	ImageName   string
	TarballPath string
	CachePath   string // Optional, path to cache the image
}

type Manifest struct {
	Config   string
	RepoTags []string
	Layers   []string
}

type ImageMetadata struct {
	LastTagTime time.Time `json:",omitempty"`
}

type RootFS struct {
	Type      string
	Layers    []string `json:",omitempty"`
	BaseLayer string   `json:",omitempty"`
}
type Image struct {
	ID          string `json:"Id"`
	RepoTags    []string
	RepoDigests []string

	Comment string
	Created string

	Author       string
	Config       *Config
	Architecture string

	Os string

	Size int64 // Size is the unpacked size of the image

	RootFS   RootFS
	Metadata ImageMetadata
}

type Container struct {
	ID       string
	Dir      string
	Image    *Image
	ImageDir string
}

type Config struct {
	Hostname    string   `json:",omitempty"` // Hostname
	User        string   `json:",omitempty"` // User that will run the command(s) inside the container, also support user:group
	AttachStdin bool     // Attach the standard input, makes possible user interaction
	Env         []string `json:",omitempty"` // List of environment variable to set in the container
	Cmd         []string `json:",omitempty"` // Command to run when starting the container

	Volumes    map[string]struct{} `json:",omitempty"` // List of volumes (mounts) used for the container
	WorkingDir string              `json:",omitempty"`
	Entrypoint []string            `json:",omitempty"`
	Labels     map[string]string   `json:",omitempty"`
}
