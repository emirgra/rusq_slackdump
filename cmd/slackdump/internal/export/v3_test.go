package export

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/rusq/fsadapter"
	"github.com/rusq/slack"

	"github.com/rusq/slackdump/v3"
	"github.com/rusq/slackdump/v3/internal/chunk"
	"github.com/rusq/slackdump/v3/internal/chunk/chunktest"
	"github.com/rusq/slackdump/v3/internal/fixtures"
	"github.com/rusq/slackdump/v3/internal/network"
	"github.com/rusq/slackdump/v3/internal/structures"
)

var (
	baseDir   = filepath.Join("..", "..", "..", "..")
	chunkDir  = filepath.Join(baseDir, "tmp", "2")
	guestDir  = filepath.Join(baseDir, "tmp", "guest")
	largeFile = filepath.Join(chunkDir, "C0BBSGYFN.json.gz")
)

func Test_exportV3(t *testing.T) {
	fixtures.SkipInCI(t)
	fixtures.SkipOnWindows(t)
	// // TODO: this is manual
	// t.Run("large file", func(t *testing.T) {
	// 	srv := chunktest.NewDirServer(chunkDir)
	// 	defer srv.Close()
	// 	cl := slack.New("", slack.OptionAPIURL(srv.URL()))

	// 	ctx := logger.NewContext(context.Background(), lg)
	// 	prov := &chunktest.TestAuth{
	// 		FakeToken:      "xoxp-1234567890-1234567890-1234567890-1234567890",
	// 		WantHTTPClient: http.DefaultClient,
	// 	}
	// 	sess, err := slackdump.New(ctx, prov, slackdump.WithSlackClient(cl), slackdump.WithLimits(slackdump.NoLimits))
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	output := filepath.Join(baseDir, "output.zip")
	// 	fsa, err := fsadapter.New(output)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	defer fsa.Close()

	// 	list := &structures.EntityList{Include: []string{"C0BBSGYFN"}}
	// 	if err := exportV3(ctx, sess, fsa, list, export.Config{List: list}); err != nil {
	// 		t.Fatal(err)
	// 	}
	// })
	t.Run("guest user", func(t *testing.T) {
		cd, err := chunk.OpenDir(guestDir)
		if err != nil {
			t.Fatal(err)
		}
		defer cd.Close()
		srv := chunktest.NewDirServer(cd)
		defer srv.Close()
		cl := slack.New("", slack.OptionAPIURL(srv.URL()))

		prov := &chunktest.TestAuth{
			FakeToken:      "xoxp-1234567890-1234567890-1234567890-1234567890",
			WantHTTPClient: http.DefaultClient,
		}
		ctx := context.Background()
		sess, err := slackdump.New(ctx, prov, slackdump.WithSlackClient(cl), slackdump.WithLimits(network.NoLimits))
		if err != nil {
			t.Fatal(err)
		}
		dir := t.TempDir()
		output := filepath.Join(dir, "output.zip")
		fsa, err := fsadapter.New(output)
		if err != nil {
			t.Fatal(err)
		}
		defer fsa.Close()

		list := &structures.EntityList{}
		if err := export(ctx, sess, fsa, list, exportFlags{}); err != nil {
			t.Fatal(err)
		}
	})
}
