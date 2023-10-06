To generate a regular expression pattern using gorx, follow these steps:

### Setup

1. Make sure you have the `OPENAI_API_KEY` environment variable set with your OpenAI API key. If it's not set, the app will prompt you to enter it.

2. Install the app by running the following command:
   ```bash
   go install github.com/niuguy/gorx
   ```
   This will install the app in your `$GOPATH/bin` directory. Make sure that directory is in your `$PATH` environment variable.

   Alternatively, you can clone the repo and run the app directly from the command line:
    ```bash
    go run main.go
    ```
    or build the app and run it:
    ```bash
    go build -o gorx main.go
    ./gorx
    ```


### Usage

The `gorx` command allows you to generate regular expression patterns in different modes. You can use the following flags to specify the mode:

- `-s, --semantic`: Generate a regular expression based on semantic meaning .
- `-p, --pattern`: Find common patterns of the given strings separated by empty space.
- `-c, --context`: Match the first input string from the second input string.

'-s' is the default mode if no flag is specified.

### Example Usage

- To find common patterns in a set of strings:
  ```
  gorx -p "string1" "string2" "string3"
  ```

- To match one string in the context of another:
  ```
  gorx -c "string_to_match" "context_string"
  ```

- To generate a regular expression pattern based on semantic meaning:
  ```
  gorx -s "ip address"
  ```

  this is equivalent to:
  ```
  gorx "ip address"
  ``` 

### Output

The app will make a request to the OpenAI API and return the generated regular expression pattern.



