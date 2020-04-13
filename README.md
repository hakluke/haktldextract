# haktldextract

Basic tool to extract base domains from URLs.

# Usage
```
cat urls.txt | haktldextract
```

Options:
```
-t      threads (number of concurrent threads to use)
-s      subdomains (dump subdomains instead of base domains) 
```

Example:
```
cat urls.txt | haktldextract -s -t 16 | tee ./subs.txt
```
