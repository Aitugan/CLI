# CLI
### DAR Final Project. 
#### This is the project I did in order to receive the certificate of DAR internship. 
#### It widely shows abilities I could gain while being an intern, like 
 - creating CLI applications, 
 - testing, 
 - docummenting, 
 - using http methods, 
 - extracting data and etc.

##### To run CLI, type $go run main.go run
You should enter some data to configuration file like:
interface, upstreams with path,method,backends,proxy method.
##### To run CLI daemonly, type "$go run main.go run -daemon" or "$go run main.go run -d"

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

###### In order to change configuration file, run CLI, then use special flag for "run" or "reload" commands: 
- $go run main.go run example.json 
###### OR 
- $go run main.go run example.json -d


##### If you mistakenly entered wrong configuration file or would like to start again, use 
- $go run main reload
##### command


