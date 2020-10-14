# go-remove-slowly

A tool to remove file slowly, truncating.

Removing large files with normal `rm` command may cause high I/O load.
This tool truncates the target file gradually, and finally removes it to avoid high I/O load.
