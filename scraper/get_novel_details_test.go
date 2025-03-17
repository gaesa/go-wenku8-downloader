package scraper

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetNovelDetail(t *testing.T) {
	novel, err := GetNovelDetails(1973)
	require.NoError(t, err)
	require.NotEmpty(t, novel)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(novel)
	require.NoError(t, err)

	// copyright animate
	// novel, err := GetNovelDetails(1587)
	// require.NoError(t, err)
	// require.NotEmpty(t, novel)
	// fmt.Printf("%v", novel.LastUpdateTime)

	// // animate
	// novel, err := GetNovelDetails(2975)
	// require.NoError(t, err)
	// require.NotEmpty(t, novel)
	// fmt.Printf("%v", novel.Desc)
}
