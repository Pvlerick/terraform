// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudplan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"
)

func TestCloud_loadBasic(t *testing.T) {
	bookmark := SavedPlanBookmark{
		RemotePlanFormat: 1,
		RunID:            "run-GXfuHMkbyHccAGUg",
		Hostname:         "app.terraform.io",
	}

	result, err := LoadSavedPlanBookmark("./testdata/plan-bookmark/bookmark.json")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(bookmark, result, cmp.Comparer(cty.Value.RawEquals)); diff != "" {
		t.Errorf("wrong result\n%s", diff)
	}
}

func TestCloud_saveBasic(t *testing.T) {
	tmp := t.TempDir()
	bookmarkPath := filepath.Join(tmp, "saved-bookmark.json")

	// verify the created path exists
	// os.Stat wants a path, not a file
	_, err := os.Stat(tmp)
	if err != nil {
		t.Fatal("Path does not exist.", err)
	} else {
		b := &SavedPlanBookmark{
			RemotePlanFormat: 1,
			RunID:            "run-GXfuHMkbyHccAGUg",
			Hostname:         "app.terraform.io",
		}
		b.Save(bookmarkPath)
	}
}
