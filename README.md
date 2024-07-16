# limeleaf.io

The website built with [Zola][1].

## Running locally

1. Ensure Zola is [installed][2]
2. Clone the repo
3. Run `zola serve` from the directory containing the `config.toml` file

The site will be running at [http://localhost:1111][3]. You can edit
any file in the repo and Zola will detect changes and automatically
rebuild the site and refresh the browser (if open).

## Deploying live

Any changes pushed to the `main` branch will automatically be pushed to
[https://limeleaf.net][4]. If you don't want that to happen then issue a
PR with the changes or set the Front Matter (the fields between the 
`+++` characters to contain `draft = true`), but a PR is probably just
easier.

```
+++
title = "This is a test page"
draft = true
+++

The markdown content of the page...
```

[1]: https://www.getzola.org
[2]: https://www.getzola.org/documentation/getting-started/installation/
[3]: http://localhost:1111
[4]: https://limeleaf.net
