package bytesource

/**
 * Build a ConfigConnector based on file contents
 */

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	api_config "github.com/wunderkraut/radi-api/operation/config"
)

const (
	// How can the config connector identify config files?  if they match this pattern
	FILE_CONFIGCONNECT_FILEMATCHPATTERN = "*.yml"
)

// Constructor for ConfigConnectYmlFiles
func New_ConfigConnectYmlFiles(paths *Paths) *ConfigConnectYmlFiles {
	return &ConfigConnectYmlFiles{
		paths: paths,
	}
}

// A ConfigConnector that looks for files
type ConfigConnectYmlFiles struct {
	paths *Paths
}

func (connect *ConfigConnectYmlFiles) convertKeyToFileName(key string) string {
	return strings.ToLower(key) + ".yml"
}

func (connect *ConfigConnectYmlFiles) findKey(key string) *Files {
	files := Files{}
	filename := connect.convertKeyToFileName(key)

	for _, pathKey := range connect.paths.Order() {
		pathRoot, _ := connect.paths.Get(pathKey)
		fileSource := pathRoot.FullPath(filename)
		files.Add(pathKey, fileSource)
	}

	return &files
}

func (connect *ConfigConnectYmlFiles) Readers(key string) api_config.ScopedReaders {
	readers := api_config.ScopedReaders{}

	files := connect.findKey(key)
	for _, fileKey := range files.Order() {
		file, _ := files.Get(fileKey)
		if reader, err := file.Reader(); err == nil {
			readers.Add(fileKey, reader)
		}
	}

	return readers
}
func (connect *ConfigConnectYmlFiles) Writers(key string) api_config.ScopedWriters {
	writers := api_config.ScopedWriters{}

	files := connect.findKey(key)
	for _, fileKey := range files.Order() {
		file, _ := files.Get(fileKey)
		if writer, err := file.Writer(); err == nil {
			writers.Add(fileKey, writer)
		}
	}

	return writers
}

// List all possible configs, so all possible config files in all
func (connect *ConfigConnectYmlFiles) List() []string {
	files := []string{}
	trackFound := map[string]bool{}

	for _, pathKey := range connect.paths.Order() {
		path, _ := connect.paths.Get(pathKey)

		dirFiles, _ := ioutil.ReadDir(path.PathString())
		for _, f := range dirFiles {
			if !f.IsDir() {
				name := f.Name()
				if matched, _ := filepath.Match(FILE_CONFIGCONNECT_FILEMATCHPATTERN, name); matched {
					if lastPeriod := strings.LastIndex(name, "."); lastPeriod > 0 {
						name = name[:lastPeriod]
					}
					if _, alreadyFound := trackFound[name]; !alreadyFound {
						files = append(files, name)
						trackFound[name] = true
					}
				}
			}
		}
	}

	return files
}
