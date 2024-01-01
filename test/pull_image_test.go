package pkg_test

import (
	"os"
	"testing"

	iza_pkg "github.com/joey00072/iza/pkg"
	types "github.com/joey00072/iza/pkg/types"
)

// TestPullImage tests the PullImage function in various scenarios.
func TestPullImage(t *testing.T) {
	// Define a struct for test cases
	testCases := []struct {
		name        string
		opts        types.ImageOptions
		expectError bool
	}{
		{
			name: "Valid Image Without Cache",
			opts: types.ImageOptions{
				ImageName:   "hello-world", // Use an image that is universally available and small in size
				TarballPath: "/tmp/hello-world.tar",
			},
			expectError: false,
		},
		{
			name: "Invalid Image Name",
			opts: types.ImageOptions{
				ImageName:   "nonexistent-image",
				TarballPath: "/tmp/nonexistent.tar",
			},
			expectError: true,
		},
		// Add more test cases as needed
	}

	// Iterate over each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := iza_pkg.PullImage(tc.opts)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for '%s', but none occurred", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for '%s', but got: %v", tc.name, err)
				}
				// Check if the tarball file exists
				if _, err := os.Stat(tc.opts.TarballPath); os.IsNotExist(err) {
					t.Errorf("Expected tarball file to be created for '%s', but it was not", tc.name)
				}
				// Cleanup
				os.Remove(tc.opts.TarballPath)
			}
		})
	}
}
