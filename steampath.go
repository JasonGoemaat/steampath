package steampath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/andygrunwald/vdf"
	"golang.org/x/sys/windows/registry"
)

type SteamApp struct {
	Id         string
	AppId      string // same as 'id', but read from manifest
	Name       string // read from manifest
	InstallDir string // read from manifest
	GamePath   string // path inside locally installed folder (adding library, steamapps/common, and installdir)
	SavePath   string // impossible to know from steam config, but can set if we know it
}

// Get steam install path
// TODO: Non-windows version
func GetSteamPath() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return "", errors.New("cannot find steam registry key")
	}
	defer k.Close()

	s, _, err := k.GetStringValue("InstallPath")
	if err != nil {
		return "", errors.New("cannot find installPath in steam registry key")
	}
	return s, nil
}

// Return map with vdf content of libraryfolders.vdf
// TODO: Maybe check for capitalization 'LibraryFolders.vdf' as well like node module?
func GetLibraryFoldersContent() (map[string]interface{}, error) {
	steamPath, err := GetSteamPath()
	if err != nil {
		return nil, err
	}

	libraryFoldersPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")

	f, err := os.Open(libraryFoldersPath)
	if err != nil {
		return nil, err
	}

	p := vdf.NewParser(f)
	m, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Return paths for steam libraries checking config and verifying existance
func GetLibraryPaths() ([]string, error) {
	vdf, err := GetLibraryFoldersContent()
	if err != nil {
		return nil, err
	}

	// libraryfolders := map[string]interface{}(vdf["libraryfolders"])
	libraryfolders, ok := vdf["libraryfolders"].(map[string]interface{})
	if !ok {
		return nil, errors.New("libraryfolders not string map")
	}

	paths := []string{}
	for _, v := range libraryfolders {
		if v == nil {
			return nil, errors.New("no value from GetLibraryFoldersContent()")
		}
		libraryPath := v.(map[string]interface{})["path"].(string)
		stat, err := os.Stat(libraryPath)
		if err != nil {
			// NOOP - don't error, just don't add missing path
			//fmt.Println("Path Error:", err)
			//fmt.Println("\t", libraryPath)
		} else if stat.IsDir() {
			paths = append(paths, libraryPath)
		} else {
			fmt.Println("Not Directory:")
			fmt.Println("\t", libraryPath)
		}
	}

	return paths, nil
}

// Get steam library installed game details given an id, or nil if not found
func GetApp(id string, args ...string) (*SteamApp, error) {
	manifestPath, librarypath, err := getAppManifestPath(id)
	if err != nil {
		return nil, err
	}

	fmt.Println("GetApp() path:", manifestPath)

	f, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}

	p := vdf.NewParser(f)
	m, err := p.Parse()
	if err != nil {
		return nil, err
	}

	appState := m["AppState"].(map[string]interface{})

	app := SteamApp{Id: id}
	app.AppId = appState["appid"].(string)
	app.Name = appState["name"].(string)
	app.InstallDir = appState["installdir"].(string)
	app.GamePath = filepath.Join(librarypath, "steamapps", "common", app.InstallDir)
	savePath, exists := knownSavePaths[id]
	if exists {
		app.SavePath = savePath
	} else if len(args) > 0 {
		app.SavePath = args[0]
	}
	// fmt.Printf("%+v\n", app)
	return &app, nil
}

// Get path(s) for remote folders given a game id.
// Searches in each library folder for a path like
// 'userdata\<XXX>\<id>\remote' where <XXX> is some
// identifier for a local user that I don't know how
// to get, though I know it isn't a user's steam id.
func GetRemotePaths(id string) ([]string, error) {
	results := []string{}
	steamPath, err := GetSteamPath()
	userdataPath := filepath.Join(steamPath, "userdata")
	entries, err := os.ReadDir(userdataPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		remotePath := filepath.Join(userdataPath, entry.Name(), id, "remote")
		stat, err := os.Stat(remotePath)
		if err == nil && stat.IsDir() {
			results = append(results, remotePath)
		}
	}
	return results, nil
}

func getAppManifestPath(id string) (string, string, error) {
	paths, err := GetLibraryPaths()
	if err != nil {
		return "", "", errors.New("notfound")
	}

	manifestName := fmt.Sprintf("appmanifest_%s.acf", id)

	for _, path := range paths {
		manifestPath := filepath.Join(path, "steamapps", manifestName)
		stat, err := os.Stat(manifestPath)
		if err == nil && !stat.IsDir() {
			return manifestPath, path, nil
		}
		// fmt.Println(id, "not found in", manifestPath)
	}

	return "", "", errors.New("notfound")
}

// Known locations for save files based on id
var knownSavePaths = map[string]string{
	"546430": `%LOCALAPPDATA%\Robotality\Pathway`,
}

var GameIds = struct {
	Pathway string
}{
	Pathway: "546430",
}
