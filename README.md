# go-url-shortner
URL Shortener
This is a simple URL shortener built using Go that redirects incoming web requests to their corresponding target URLs. It acts as a basic HTTP server that listens for incoming requests and checks if the requested path matches any predefined mappings. If a match is found, the server will redirect the user to the target URL, similar to a typical URL shortener service.

# How It Works
The URL shortener uses an http.Handler implementation to handle incoming requests. The main logic is built around the concept of mapping paths to their corresponding target URLs. The server maintains a list of path-to-URL mappings and checks incoming requests against this list to determine if a redirection is required.

# Getting Started
To use the URL shortener, follow these steps:

# Clone the repository:

Copy code
git clone https://github.com/mobster1425/go-url-shortener.git
Navigate to the project directory:


# Copy code
cd go-url-shortener
Build and run the application:

go run main.go -yaml=paths.yaml -db=urlshort.db -json=data.json


The URL shortener will start listening on port 8080.

# Adding Redirections
To add new redirections, you need to define the mappings in either YAML or JSON format. The mappings consist of a path and its corresponding target URL. The server will use these mappings to redirect incoming requests to the desired destinations.

# YAML Format
The YAML format should be as follows:

yaml
Copy code
- path: /short
  url: https://www.somesite.com/a-very-long-url

- path: /dogs
  url: https://www.somesite.com/a-story-about-dogs
Save these mappings in a file named paths.yaml in the data directory.

# JSON Format
The JSON format should be as follows:

json
Copy code
[
  {
    "path": "/short",
    "url": "https://www.somesite.com/a-very-long-url"
  },
  {
    "path": "/dogs",
    "url": "https://www.somesite.com/a-story-about-dogs"
  }
]
Save these mappings in a file named data.json in the data directory.

# How the Redirection Works
When the URL shortener receives an incoming request, it will look for a match between the requested path and the mappings defined in the YAML and JSON files. If a match is found, the server will redirect the user to the corresponding target URL. If no match is found, the server will return a "Not Found" response.



# License
This project is licensed under the MIT License. See the LICENSE file for more details.





