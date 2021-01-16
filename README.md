# normalize-srt

The purpose of this program is to normalize a specific kind of malformed srt file, which I downloaded from Internet.

It looks like:

```
1
1

00:00:00,180  -->  00:00:01,250
Okay, so now let's talk
2

2

00:00:01,250  -->  00:00:04,150
about the United states.
3

...
```

That's annoying, because my media player can't recognize it. Even though some players work, it will display a redundant number.

So far, I observed every timecode followed by a single line of caption text, so my algorithm is to extract every pair of "the valuable" from the source file and put it into an array. Then re-assign the sequence number when output. (plz pay attention if this isn't in your case.)

Then, this is expected:

```
1
00:00:00,340  -->  00:00:01,250
Okay, we need talk

2
00:00:01,250  -->  00:00:04,150
about the United states.

3
...
```

Also, this program can recursively process all srt files in a directory tree, by specifying a path.
