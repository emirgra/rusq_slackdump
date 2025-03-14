package source

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
	// STdump is the storage type for the dump format.
	STdump
	// STAvatar is the storage type for the avatar storage.
	STAvatar
)

// Set translates the string value into the ExportType, satisfies flag.Value
// interface.  It is based on the declarations generated by stringer.
//
// It is imperative that the stringer is generated prior to calling this
// function.
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
	STdump:       DumpFilepath,
	STnone:       func(*slack.Channel, *slack.File) string { return "" },
}
