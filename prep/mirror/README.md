# web mirror

## Problem statement (Exercise 8.7)

Write a concurrent program that creates a local mirror of a web site,
fetching each readable page and writing it to a directory on the local disk. 
Only pages within the original domain (for instance, golang.org) should be 
fetched.  URLs within mirrored pages should be altered as needed so that 
they refer to the mirrored page, not the original.