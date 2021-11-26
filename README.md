# ado

This is a tool to extract TODOs, NOTEs etc or user provided terms from given files and/or directories.

**DEPRECIATED**: My project [seek](https://gitlab.com/utt_meelis/seek) has cleaner code and includes the core functionality/idea of this project. However, will keep this project here, since it has some nuances, that `seek` does not have.

## Usage

ado [SUBCOMMAND | [FLAG=\<value>] [filname(s) | dirname(s)]]

### SUBCOMMANDS

**help**

```
Prints help.
```

### FLAGS

**-file=filename**

```
The filename where the TODOs, NOTEs etc are saved in the current working directory.
```

**-indent=int**

```
The size of indentation between filepath and ado. (default 60)
```

**-ignore=filename**

```
Ignore file, where each line represents one directory or file that is ignored. If when .adoignore exist in the current directory, this flag is not necessary. (default .adoignore)
```

**-search=string**

```
Search for a specific string (regexp allowed); will overwrite the default search keywords (see keywords).
```

**-add=filename**

```
File, where each line represents one additional keyword (regexp allowed).
```

**-depth=int**

```
The depth of directory structure recursion, -1 is exhaustive recursion. (default -1)
```

### DEFAULT SEARCH KEYWORDS

```
TODO:|NOTE:|HACK:|DEBUG:|FIXME:|REVIEW:|BUG:|TEST:|TESTME:|MAYBE:
```

## Examples

- Gets recursively all ados from current working directory.
    ```sh
    ado
    ```
- Prints help for `ado` program.
    ```sh
    ado help
    ```
- Gets recursively all ados from current and parent directory.
    ```sh
    ado . ..
    ado ./ ../
    ado . ../
    ado ./ ..
    ```
- Gets recursively all ados from parent directory and <file1>.
    ```sh
    ado <file1> ../
    ```
- Gets all ados from <file1> and in addition saves them to test.txt in the current working directory.
    ```sh
    ado -file=test.txt <file1>
    ado --file=test.txt <file1>
    ```
- Gets recursively all ados from current directory and prints them so that ados start at column 100.
    ```sh
    ado -indent=100
    ado --indent=100
    ```
- Gets recursively all ados from current directory, ignoring files and directories mentioned in the given file.
    ```sh
    ado -ignore=.adoignore
    ado --ignore=.adoignore
    ```
- Gets recursively all 'RandomString' mentions from current directory.
    ```sh
    ado -search=RandomString
    ado --search=RandomString
    ```
- Gets recursively all ados from current directory, including the keywords mentioned in the given file.
    ```sh
    ado -add=.adoadd
    ado --add=.adoadd
    ```
- Gets recursively all ados from current directory, until subdirectory depth is 2 (including).
    ```sh
    ado -depth=2
    ado --depth=2
    ```

### Example output

```
/path/to/directory/ado/ado.go:50:                HACK: not the most elegant solution, but will do for now.
```

## Author

Meelis Utt
