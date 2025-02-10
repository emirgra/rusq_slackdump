package fileproc

import (
	"fmt"
	"strings"

	"github.com/rusq/slack"
)

type StorageType uint8

//go:generate stringer -type=StorageType -trimprefix=ST
const (
	STnone StorageType = iota
	// STstandard is the storage type for the standard file storage.
	STstandard
	// STmattermost is the storage type for Mattermost.
	STmattermost
)

// Set translates the string value into the ExportType, satisfies flag.Value
// interface.  It is based on the declarations generated by stringer.
func (e *StorageType) Set(v string) error {
	v = strings.ToLower(v)
	for i := 0; i < len(_StorageType_index)-1; i++ {
		if strings.ToLower(_StorageType_name[_StorageType_index[i]:_StorageType_index[i+1]]) == v {
			*e = StorageType(i)
			return nil
		}
	}
	return fmt.Errorf("unknown format: %s", v)
}

var StorageTypeFuncs = map[StorageType]func(_ *slack.Channel, f *slack.File) string{
	STmattermost: MattermostFilepath,
	STstandard:   StdFilepath,
	STnone:       func(_ *slack.Channel, f *slack.File) string { return "" },
}
