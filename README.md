# tfstate
Restful Terraform remote state server

Refer to the Documentation [here](https://www.terraform.io/docs/state/remote/http.html) on how to setup your terraform to talk to a resftful server
Currently this only supports saving the State to MongoDB after the Restful server receives it, but you can have more implementations.
Look at mongo.go for a Sample implementation and storeage.go is the basic interface that every engine needs to implement.

Make sure your GOPATH and all is set.

## Getthing Running
  * Do `go get github.com/meson10/tfstate`
  * Execute `tfstate`
  
You can run `tfstate --help` to check usage:

```shell
piyush:~  λ tfstate --help
Usage of tfstate:
  -config string
    	Location of the yaml config file
```

### Sample Configuration file

```yaml
mongo:
  host: hello.mlab.com:15194
  database: terraform
  username: transformer
  password: 0hS0sw33t
```

### Default Configuration

```yaml
mongo:
  host: "127.0.0.1:27017"
  database: terraform

```
