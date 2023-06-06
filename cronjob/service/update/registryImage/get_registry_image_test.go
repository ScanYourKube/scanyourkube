package registry_image

import (
	"fmt"
	dto_service_update "github.com/scanyourkube/cronjob/dto/service/update"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var registryResp = `{
	"schemaVersion": 2,
	"mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
	"manifests": [
	  {
		"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
		"digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
		"size": 7143,
		"platform": {
		  "architecture": "ppc64le",
		  "os": "linux"
		}
	  },
	  {
		"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
		"digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
		"size": 7682,
		"platform": {
		  "architecture": "amd64",
		  "os": "linux",
		  "features": [
			"sse4"
		  ]
		}
	  }
	]
  }`

type validImageTest struct {
	arg1            dto_service_update.UpdateServiceImageDto
	expectedName    string
	expectedVersion string
	versionExpected bool
	errorExpected   bool
}

type shouldUpdateImageTest struct {
	arg1          dto_service_update.UpdateServiceImageDto
	arg2          string
	shouldUpdate  bool
	errorExpected bool
}

var validImageTests = []validImageTest{
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "trittale/thesis_nginx:1.0.0",
			ResourceHash: "980475d36f155d10a2f3abba2c074192cd6d14ec4361e780abcc787496909ae6",
		},
		expectedName:    "trittale/thesis_nginx",
		expectedVersion: "1.0.1",
		versionExpected: true,
		errorExpected:   false,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "postgres:latest",
			ResourceHash: "6b91d38a9c596fa4e6a1276f6f81810882d9f292a09f9cf2647c6a554c8b6d00",
		},
		expectedName:    "library/postgres",
		expectedVersion: "latest",
		versionExpected: true,
		errorExpected:   false,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "postgres:bullseye",
			ResourceHash: "6b91d38a9c596fa4e6a1276f6f81810882d9f292a09f9cf2647c6a554c8b6d00",
		},
		expectedName:    "library/postgres",
		expectedVersion: "bullseye",
		versionExpected: true,
		errorExpected:   false,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "postgres:15",
			ResourceHash: "6b91d38a9c596fa4e6a1276f6f81810882d9f292a09f9cf2647c6a554c8b6d00",
		},
		expectedName:    "library/postgres",
		expectedVersion: "15",
		versionExpected: true,
		errorExpected:   false,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "trittale/thesis_nginx123:1.0.0",
			ResourceHash: "980475d36f155d10a2f3abba2c074192cd6d14ec4361e780abcc787496909ae6",
		},
		expectedName:    "trittale/thesis_nginx123",
		expectedVersion: "",
		versionExpected: false,
		errorExpected:   true,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "alpine:3.16",
			ResourceHash: "e28792ec7904bff56f22df296d78cc1188caf347bd824570d0ecf235e4f6e4cd",
		},
		expectedName:    "library/alpine",
		expectedVersion: "3.18",
		versionExpected: true,
		errorExpected:   false,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "postgres:14.7-alpine",
			ResourceHash: "c06405f9394f2be49ba284304758c50a26770c1fe3d4bcce5ed877e617a3b137",
		},
		expectedName:    "library/postgres",
		expectedVersion: "14.8-alpine",
		versionExpected: true,
		errorExpected:   false,
	},
}

var shouldUpdateImageTests = []shouldUpdateImageTest{
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "trittale/thesis_nginx:1.0.0",
			ResourceHash: "980475d36f155d10a2f3abba2c074192cd6d14ec4361e780abcc787496909ae6",
		},
		arg2:          "1.0.0",
		shouldUpdate:  false,
		errorExpected: false,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "trittale/thesis_nginx123:1.0.0",
			ResourceHash: "980475d36f155d10a2f3abba2c074192cd6d14ec4361e780abcc787496909ae6",
		},
		arg2:          "1.0.0",
		shouldUpdate:  false,
		errorExpected: true,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "trittale/thesis_nginx:1.0.0",
			ResourceHash: "dummyValue",
		},
		arg2:          "1.0.0",
		shouldUpdate:  true,
		errorExpected: false,
	},
	{
		arg1: dto_service_update.UpdateServiceImageDto{
			ResourceName: "redis:5",
			ResourceHash: "fc5ecd863862f89f04334b7cbb57e93c9790478ea8188a49f6e57b0967d38c75",
		},
		arg2:          "5",
		shouldUpdate:  false,
		errorExpected: false,
	},
}

func TestValidImages(t *testing.T) {
	for _, test := range validImageTests {
		image, err := GetRegistryImageToUpdateTo(test.arg1)
		if err != nil && !test.errorExpected {
			t.Error(err)
			continue
		}

		if image == nil && test.versionExpected {
			t.Error("expected version")
			continue
		}

		if image != nil && test.versionExpected {
			t.Logf("%s:%s", image.ImageName, image.Tag)
			if image.ImageName != test.expectedName {
				t.Error("Name mismatch")
			}

			if image.Tag != test.expectedVersion {
				t.Errorf("name: %s version: %sVersion mismatch", image.ImageName, image.Tag)
			}

			continue
		}

	}
}

func TestShouldUpdateImage(t *testing.T) {
	for _, test := range shouldUpdateImageTests {
		shouldUpdate, err := CompareCurrentDigestWithNewVersion(test.arg1, test.arg2)
		if err != nil && !test.errorExpected {
			t.Error(err)
			continue
		}

		if shouldUpdate == nil && test.shouldUpdate {
			t.Error("expected should update")
			continue
		}

		if shouldUpdate != nil && *shouldUpdate == test.shouldUpdate {
			t.Logf("got %t", *shouldUpdate)
			continue

		}

		if shouldUpdate != nil && *shouldUpdate != test.shouldUpdate {
			t.Errorf("got %t expected %t", *shouldUpdate, test.shouldUpdate)
			continue

		}
	}
}

func TestDigestDifferenceWithMockedRequest(t *testing.T) {
	cases := []struct {
		input dto_service_update.UpdateServiceImageDto
		want  bool
		err   bool
	}{
		{
			input: dto_service_update.UpdateServiceImageDto{
				ResourceName: "trittale/thesis_nginx:1.0.0",
				ResourceHash: "980475d36f155d10a2f3abba2c074192cd6d14ec4361e780abcc787496909ae6",
			},
			want: true,
			err:  false,
		},
		{
			input: dto_service_update.UpdateServiceImageDto{
				ResourceName: "trittale/thesis_nginx:1.0.0",
				ResourceHash: "5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
			},
			want: false,
			err:  false,
		},
		{
			input: dto_service_update.UpdateServiceImageDto{
				ResourceName: "trittale/thesis_nginx:1.0.0",
				ResourceHash: "",
			},
			want: true,
			err:  false,
		},
		{
			input: dto_service_update.UpdateServiceImageDto{
				ResourceName: "",
				ResourceHash: "a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
			},
			want: false,
			err:  true,
		},
	}
	for _, c := range cases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Docker-Content-Digest", "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270")
			fmt.Fprintln(w, registryResp)
		}))
		defer ts.Close()

		imageDto := c.input
		url := strings.ReplaceAll(ts.URL, "https://", "")
		url = strings.ReplaceAll(url, "http://", "")
		imageDto.ResourceName = url + "/" + imageDto.ResourceName

		isDifferent, err := CompareCurrentDigestWithNewVersion(imageDto, "v0.2.3")

		if err != nil && !c.err {
			t.Errorf("error while getting digest: %s", err)
		}

		if !c.err && *isDifferent != c.want {
			t.Errorf("unexpected is different: %t, expected: %t, digest was: %s", *isDifferent, c.want, c.input.ResourceHash)
		}

	}

}
