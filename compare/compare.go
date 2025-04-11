package compare

import (
	"fmt"
	"github.com/tobischo/gokeepasslib/v3"
	"strings"
)

type EntryData struct {
	Path     string
	Title    string
	Username string
	Password string
	URL      string
	Notes    string
}

func flattenEntries(group gokeepasslib.Group, path string, entries *[]EntryData) {
	currentPath := strings.Trim(path+"/"+group.Name, "/")

	for _, entry := range group.Entries {
		*entries = append(*entries, EntryData{
			Path:     currentPath,
			Title:    entry.GetTitle(),
			Username: entry.GetContent("UserName"),
			Password: entry.GetPassword(),
			URL:      entry.GetContent("URL"),
			Notes:    entry.GetContent("Notes"),
		})
	}

	for _, subFolder := range group.Groups {
		flattenEntries(subFolder, currentPath, entries)
	}
}

func buildMap(entries []EntryData) map[string]EntryData {
	m := make(map[string]EntryData)
	for _, e := range entries {
		key := e.Path + "/" + e.Title
		m[key] = e
	}
	return m
}

func CompareDatabases(db1, db2 *gokeepasslib.Database) string {
	var entries1, entries2 []EntryData

	flattenEntries(*&db1.Content.Root.Groups[0], "", &entries1)
	flattenEntries(*&db2.Content.Root.Groups[0], "", &entries2)

	entryMap1 := buildMap(entries1)
	entryMap2 := buildMap(entries2)

	result := "🔍 Comparing two KeePass databases:\n"

	for key, e1 := range entryMap1 {
		if e2, exists := entryMap2[key]; exists {
			var diffs []string

			if e1.Username != e2.Username {
				diffs = append(diffs, fmt.Sprintf("Username: '%s' → '%s'", e1.Username, e2.Username))
			}
			if e1.Password != e2.Password {
				diffs = append(diffs, "Password: Changed")
			}
			if e1.URL != e2.URL {
				diffs = append(diffs, fmt.Sprintf("URL: '%s' → '%s'", e1.URL, e2.URL))
			}
			if e1.Notes != e2.Notes {
				diffs = append(diffs, "Notes: Changed")
			}
			if len(diffs) > 0 {
				result += fmt.Sprintf("✏️  Changed: %s\n", key)
				for _, d := range diffs {
					result += fmt.Sprintf("    - %s\n", d)
				}
			}
		} else {
			result += fmt.Sprintf("➖ Missing in DB2: %s\n", key)
		}
	}

	// Добавленные в DB2
	for key := range entryMap2 {
		if _, ok := entryMap1[key]; !ok {
			result += fmt.Sprintf("➕ Added to DB2: %s\n", key)
		}
	}

	return result
}
