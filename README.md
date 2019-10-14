## Jokester
1. Place /jokeServer in your Go workspace.  
2. Build the executeable found in /cmd/jokeServer.  
3. Executeable can be ran with an optional parameter of "port" defining the port to listen for requests on, defaulting to 5000. 
   Example: ./jokeServer -port=5000
4. Request a joke: curl http://localhost:5000 (replace 5000 with set port if cli parameter set)
