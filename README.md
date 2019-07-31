# rehentai-content-server
A simple ipfs client. Please exec after launching ipfs daemon.

## System Requirements
### Install Go
The build process for requires Go 1.12 or higher.

### Install go-ipfs
[https://github.com/ipfs/go-ipfs#install](https://github.com/ipfs/go-ipfs#install)

### Install Dependency
```shell
$ go get -u github.com/ipfs/go-ipfs-api
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
