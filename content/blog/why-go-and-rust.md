+++
title = "Why Go and Rust?"
date = 2024-09-27
authors = ["Blain Smith"]

[taxonomies]
tags = ["Engineering"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

I was about 80% finished with this post and I was going down the road of
listing out all the great things both Go and Rust offer for us to build
and maintain software for our clients and product, but then I threw it
all out. I realized that posts like that are a dime a dozen and there
are so many resources out there comparing, contrasting, and listing
features of both languages, but none of that really explains WHY we, the
humans Erik, John, and I, like using Go and Rust. So instead I thought
it would be better to personalize this a bit and grab some quotes from
each of us about why we like each of them, unedited, unfiltered, and
raw!

<!--more-->

<figure style="text-align: center;">
  <img style="height: 200px;" src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" alt="Go logo: Courteousy of https://go.dev" />
</figure>

Erik says:

> My first experience writing production code was with PHP and it wasn’t
> until I was exposed to Ruby that I had my first “ah-ha!” moment that
> there was a much more pleasant way to do things. Go was that second
> moment. The standard library and its documentation (and generally the
> community package documentation), static typing, and overall
> simplicity of the language instantly clicked and felt like it made me
> a more productive developer while making systems that were performant
> and correct. Removing the need to pull in a myriad of external
> dependencies and libraries and, in turn, having the ability to rely on
> Go’s standard library for networking, I/O and testing has helped me
> many times to quickly prototype a system I’m confident will work the
> way I intend and is generally already production-ready.

---

John says:

> Go revolutionized the development of network services and web
> applications. Its comprehensive standard library, simplicity, and
> integrated testing tools (among other features) enable teams to
> rapidly build and refine sophisticated, highly scalable network
> software products with minimal bugs–things that no Product Manager
> will ever complain about.

---

Blain says:

> Go will always hold a special place in my heart. Switching from Node
> to Go all those years ago made me want to return to college and focus
> more on writing software with respect to the machines it is running
> on. Even though it is a garbage collected language you can write Go
> efficiently by understanding how and when memory is allocated and
> deallocated. However, the biggest selling point for Go for me is the
> amazing standard library that makes writing networking and I/O code
> such a breeze with net, net/http, io, and bufio. I remember standing
> by that once when I decided to live code a toy key/value database in
> Go at a conference. Just a blast to write.

<figure style="text-align: center;">
  <img style="height: 200px;" src="https://www.rust-lang.org/static/images/rust-logo-blk.svg" alt="Rust logo: Courteousy of https://rust-lang.org" />
</figure>

Erik says:

> Learning Rust has given me a new perspective on development. It’s been
> another one of the “ah-ha!” moments along my software journey that has
> made me more mindful in terms of memory and the ownership of data
> within a program. Being able to write performant, safe code in a
> concise and simple way is such a joy. The somewhat steeper learning
> curve really pays off when you’re more confident and fearless in the
> code that you’re writing.

---

John says:

> As a Product Manager, I’m excited about Rust because it opens up a new
> world of interesting customers and projects for us at Limeleaf.
> Developers writing in modern languages like Rust are at the leading
> edge of innovation, using it to build everything from complex web apps
> to IoT devices. It’s really fun and rewarding to mentor companies as 
> they implement the language.

---

Blain says:

> I actually really love the borrow checker. I spent a lot of years
> writing backend systems with dynamic languages and I still have
> nightmares about null references and panics in production. Even though
> I like Go, it has its pain points with nil references as well. Rust’s
> memory model in general is just fantastic at getting rid of an entire
> class of problems and actually makes me a better programmer when
> thinking about ownership. A language that can do all of that to be
> used in not only distributed web applications and systems, but also
> embedded microcontroller systems is a win/win/win in my book.
