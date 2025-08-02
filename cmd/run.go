package cmd

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
)

func Command() (string, string) {
	app := kingpin.New("sifu-box", "A quick and simple application for transforming config file into sing-box format")
	app.Version("1.1.0")
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	run := app.Command("run", "Boot up the application.")
	config := run.Flag("config", "Path of the config directory").Short('c').PlaceHolder("/opt/sifubox/config").Required().String()
	dir := run.Flag("dir", "Path of the working directory").Short('d').PlaceHolder("/opt/sifubox").Default("/opt/sifubox").String()
	run.Help(`Boot up the application. The flags "--config" and "--dir" is required.`)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	return *config, *dir
}
