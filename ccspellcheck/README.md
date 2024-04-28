# Build Your Own Spell Checker Using A Bloom Filter

References:
- https://codingchallenges.fyi/challenges/challenge-bloom


## Build & Run

Create list of words that we want to insert into the Bloom Filter.

```console
$ cp /usr/share/dict/words words.txt
```

```console
$ make
```

```console
$ build/ccspellcheck homer
2024/04/28 19:59:41 buildDictionary: added 235976 words
2024/04/28 19:59:41 'homer' found!
```
