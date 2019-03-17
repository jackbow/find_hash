![output](https://i.imgur.com/t9YtwnT.gif)

| Language | Time to check 10mil hashes  | Speedups                 |
| :-----   | :-------------------------- | :----------------------- |
| Python   | 3882.30sec (~1hr)           | NA                       |
| Go       | 4.57sec                     | 277x Python              |
| Rust     | 1.66sec                     |  2.75x Go; 2,388x Python |

Reimplemented twice because I was curious about go, rust, and language speed differences. This was my first time trying go, rust, and multithreaded programming outside my OS class. Each reimplementation improves upon the last (better code, multithreading, live output, etc).

Used to find strings with md5 hashes which when converted to ASCII make unsanitized sql statements server-side evaluate to true regardless of correct password. Used for a sql injection as part of a web security project in my computer security class.
