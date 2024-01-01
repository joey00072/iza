package types

// PullImageOptions holds the options for pulling an image.
type ImageOptions struct {
	ImageName   string
	TarballPath string
	CachePath   string // Optional, path to cache the image
}
