package cmd

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
)

func InitCmd() (*string, *string, *string, *bool) {
	app := kingpin.New("sifu-box", "A quick and simple application for transforming config file into sing-box format")
	app.Version("1.0.0")
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	run := app.Command("run", "Boot up the application.")
	work := run.Flag("workdir", "Path of the working directory").Short('w').PlaceHolder("/opt/sifubox/lib").Required().String()
	config := run.Flag("config", "Path of the config directory").Short('c').PlaceHolder("/opt/sifubox/config").Required().String()
	listen := run.Flag("listen", "Address to listen on").Short('l').PlaceHolder(":8080").Default(":8080").String()
	server := run.Flag("server", "Mod of the application run").Short('s').Bool()
	run.Help(`Boot up the application. The flags "--config" and "--workdir" are required. If the flag "--listen" is empty, it will take ":8080" by default.`)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *server && *listen == "" {
		kingpin.Fatalf("When --server is true, --listen must be specified.")
	}
	return work, config, listen, server
}