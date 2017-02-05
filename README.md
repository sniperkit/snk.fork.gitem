# Gitem
Gitem is a command line tool that helps you make sense of the corpus of code that an organization or user has in github.  It supports
 the enumeration and checkout of all or a subset of the projects that a github organization or user has.

# Supported Authentication
## Basic
- username, password

## Oauth2 Token
- token

# Commands
## Clone your repositories
```
gitem clone repos -user="" -password=""
```

## Clone an owner's repositories
```
gitem clone repos -owner=""
```

## Clone an organization's repositories
```
gitem clone repos -org=""
```

## List your repositories
```
gitem list repos -user="" -password=""
```

## List an owner's repositories
```
gitem list repos -owner=""
```

## List an organization's repositories
```
gitem list repos -org=""
```

## List contributors
```
gitem list contributors -owner="" -repo=""
```

## List languages
```
gitem list languages -owner="" -repo=""
```

## List teams
```
gitem list teams -owner="" -repo=""
```

## List tags
```
gitem list tags -owner="" -repo=""
```