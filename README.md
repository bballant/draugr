# Draugr: Toy Search Engine

Draugr is a toy search engine written in Go. It uses TF-IDF (Term Frequency-Inverse Document Frequency) to rank documents based on their relevance to a given query. The search engine is designed to work with text files, such as stories and essays.

## Getting Started

TODO: This isn't how it works

1. Build the executable:

   ```
   go build
   ```

2. Index your documents:

   ```
   ./draugr index -dir test_files/stories
   ```

3. Start searching:

   ```
   ./draugr search -query "your search query"
   ```

## Testing

To run the test suite, execute the following command:

```
go test ./...
```

## License

This project is licensed under the [MIT License](./LICENSE).
