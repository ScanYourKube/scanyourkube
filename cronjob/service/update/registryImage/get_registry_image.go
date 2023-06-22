package registry_image

import (
	"fmt"
	"strings"

	dto_service_update "github.com/scanyourkube/cronjob/dto/service/update"

	"github.com/Masterminds/semver"
	"github.com/docker/distribution/reference"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	log "github.com/sirupsen/logrus"
)

type ImageRegistryStatus struct {
	RegistryImage       RegistryImage
	NewVersionAvailable bool
}

type RegistryImage struct {
	ImageName string
	Tag       string
}

func GetImageRegistryStatus(updateImage dto_service_update.UpdateServiceImageDto) (ImageRegistryStatus, error) {
	registryImage, err := GetRegistryImageToUpdateTo(updateImage)
	if err != nil {
		log.Error(err)
		return ImageRegistryStatus{}, err
	}
	newVersionAvailable, err := CompareCurrentDigestWithNewVersion(updateImage, registryImage.Tag)
	if err != nil {
		log.Error(err)
		return ImageRegistryStatus{}, err
	}

	return ImageRegistryStatus{
		RegistryImage:       *registryImage,
		NewVersionAvailable: *newVersionAvailable,
	}, nil
}

func CompareCurrentDigestWithNewVersion(
	updateImage dto_service_update.UpdateServiceImageDto,
	newVersion string) (*bool, error) {
	named, err := reference.ParseNormalizedNamed(updateImage.ResourceName)
	if err != nil {
		return nil, err
	}
	named = reference.TrimNamed(named)
	named, err = reference.WithTag(named, newVersion)
	if err != nil {
		return nil, err
	}
	ref, err := name.ParseReference(named.String(), name.Insecure)
	if err != nil {
		return nil, err
	}

	newImg, err := remote.Head(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, err
	}

	newDigest := newImg.Digest

	/*
		If the new version is the same as the old version,
		we should not update the image. As the digest might have sha256: as prefix we remove that
	*/
	digestString := strings.Replace(newDigest.String(), "sha256:", "", -1)
	log.Debugf("Local digest: %s, Remote digest: %s", updateImage.ResourceHash, newDigest.String())
	isNotSameImage := digestString != updateImage.ResourceHash

	return &isNotSameImage, nil
}

func GetRegistryImageToUpdateTo(updateImage dto_service_update.UpdateServiceImageDto) (*RegistryImage, error) {
	named, err := reference.ParseNormalizedNamed(updateImage.ResourceName)
	if err != nil {
		return nil, err
	}

	log.Debugf("Got named registry image: %v", named)
	named = reference.TagNameOnly(named)
	tagged, isTagged := named.(reference.NamedTagged)

	if !isTagged {
		return nil, fmt.Errorf("failed to tag registry image with default tag")
	}
	log.Debugf("Got tagged registry image: %v", tagged)
	ref, err := name.ParseReference(updateImage.ResourceName, name.Insecure)
	log.Debugf("Got registry image reference: %v and will connect over %s", ref, ref.Context().Registry.Scheme())
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var imageVersion *semver.Version

	if strings.ToLower(tagged.Tag()) == "latest" {
		return &RegistryImage{
			ImageName: removeRepositoryInformationFromImageName(tagged.Name()),
			Tag:       tagged.Tag(),
		}, nil
	}

	imageVersion, err = semver.NewVersion(tagged.Tag())

	if err != nil {
		return &RegistryImage{
			ImageName: removeRepositoryInformationFromImageName(tagged.Name()),
			Tag:       tagged.Tag(),
		}, nil
	}

	if checkIfTagIsOnlyMajorVersion(tagged.Tag()) {
		return &RegistryImage{
			ImageName: removeRepositoryInformationFromImageName(tagged.Name()),
			Tag:       tagged.Tag(),
		}, nil
	}

	log.Debugf("Getting image tags from registry")
	tags, err := remote.List(ref.Context())
	log.Debugf("Got image tags from registry %v", tags)
	if err != nil {
		return nil, err
	}

	newestRelease := imageVersion
	for _, tag := range tags {
		newestRelease = getNewerTagVersion(tag, newestRelease)
	}

	log.Debugf("Newest release %v", newestRelease)

	return &RegistryImage{
		ImageName: removeRepositoryInformationFromImageName(tagged.Name()),
		Tag:       newestRelease.Original(),
	}, nil
}

func removeRepositoryInformationFromImageName(imageName string) string {
	_, name, found := strings.Cut(imageName, "/")
	if found {
		return name
	}

	return imageName
}

func getNewerTagVersion(tag string, newestVersion *semver.Version) *semver.Version {
	semverTag, err := semver.NewVersion(tag)
	if err != nil {
		return newestVersion
	}

	if checkIfTagIsOnlyMajorAndMinorVersion(newestVersion.Original()) && !checkIfTagIsOnlyMajorAndMinorVersion(tag) {
		return newestVersion
	}

	if semverTag.GreaterThan(newestVersion) &&
		semverTag.Major() == newestVersion.Major() &&
		semverTag.Prerelease() == newestVersion.Prerelease() {
		return semverTag
	}

	return newestVersion
}

func checkIfTagIsOnlyMajorVersion(tag string) bool {
	return len(strings.Split(tag, ".")) == 1
}

func checkIfTagIsOnlyMajorAndMinorVersion(tag string) bool {
	return len(strings.Split(tag, ".")) == 2
}
