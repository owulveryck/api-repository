package business

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducts_MarshalUmnarshal(t *testing.T) {
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
	products := &Products{
		Elements: []*Product{p, p},
	}
	b, err := json.Marshal(products)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	// END_TEST_PRODUCT OMIT
	var pu *Products
	err = json.Unmarshal(b, &pu)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(products, pu, "The two structure should be the same.")
}
