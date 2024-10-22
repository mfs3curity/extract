extract headers and source body and href and src->js and comments

#### install 
```
go install github.com/mfs3curity/extract@latest
```

how use
```
subfinder -d target.com -silent | httpx -silent | extract
```

after use 
After executing the command, a folder with a random name will be created containing several files. (comments.txt, headers.txt,href.txt,js.text,body.txt)

use now
```
grep -iHnr "word"
```
