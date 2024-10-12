package steampath

import (
	"testing"
)

func TestSteamPath(t *testing.T) {
	result, err := GetSteamPath()
	// fmt.Println("GetSteamPath():", result)
	if err != nil {
		t.Fatal("GetSteampath() error:", err)
	}
	if result == "" {
		t.Fatal("GetSteampath() empty result:", result)
	}
}

func TestGetLibraryFoldersContent(t *testing.T) {
	result, err := GetLibraryFoldersContent()
	if err != nil {
		t.Fatal("GetLibraryFolders error:", err)
	}
	if result == nil {
		t.Fatal("No content from GetLibraryFoldersContent()")
	}
	libfolders := result["libraryfolders"]
	if libfolders == nil {
		t.Fatal("libraryfolders key has no content")
	}
}

func TestGetLibraryPaths(t *testing.T) {
	folders, err := GetLibraryPaths()
	if err != nil {
		t.Fatal("Error calling GetLibraryPaths()")
	}
	if folders == nil {
		t.Fatal("No folders from GetLibraryPaths()")
	}
	for _, folder := range folders {
		if folder == "" {
			t.Fatal("Empty folder")
		}
	}
}

func TestGetApp(t *testing.T) {
	ID := "546430"
	app, err := GetApp(ID)
	if err != nil {
		t.Fatal("Error in GetApp():", err)
	}

	if app == nil {
		t.Fatal("App not found:", ID)
	}
}

func TestGetRemotePaths(t *testing.T) {
	entries, err := GetRemotePaths(GameIds.Pathway)
	if err != nil {
		t.Fatal("Error:", err)
	}
	if entries == nil {
		t.Fatal("No entries")
	}
	// for _, entry := range entries {
	// 	fmt.Println("FOUND:", entry)
	// }
}
