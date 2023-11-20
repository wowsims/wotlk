package cmd

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wowsims/classic/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
	goproto "google.golang.org/protobuf/proto"
)

var decodeLinkCmd = &cobra.Command{
	Use:   "decodelink [link]",
	Short: "decode wowsims link/url",
	Long:  "decode wowsims link/url",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return decodeLink(args[0])
	},
}

var errInvalidLink = errors.New("invalid wowsims export link")

func decodeLink(link string) error {
	parts := strings.Split(link, "#")
	switch {
	case len(parts) != 2:
		return errInvalidLink
	case parts[1] == "":
		return errInvalidLink
	}

	raw, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("cannot decode proto from link: %w", err)
	}

	r, err := zlib.NewReader(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("cannot create zlib reader: %w", err)
	}
	defer r.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return fmt.Errorf("reading zlib data failed: %w", err)
	}

	var settings goproto.Message
	if strings.Contains(link, "/raid/") {
		settings = &proto.RaidSimSettings{}
	} else {
		settings = &proto.IndividualSimSettings{}
	}

	if err := goproto.Unmarshal(buf.Bytes(), settings); err != nil {
		return fmt.Errorf("cannot unmarshal raw proto: %w", err)
	}

	fmt.Println(protojson.Format(goproto.Message(settings)))
	return nil
}
