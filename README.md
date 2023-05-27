# Draugr: Toy Search Engine

Draugr is a toy search engine written in Go. It uses TF-IDF (Term Frequency-Inverse Document Frequency) to rank documents based on their relevance to a given query. The search engine is designed to work with text files, such as stories and essays.

## Getting Started


1. Install the executable:

   ```
   go install
   ```
2. Display help text:

   ```
   $ draugr -help
   Usage of draugr:
     -client
         Run as client
     -debug
         Print a lot of stuff
     -dir string
         index dir (default ".")
     -exts string
         file extensions to filter for (default ".txt .md .scala .go .hs .ts")
     -help
         Show help
     -search string
         search terms
     -serve
         Run as service

   ```
3. Run a local search:

   ```
   $ draugr -search water
   test_files/essays/The-Hudson-River-And-Its-Early-Names-Susan-Fenimore-Cooper.txt
   test_files/stories/2BR02B-Kurt-Vonnegut.txt
   ```

4. Index a directory in the background:

   ```
   draugr -serve -dir ~/code/go &
   ```

5. Search it:

   ```
   $ draugr -client -search water
   /home/bballant/code/go/draugr/test_files/essays/The-Hudson-River-And-Its-Early-Names-Susan-Fenimore-Cooper.txt
   /home/bballant/code/go/draugr/test_files/stories/2BR02B-Kurt-Vonnegut.txt
   ```

## Testing

To run the test suite, execute the following command:

```
go test ./...
```

## License

This project is licensed under the [MIT License](./LICENSE).
