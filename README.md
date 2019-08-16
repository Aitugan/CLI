# CLI
DAR Final Project

# To run it, type $go run main.go run
You should enter some data to configuration file like:
interface, upstreams with path,method,bacends,proxy method.

This is how default configuration file looks like:  

{
  "config": [  

    {    
      "interface": ":8081",  
      "upstreams": [  
        {  
          "path": "/example2",  
          "method": "GET",  
          "backends": [  
            "https://www.apple.com/",  
            "https://www.microsoft.com/"  
          ],  
          "proxyMethod": "anycast"  
        },  
        {  
          "path": "/example1",  
          "method": "GET",  
          "backends": [  
            "https://github.com",  
            "https://netflix.com/"  
          ],  
          "proxyMethod": "round-robin"  
        }
      ]
    },

    {
      "interface": ":8080",
      "upstreams": [
        {
          "path": "/example2",
          "method": "GET",
          "backends": [
            "https://www.yandex.com/",
            "https://www.dodopizza.kz/"
          ],
          "proxyMethod": "anycast"
        },
        {
          "path": "/example1",
          "method": "GET",
          "backends": [
            "https://www.google.com/",
            "https://www.youtube.com"
          ],
          "proxyMethod": "round-robin"
        }
      ]
    },

    {
      "interface": ":8082",
      "upstreams": [
        {
          "path": "/example2",
          "method": "GET",
          "backends": [
            "https://www.duckduckgo.com/",
            "https://en.wikipedia.org/"
          ],
          "proxyMethod": "anycast"
        },
        {
          "path": "/example1",
          "method": "GET",
          "backends": [
            "https://code.visualstudio.com/",
            "https://sublimetext.com/"
          ],
          "proxyMethod": "round-robin"
        }
      ]
    }
  ]
}

If you mistakenly entered wrong configuration file or would like to start again, use 
# $go run main reload
command


