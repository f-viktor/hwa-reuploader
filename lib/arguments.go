package hwapro

import (
	"flag"
	"fmt"
	"os"
)

type inputArgs struct {
	Email  string
	Pass   string
	Advert string
	Save   bool
	Upload string
}

// parse command line arguments
func ParseArgs() inputArgs {
	email := flag.String("u", "example@domain.com", "User e-mail address")
	pass := flag.String("p", "P4ssw0rd!", "User password")
	advert := flag.String("a", "https://hardverapro.hu/...", "Full url of an Archived hwapro advert")
	save := flag.Bool("s", false, "Save the advertisement only (do not repost)")
	upload := flag.String("l", "ads/ad_title", "Post from local folder (set local folder path)")

	flag.Parse()

	args := inputArgs{*email, *pass, *advert, *save, *upload}

	if args.Email == "example@domain.com" || args.Pass == "P4ssw0rd!" {
		flag.Usage()
		os.Exit(2)
	}

	if args.Save && args.Upload != "ads/ad_title" {
		fmt.Println("[!] Do not specify Upload location if you only intend to Save")
		flag.Usage()
		os.Exit(2)
	}

	if args.Save && args.Advert == "https://hardverapro.hu/..." {
		fmt.Println("[!] Specify Advert URL if you intend to Save")
		flag.Usage()
		os.Exit(2)
	}

	if args.Advert != "https://hardverapro.hu/..." && args.Upload != "ads/ad_title" {
		fmt.Println("[!] Do not specify Advert URL if you are uploading from local Save")
		flag.Usage()
		os.Exit(2)
	}

	return args
}
