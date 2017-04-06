package coverartarchive

import "encoding/json"

type Thumbnails struct {
	Large string `json:"large"`
	Small string `json:"small"`
}

type Image struct {
	ID         json.Number `json:"id,Number"`
	URL        string      `json:"image"`
	Edit       int         `json:"edit"`
	Approved   bool        `json:"approved"`
	Back       bool        `json:"back"`
	Front      bool        `json:"front"`
	Comment    string      `json:"comment"`
	Thumbnails Thumbnails  `json:"thumbnails"`
	Types      []string    `json:"types"`
}

type CoverArtResponse struct {
	Images  []*Image `json:"images"`
	Release string   `json:"release"`
}

func (c *CoverArtResponse) Front() *Image {
	for _, image := range c.Images {
		if image.Front && !image.Back {
			return image
		}
	}

	return nil
}

func (c *CoverArtResponse) FrontSmallThumbnailURL() string {
	for _, image := range c.Images {
		if image.Front && !image.Back {
			return image.Thumbnails.Small
		}
	}

	return ""
}

func (c *CoverArtResponse) FrontLargeThumbnailURL() string {
	for _, image := range c.Images {
		if image.Front && !image.Back {
			return image.Thumbnails.Large
		}
	}

	return ""
}

func (c *CoverArtResponse) Back() *Image {
	for _, image := range c.Images {
		if image.Back && !image.Front {
			return image
		}
	}

	return nil
}

func (c *CoverArtResponse) BackSmallThumbnailURL() string {
	for _, image := range c.Images {
		if image.Back && !image.Front {
			return image.Thumbnails.Small
		}
	}

	return ""
}

func (c *CoverArtResponse) BackLargeThumbnailURL() string {
	for _, image := range c.Images {
		if image.Back && !image.Front {
			return image.Thumbnails.Large
		}
	}

	return ""
}
