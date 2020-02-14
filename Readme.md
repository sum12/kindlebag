kindlebag 
---------

(For jailbroken kindles only)

Small golang based binary to download all articles from wallabag directly to kindle.
No email to @kindle address.

Using the build command (below) create a binary and copy the birnay and config.json to kindle


build:
-----
```
env GOOS=linux GOARCH=arm GOARM=7 go build
```


install:
-------


#### Option1

Jailbroken kindle gives you ssh access to kindle
copy the binary and config.json somewhere on kindle using ssh
execute command
```
./kindlebag -config config.json -outfolder /mnt/us/documents
```


#### Option2

Use this IF you know what are KUAL exteionsions.
use the KUALextensions directory for sample.
- copy its contents (using ssh) to /mnt/us/extensions/kindlebag
- copy the binary (built above) to location pointed in [kindlebag.sh](KUALextension/bin/kindlebag.sh)


## License


wtfpl

I AM NOT RESPONSIBLE TO OUTCOME OF WHAT ANYONE DOES WITH THIS REPO
