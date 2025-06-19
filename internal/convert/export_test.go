package convert

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/rusq/fsadapter"
	"github.com/rusq/slack"

	"github.com/rusq/slackdump/v3/internal/chunk"
	"github.com/rusq/slackdump/v3/internal/fixtures"
	"github.com/rusq/slackdump/v3/source"
)

const (
	testSrcDir = "../../tmp/ora600" // TODO: fix manual nature of this/obfuscate
)

var testLogger = slog.Default()

func TestChunkToExport_Validate(t *testing.T) {
	fixtures.SkipInCI(t)
	fixtures.SkipIfNotExist(t, testSrcDir)
	srcDir, err := chunk.OpenDir(testSrcDir)
	if err != nil {
		t.Fatal(err)
	}
	defer srcDir.Close()
	src := source.OpenChunkDir(srcDir, true)
	testTrgDir := t.TempDir()

	type fields struct {
		Src       source.Sourcer
		Trg       fsadapter.FS
		opts      options
		UploadDir string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{}, true},
		{"no source", fields{Trg: fsadapter.NewDirectory(testTrgDir)}, true},
		{"no target", fields{Src: src}, true},
		{
			"valid, no files",
			fields{
				Src: src,
				Trg: fsadapter.NewDirectory(testTrgDir),
				opts: options{
					includeFiles: false,
				},
			},
			false,
		},
		{
			"valid, include files, but no location functions",
			fields{
				Src: src,
				Trg: fsadapter.NewDirectory(testTrgDir),
				opts: options{
					includeFiles: true,
				},
			},
			true,
		},
		{
			"valid, include files, with location functions",
			fields{
				Src: src,
				Trg: fsadapter.NewDirectory(testTrgDir),
				opts: options{
					includeFiles: true,
					trgFileLoc: func(*slack.Channel, *slack.File) string {
						return ""
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ToExport{
				src:  tt.fields.Src,
				trg:  tt.fields.Trg,
				opts: tt.fields.opts,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ChunkToExport.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChunkToExport_Convert(t *testing.T) {
	fixtures.SkipInCI(t)
	fixtures.SkipIfNotExist(t, testSrcDir)
	cd, err := chunk.OpenDir(testSrcDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cd.Close()
	testTrgDir, err := os.MkdirTemp("", "slackdump")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(testTrgDir)
	// var testTrgDir = t.TempDir()
	fsa, err := fsadapter.NewZipFile(filepath.Join(testTrgDir, "slackdump.zip"))
	if err != nil {
		t.Fatal(err)
	}
	defer fsa.Close()
	src := source.OpenChunkDir(cd, true)
	c := NewToExport(src, fsa, WithIncludeFiles(true))

	ctx := t.Context()
	c.opts.lg = testLogger
	if err := c.Convert(ctx); err != nil {
		t.Fatal(err)
	}
}
