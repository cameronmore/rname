# usage

`rname` takes two arguments and ignores anything beyond that. The first argument is the name or path of the thing to be renammed, and the second argument is the new name.

1. Both arguments must be the same kind of thing (e.g., both files, both directories). 
2. You may provide a path in your arguments, and if the path does not fully exist, it will be created
3. If the destination already exists, `rname` will not write over the existing path and will add '-duplicate' to the filename of any duplicates.
4. `rname` will not make copies of directories if the destination directory already exists. It will add the renammed files to the existing directories.

## examples

suppose I'm in my `dev/` directory and I want to rename a project from `my-react-app` to `my-vue-app`. Just run `rname my-react-app my-vue-app`.

What if you've already created `my-vue-app`? by default, `rname` will make a duplicate (`my-vue-app_duplicate`). optionally, `rname` will perform a merge of all the contents of `my-react-app` into the existing `my-vue-app` (by adding the `-m` flag). If there are any potential duplicate files in the directory merge, the ones coming from `my-react-app` like `index.html` will be renammed `index_duplicate.html`.
