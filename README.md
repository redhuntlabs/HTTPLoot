# HTTPLoot
An automated tool which can simultaneously crawl, fill forms, trigger error/debug pages and "loot" secrets out of the client-facing code of sites.

## Usage
To use the tool, you can grab any one of the pre-built binaries from the [Releases](https://github.com/redhuntlabs/HTTPLoot/releases) section of the repository. If you want to build the source code yourself, you will need Go > 1.16 to build it. Simply running `go build` will output a usable binary for you.

Additionally you will need two json files ([lootdb.json](https://github.com/redhuntlabs/HTTPLoot/blob/master/lootdb.json) and [regexes.json](https://github.com/redhuntlabs/HTTPLoot/blob/master/regexes.json)) alongwith the binary which you can get from the repo itself. Once you have all 3 files in the same folder, you can go ahead and fire up the tool.

Video demo:

[![video](https://user-images.githubusercontent.com/39941993/168653593-9551b6be-0eb7-4fa8-85ee-0de8e4506fe6.png)](https://www.youtube.com/watch?v=qc8Mm2O5t6Q)

Here is the help usage of the tool:
```s
$ ./httploot --help
      _____
       )=(
      /   \     H T T P L O O T
     (  $  )                  v0.1
      \___/

[+] HTTPLoot by RedHunt Labs - A Modern Attack Surface (ASM) Management Company
[+] Author: Pinaki Mondal (RHL Research Team)
[+] Continuously Track Your Attack Surface using https://redhuntlabs.com/nvadr.

Usage of ./httploot:
  -concurrency int
        Maximum number of sites to process concurrently (default 100)
  -depth int
        Maximum depth limit to traverse while crawling (default 3)
  -form-length int
        Length of the string to be randomly generated for filling form fields (default 5)
  -form-string string
        Value with which the tool will auto-fill forms, strings will be randomly generated if no value is supplied
  -input-file string
        Path of the input file containing domains to process
  -output-file string
        CSV output file path to write the results to (default "httploot-results.csv")
  -parallelism int
        Number of URLs per site to crawl parallely (default 15)
  -submit-forms
        Whether to auto-submit forms to trigger debug pages
  -timeout int
        The default timeout for HTTP requests (default 10)
  -user-agent string
        User agent to use during HTTP requests (default "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0")
  -verify-ssl
        Verify SSL certificates while making HTTP requests
  -wildcard-crawl
        Allow crawling of links outside of the domain being scanned
```

### Concurrent scanning
There are two flags which help with the concurrent scanning:
- `-concurrency`: Specifies the maximum number of sites to process concurrently.
- `-parallelism`: Specifies the number of links per site to crawl parallely.

Both `-concurrency` and `-parallelism` are crucial to performance and reliability of the tool results.

### Crawling
The crawl depth can be specified using the `-depth` flag. The integer value supplied to this is the maximum chain depth of links to crawl grabbed on a site.

An important flag `-wildcard-crawl` can be used to specify whether to crawl URLs outside the domain in scope.

> __NOTE__: Using this flag might lead to infinite crawling in worst case scenarios if the crawler finds links to other domains continuously.

### Filling forms
If you want the tool to scan for debug pages, you need to specify the `-submit-forms` argument. This will direct the tool to autosubmit forms and try to trigger error/debug pages _once a tech stack has been identified successfully_.

If the `-submit-forms` flag is enabled, you can control the string to be submitted in the form fields. The `-form-string` specifies the string to be submitted, while the `-form-length` can control the length of the string to be randomly generated which will be filled into the forms.

### Network tuning
Flags like:
- `-timeout` - specifies the HTTP timeout of requests.
- `-user-agent` - specifies the user-agent to use in HTTP requests.
- `-verify-ssl` - specifies whether or not to verify SSL certificates.

### Input/Output
Input file to read can be specified using the `-input-file` argument. You can specify a file path containing a list of URLs to scan with the tool. The `-output-file` flag can be used to specify the result output file path -- which by default goes into a file called `httploot-results.csv`.

## Further Details
Further details about the research which led to the development of the tool can be found on our [RedHunt Labs Blog](https://redhuntlabs.com/blog/the-http-facet-httploot.html).

## License & Version
The tool is licensed under the MIT license. See LICENSE.

Currently the tool is at v0.1.

## Credits
The RedHunt Labs Research Team would like to extend credits to the creators & maintainers of [shhgit](https://github.com/eth0izzle/shhgit) for the regular expressions provided by them in their repository.

**[`To know more about our Attack Surface Management platform, check out NVADR.`](https://redhuntlabs.com/nvadr)**
