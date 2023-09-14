# Completion exercise

Write a Bash program that shows summary information from wikipedia. 
`wiki walrus` should show the first sentence of the "Walrus"
article, along with a list of section headings.  `wiki walrus anatomy`
should show the first sentence of the anatomy section, along with a
list of subsections headings.

* You are welcome to use `curl` to retrieve the page itself, as well as
any commands that you'd expect to find on a typical Unix family OS
distribution.

* As a stretch goal, add support for tab completion in Bash, i.e.
`wiki walr[tab]` should expand to walrus, and `wiki walrus anat[tab]`
should expand to anatomy.

* As another stretch goal, add a flag to support querying multiple pages
in parallel.
