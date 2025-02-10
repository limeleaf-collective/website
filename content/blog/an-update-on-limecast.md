+++
title = "An Update on Limecast"
date = 2024-11-14
draft = false
authors = ["John Luther"]

[taxonomies]
tags = ["Limecast", "Product"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

Online content distribution has become highly centralized and increasingly dominated by a few tech monopolies.

We believe these companies prioritize profit and growth over user privacy and freedom of speech. What's more concerning is that their leaders are [cozying up to politicians](https://prospect.org/politics/2024-10-22-silicon-valley-billionaires-supporting-trump/) who have [repeatedly promised to crack down on the press and free expression](https://www.theatlantic.com/politics/archive/2024/11/donald-trump-hates-free-speech/680515/).

<!-- more -->

We believe this trend is dangerous, and we want to do our small part to resist online censorship. To that end, we've shifted the development plan for [Limecast](https://limecast.net), the podcasting platform we're working on. 

We will offer Limecast as a paid service someday, but our new priority is to develop it into a free, open-source, privacy-first, decentralized podcast management, distribution, and monetization platform.

We hear you thinking, *That is a lot of words, what do they mean?* In brief:

- Mandatory HTTPS transport between podcast storage and listeners.
- Provide a decentralized content storage and distribution option. We're still discussing what technology to use for this ([IPFS](https://ipfs.tech/), [Storj](https://github.com/storj/storj), [Iroh](https://www.iroh.computer/), etc.) and [welcome your ideas](https://codeberg.org/limeleaf/limecast/issues). Local server filesystem and CDN-like storage will still be options, too.
- Simple deployment and management to be platform agnostic. Run it on Kubernetes or a Raspberry Pi and everywhere in between. The whole app will be contained in a single binary file.
- Written in [Rust](https://rust-lang.org) and backed by [SQLite](https://sqlite.org) to be type-safe, performant, and straightforward. Pure HTML, CSS, and minimal JS. No SPA, no bloated JS frameworks, no JSON API. Just plain old form POST-ing, cookies, and other web standards.
- Open standards for content syndication and discovery with RSS and ActivityPub.
- Privacy-first analytics to help podcasters understand their audience while not compromising that audience's privacy.
- Licensed under the [Mozilla Public License](https://www.mozilla.org/en-US/MPL/), which we feel is a reasonable middle ground between permissive licenses like MIT and strong copyleft licenses like the GPL.
- Integrate with [Open Collective](https://opencollective.org) so podcasters can get paid for their work.

It's the early days of the project, but we've implemented some basic functionality (creating podcasts and adding/editing episode metadata, etc.), which you can demo on your local system by cloning [the Limecast repository on Codeberg](https://codeberg.org/limeleaf/limecast).

To learn more, visit [limecast.coop](https://limecast.coop/). If you want to get involved, check out [the repo](https://codeberg.org/limeleaf/limecast)), support us on [Open Collective](https://opencollective.com/limeleaf/projects/limecast), or contact us at [info@limeleaf.coop](mailto:info@limeleaf.coop).
