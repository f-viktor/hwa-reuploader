package main

import (
	"hwapro/lib"
)

func main() {
	// parse command-line arguments
	args := hwapro.ParseArgs()
	sess := hwapro.UserSession{}
	sess.Login(args.Email, args.Pass)

	if args.Save && args.Advert != "https://hardverapro.hu/..." {
		ad := hwapro.ParseAdvertisment(args.Advert)
		ad.SaveAdvertisment()
	} else if args.Upload != "ads/ad_title" {
		ad := hwapro.LoadAdvertisment(args.Upload)
		ad.RepostSaved(&sess)
	} else if args.Advert != "https://hardverapro.hu/..." {
		ad := hwapro.ParseAdvertisment(args.Advert)
		ad.SaveAdvertisment()
		ad.RepostSaved(&sess)
	}
}
