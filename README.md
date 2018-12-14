# tplwalker

Package tplwalker is a tool to copy and transform a directory of templates
and other files to a target directory.  It will re-create subdirectories,
regular files and symlinks from a source directory.  Files that end with a
configurable suffix, will be stripped of that suffix and be digested as a Go
text template and executed with some data object.

It was once intended as part of an engine that automatically creates
directories of configuration files, but it has not been seen in action, yet.
