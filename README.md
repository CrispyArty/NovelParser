# NovelParser

Parse multiple chapters of a chosen novel from various sites and compile them into an EPUB for convenient reading.
Count of chapters can be setup and store for each novel in config

Commands:
-  `novelparser add [novel_name url] [options]`\
    `-b [batch_size]`  - will set how many chapters per book should be compiled.\
    Adds a novel to config. `url` - url of first chapter in novel

-  `novelparser parse [novel_name] [options]`\
    `-c [count]` - How many book you need to download. Default is 1\
    `-e` - Email downloaded books to kindle e-book\
    Parse amount of chapters to fill `count` books.\
    Books will be saved as epub format file in `"uploads/{novel_name}/Chapers {first}-{last}.epub"`


## Additional
Parser doesn’t use GoCoroutines(multithreading) for parsing multiple chapters at once because sites will throttle requests from your ip.\
But Go coroutines are used for email delivery.

Email delivery is done by oauth2 protocol.\t
If token is missing or expired you will be prompted to open a link to authorize NovelParser application.
