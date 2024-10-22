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
create folder random name and  in side folder files  text (comments.txt, headers.txt,href.txt,js.text,body.txt)

use now
```
grep -iHnr "word"
```
