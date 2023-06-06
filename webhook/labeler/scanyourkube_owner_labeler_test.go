package labeler

import (
	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	"testing"
)

func TestScanYourKubeOwnerLabelerLabel(t *testing.T) {
	// Create an instance of ScanYourKubeOwnerLabeler
	labeler := ScanYourKubeOwnerLabeler{
		Owner:     "u-user123",
		Operation: admissionv1.Create,
	}

	// Create a sample labels map
	labels := map[string]string{
		"app":  "my-app",
		"team": "my-team",
	}

	// Call the Label method
	labeled := labeler.Label(labels)

	// Verify the modification
	assert.Equal(t, "u-user123", labeled["scanyourkube.io/owner"])
	assert.Equal(t, "my-app", labeled["app"])
	assert.Equal(t, "my-team", labeled["team"])
}

func TestScanYourKubeOwnerLabelerInvalidUserLabel(t *testing.T) {
	// Create an instance of ScanYourKubeOwnerLabeler
	labeler := ScanYourKubeOwnerLabeler{
		Owner:     "not-a-user",
		Operation: admissionv1.Update,
	}

	// Create a sample labels map
	labels := map[string]string{
		"app":  "my-app",
		"team": "my-team",
	}

	// Call the Label method
	labeled := labeler.Label(labels)

	// Verify the modification
	assert.Empty(t, labeled["scanyourkube.io/owner"])
	assert.Equal(t, "my-app", labeled["app"])
	assert.Equal(t, "my-team", labeled["team"])
}
