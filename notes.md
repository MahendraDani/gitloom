# The Git Way
- a git repository is made up of __Git Objects__ (there are only four types)
- all the objects are stored in __Git Object Database__ which is the __Git directory__ (`.git`)
- each object is compressed (with zlib) and referenced using its SHA-1 value

## Git Object Types
### 1. The Blob Object
- in git, contents of a file are stored as __blobs__.
  - only contents of files are stored in blobs, not the files.
  - the names and modes of files are not stored.
  - Why? __if you have two files anywhere in your project that are exactly the same, even if the have different names, Git will only store the blob once__.

### 2. The Tree Object
- directories in Git correspond to trees.
- A tree is a simple list of trees and blobs that the tree contains, along  with the names and modes of those trees and blobs.

```
tree [content size]\0
<filemode> <object_type> <sha1> <filename>

100644 blob ada123 README
100644 blob sdc13a hello.txt
040000 tree ace12s lib
```

### 3. The Commit Object
- It simply points to a tree and keeps an author, committer, message and any parent commits that directly preceded it.

```
commit [content size]\0

tree a123ed
parent (if exists) 12aac3
author <user_name>
          <user_email>
committer <user_name>
          <user_email>

<commit_message>
```


- parent is the SHA-1 of the last commit, so merging two branches will have SHA-1 of both parents.

## 4. The Tag Object
- This is an object that provides a permanent shorthand name for a particular commit. 
- It contains an object, type, tag, tagger and a message.