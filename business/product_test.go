package business

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_MarshalUmnarshal(t *testing.T) {
	assert := assert.New(t)

	// START_TEST_PRODUCT OMIT
	link, _ := url.Parse("https://localhost:8080/link")
	imageLink, _ := url.Parse("https://localhost:8080/image")
	p := &Product{
		SKU:                 "id",
		Title:               "title",
		Description:         "description",
		Link:                link,
		ImageLink:           imageLink,
		AdditionalImageLink: imageLink,
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	// END_TEST_PRODUCT OMIT
	var pu *Product
	err = json.Unmarshal(b, &pu)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(p, pu, "The two structure should be the same.")
}
