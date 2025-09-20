package cmd

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
)

func Command() (string, string, string) {
	app := kingpin.New("sifu-box", "A quick and simple application for transforming config file into sing-box format")
	app.Version("1.1.0")
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	run := app.Command("run", "Boot up the application.")
	listen := run.Flag("listen", "Listen address").Short('l').PlaceHolder(":8080").Default(":8080").String()
	config_path := run.Flag("config", "Path of the config directory").Short('c').PlaceHolder("/opt/sifubox/config").Required().String()
	work_dir := run.Flag("dir", "Path of the working directory").Short('d').PlaceHolder("/opt/sifubox").Default("/opt/sifubox").String()
	run.Help(`Boot up the application. The flags "--config" and "--dir" is required.`)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	return *config_path, *work_dir, *listen
}
