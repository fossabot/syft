package java

import (
	"fmt"

	"github.com/anchore/syft/internal/file"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/pkg/cataloger/generic"
	"github.com/anchore/syft/syft/source"
)

var genericZipGlobs = []string{
	"**/*.zip",
}

// TODO: when the generic archive cataloger is implemented, this should be removed (https://github.com/anchore/syft/issues/246)

// parseZipWrappedJavaArchive is a parser function for java archive contents contained within arbitrary zip files.
func parseZipWrappedJavaArchive(_ source.FileResolver, _ *generic.Environment, reader source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	contentPath, archivePath, cleanupFn, err := saveArchiveToTmp(reader.AccessPath(), reader)
	// note: even on error, we should always run cleanup functions
	defer cleanupFn()
	if err != nil {
		return nil, nil, err
	}

	// we use our zip helper functions instead of that from the archiver package or the standard lib. Why? These helper
	// functions support zips with shell scripts prepended to the file. Specifically, the helpers use the central
	// header at the end of the file to determine where the beginning of the zip payload is (unlike the standard lib
	// or archiver).
	fileManifest, err := file.NewZipFileManifest(archivePath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read files from java archive: %w", err)
	}

	// look for java archives within the zip archive
	return discoverPkgsFromZip(reader.Location, archivePath, contentPath, fileManifest, nil)
}
