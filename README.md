simple gist client for console.

# Install
``` sh
	$ go get github.com/macrat/gist
```

You must set username and personal access token of github into environment variables.
Please write like this into .bashrc if you use bash.

``` sh
	export GIT_USERNAME=your-ID-of-github
	export GIT_TOKEN=your-personal-access-token-of-github
```

# Usage
## View your gists
``` sh
	$ gist
```

If you want suppress number of items, please use `-n` option.

``` sh
	$ gist -n 5
```

Display 5 items in this example.

You can show more information by `-v` option.

## Show gist
You can show gist by ID or index.

ID is permanent identifier that set by github.
This way is faster than way that use index number.

Index is relative number from latest gist that starts from 0.
Be careful, this number is not permanent because index will change if you added new gist.

You can use by same way both of ID and index.
In this document, use index. Please replace with ID if you want use ID.

``` sh
	$ gist 0
```

This example will show contents of latest gist.
If you want see, please use `-v` option.

Files will splitted like head command if gist has multiple files.
You can change number of items with `-n` option.

``` sh
	$ gist -n 2 0
```

This example will show only first and second file of gist.

You can specifiy filename like this.

``` sh
	$ gist 0/test.txt
```

This example will show contents of `test.txt` in the gist.

## Create new gist
Input into stdin, and give filename to `-c` option.
Like this.

``` sh
	$ echo 'hello world' | gist -c 'filename.txt'
```

This example make new gist that named `filename.txt`.
And contents of new gist is `hello world`.

If you want write description, please use `-d` option.

``` sh
	$ echo 'this is test' | gist -c 'newgist.txt' -d 'this is description'
```

## Update gist
Replace gist file with stdin.

``` sh
	$ echo 'new content' | gist -u 0
```

If some files exist, you must specify file name.

``` sh
	$ echo 'new content' | gist -u 0/name.txt
```

You can update description with `-d` option..

``` sh
	$ echo 'new content' | gist -u -d 'new description' 0
```

There is no way to update description without update content.
Please replace with same content if you want.

## Delete gist
There is no way yet.

# License
[MIT License](https://opensource.org/licenses/MIT) (c)2016 [MacRat](http://blanktar.jp)
