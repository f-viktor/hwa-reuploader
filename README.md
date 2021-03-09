# hwa-reuploader

Tool for re-uploading Archived advertisements to hardverapro.hu

## Usage

%hwa-uploader go run main.go -h  
Usage of /tmp/go-build1506704604/b001/exe/main:  
  -a  Full url of an Archived hwapro advert (default "https://hardverapro.hu/...")  
  -l  Post from local folder (set local folder path) (default "ads/ad_title")  
  -p  User password (default "P4ssw0rd!")  
  -s	Save the advertisement only (do not repost)  
  -u  User e-mail address (default "example@domain.com")  

Look, it's not that difficult. You will always need to set `-u` and `-p` to log in.  
There are 3 things you can do after this:  

### Reupload an archived advert:

`go run main.go -u example@domain.com -p yourpassword -a "https://hardverapro.hu/apro/youradverturl.html"`

This will save the given ad, to `./ads` and then immediately repost it.

### Save an archived advert for modification and later Reupload

`go run main.go -u example@domain.com -p yourpassword -s -a "https://hardverapro.hu/apro/youradverturl.html"`

This will only save this advert to `./ads`. You can modify the info text file or the images if you wish. The thumbnail image will always be the last image in the LocalImages list in the list (in the info file).

### Upload a previously saved advert

`go run main.go -u example@domain.com -p yourpassword -l ads/youradvert/`

This will reupload your ad from the path you set. Set the path of the folder that contains the info file.
