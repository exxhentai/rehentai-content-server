# rehentai-content-server
A simple ipfs client. Please exec after launching ipfs daemon.

## System Requirements
### Install Go
The build process for requires Go 1.12 or higher.

### Install go-ipfs
[https://github.com/ipfs/go-ipfs#install](https://github.com/ipfs/go-ipfs#install)

### Install Dependency
```shell
$ go get -u -v github.com/ipfs/go-ipfs-api
$ go get -u -v github.com/gorilla/websocket
```

## Usage
Since the frontend server is still WIP, I use `curl` to test file uploading through http request.

```shell=
$ go run main --help
  -port int
         (default 1234)

# upload single file
$ curl -v -F "file[]=@test01.zip" http://localhost:1234/

# upload multiple file
$ curl -v -F "file[]=@test01.zip" \
          -F "file[]=@test02.zip" \
          http://localhost:1234/
```

### File Extension
This server currently only accept file type `zip`. This file will be decompress, and all decompressed files in type `png`, `jpg`, `jpeg`, `bmp` will be uploaded to ipfs node.

### Response
The IPFS hashes will be sent by http response in following format.

```shell
# upload single file
$ curl -v -F "file[]=@test01.zip" http://localhost:1234/
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 1234 (#0)
> POST / HTTP/1.1
> Host: localhost:1234
> User-Agent: curl/7.58.0
> Accept: */*
> Content-Length: 14152942
> Content-Type: multipart/form-data; boundary=------------------------06f9a775df249573
> Expect: 100-continue
>
< HTTP/1.1 100 Continue
< HTTP/1.1 200 OK
< Date: Thu, 01 Aug 2019 10:01:37 GMT
< Content-Length: 129
< Content-Type: text/plain; charset=utf-8
<
QmQg5BL1eeLX6gPhfktLfDxNaGFeR6Z5xmS1d12YoEjHLo 1208408-[Grinp (Neko Toufu)] Onii-chan wa Oshimai!  別當歐尼醬了! [Chinese]
* Connection #0 to host localhost left intact
```
