# CLI
DAR Final Project

# To run it, type $go run main.go run
You should enter some data to configuration file like:
interface, upstreams with path,method,bacends,proxy method.

This is how default configuration file looks like:
{
    "interface": ":8080",
    "upstreams": [
      {
        "path": "example1",
        "method": "GET",
        "backends": [
          "http://server1.com:9090/asd",
          "http://server2.com:9090/asd"
        ],
        "proxyMethod": "round-robin"
      },
      {
        "path": "example2",
        "method": "GET",
        "backends": [
          "http://server1.com:9090/asd",
          "http://server2.com:9090/asd"
        ],
        "proxyMethod": "anycast"
      }
    ]
  }

If you mistakenly entered wrong configuration file or would like to start again, use 
# $go run main reload
command


