# NovelParser

Parse multiple small chapters of a desired novel from various sites and compile them into an EPUB book with around 10 chapters for convenient reading.
Count of chapters can be setup and store for each novel in config

Command examples:
- `novelparser add novelname https://example.com/novelname/chapter-1 --batch-size 10`\
store new novel and first chapter in config.

- `novelparser parse novelname 10`\
parse 10 books, each with 10(batch-size) chapters.


