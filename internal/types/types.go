package types

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/getnf/embellish/internal/utils"
)

// Fonts

type NerdFonts struct {
	Version string `json:"tag_name"`
	Fonts   []Font `json:"assets"`
}

type Font struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	BrowserDownloadUrl string `json:"browser_download_url"`
	AvailableVersion   string
	InstalledVersion   string
}

func (fs NerdFonts) GetVersion() string {
	return fs.Version
}

func (fs NerdFonts) GetFonts() []Font {
	isTar := func(f Font) bool { return f.ContentType == "application/x-xz" }
	fonts := utils.Filter(fs.Fonts, isTar)
	sort.Slice(fonts, func(i, j int) bool { return strings.ToLower(fonts[i].Name) < strings.ToLower(fonts[j].Name) })
	return fonts
}

func (fs NerdFonts) GetFont(f string) Font {
	isWantedFont := func(x Font) bool { return x.Name == f }
	font := utils.Filter(fs.Fonts, isWantedFont)
	return font[0]
}

func (fs NerdFonts) GetFontsNames() []string {
	fontNames := utils.Fold(fs.Fonts, func(f Font) string {
		return f.Name
	})
	sort.Slice(fontNames, func(i, j int) bool { return strings.ToLower(fontNames[i]) < strings.ToLower(fontNames[j]) })
	return fontNames
}

func (f *Font) AddInstalledVersion(ver string) {
	f.InstalledVersion = ver
}

func (f *Font) AddAvailableVersion(ver string) {
	f.AvailableVersion = ver
}

// Command line argumetns

type InstallCmd struct {
	Fonts []string `arg:"positional" help:"list of space separated fonts to install"`
}

type UninstallCmd struct {
	Fonts []string `arg:"positional" help:"list of space separated fonts to uninstall"`
}

type ListCmd struct {
	Installed bool `arg:"-i" help:"list only installed fonts"`
}

type UpdateCmd struct {
	Update bool `default:"true"`
}

type Args struct {
	Install    *InstallCmd   `arg:"subcommand:install" help:"install fonts"`
	Uninstall  *UninstallCmd `arg:"subcommand:uninstall" help:"uninstall fonts"`
	List       *ListCmd      `arg:"subcommand:list" help:"list fonts"`
	Update     *UpdateCmd    `arg:"subcommand:update" help:"update installed fonts"`
	KeepTars   bool          `arg:"-k" help:"Keep archives in the download location"`
	ForceCheck bool          `arg:"-f" help:"Force checking for updates"`
}

func (Args) Version() string {
	return "getnf v1.0.0"
}

// paths

type Paths struct {
	Download string
	Install  string
	Db       string
}

func NewPaths() *Paths {
	paths := &Paths{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	switch os := utils.OsType(); os {
	case "linux", "solaris", "openbsd", "freebsd", "netbsd":
		paths.Download = filepath.Join(homeDir, "Downloads", "embellish")
		paths.Install = filepath.Join(homeDir, ".local", "share", "fonts")
		paths.Db = filepath.Join(homeDir, ".local", "share", "embellish")
	case "darwin":
		paths.Download = filepath.Join(homeDir, "Downloads", "embellish")
		paths.Install = filepath.Join(homeDir, "Library", "Fonts")
		paths.Db = filepath.Join(homeDir, "Library", "embellish")
	case "windows":
		paths.Download = filepath.Join(homeDir, "Downloads", "embellish")
		paths.Install = filepath.Join("C:\\Windows", "Fonts")
		paths.Db = filepath.Join(homeDir, "AppData", "embellish")
	default:
		log.Fatalln("unsupported operating system")
	}

	os.MkdirAll(paths.Download, 0755)
	os.MkdirAll(paths.Install, 0755)
	os.MkdirAll(paths.Db, 0755)

	return paths
}

func (p *Paths) GetDownloadPath() string {
	return p.Download
}

func (p *Paths) GetInstallPath() string {
	return p.Install
}

func (p *Paths) GetDbPath() string {
	return p.Db
}

// params

type GuiParams struct {
	Data         NerdFonts
	Database     *sql.DB
	DownloadPath string
	ExtractPath  string
}
