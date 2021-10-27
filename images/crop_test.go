package images_test

import (
	"image"
	"testing"

	"github.com/jphacks/B_2121_server/images"
)

func TestGetCroppedRectangle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input image.Rectangle
		want  image.Rectangle
	}{
		{
			name:  "Square",
			input: image.Rect(0, 0, 100, 100),
			want:  image.Rect(0, 0, 100, 100),
		},
		{
			name:  "Portrait1",
			input: image.Rect(0, 0, 100, 200),
			want:  image.Rect(0, 50, 100, 150),
		},
		{
			name:  "Portrait2",
			input: image.Rect(0, 0, 100, 199),
			want:  image.Rect(0, 49, 100, 149),
		},
		{
			name:  "Landscape1",
			input: image.Rect(0, 0, 200, 100),
			want:  image.Rect(50, 0, 150, 100),
		},
		{
			name:  "Landscape2",
			input: image.Rect(0, 0, 199, 100),
			want:  image.Rect(49, 0, 149, 100),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			img := &images.Image{
				Image:  image.NewRGBA(test.input),
				Format: "png",
			}

			gotImage, err := img.CropToSquare()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			got := gotImage.Bounds()
			if test.want != got {
				t.Fatalf("getCroppedRectangle got=%v, want=%v", got, test.want)
			}
		})
	}
}
