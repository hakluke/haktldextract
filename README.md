# haktldextract

Basic tool to extract domains/subdomains from URLs en masse

### Installation

```
go install github.com/hakluke/haktldextract@latest
```

### Usage
```
cat urls.txt | haktldextract
```

Options:
```
-t      threads (number of concurrent threads to use, default is 16)
-s      subdomains (dump subdomains instead of base domains) 
```

Example:
```
cat urls.txt | haktldextract -s -t 16 | tee ./subs.txt
```
