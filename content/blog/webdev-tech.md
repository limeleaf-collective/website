+++
title = "Go and Rust for the Small Web: How We Build Faster, Simpler Apps"
date = 2025-02-26
draft = false
authors = ["Blain Smith"]

[taxonomies]
tags = ["Engineering"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

At Limeleaf, we chose to specialize in Go and Rust for our clients and [we've written about][9] why we made that choice. However, now that we're building our own products, we need to develop applications for the web, too.

<!-- more -->

## TL;DR

- Use [`templ`][0] in Go and [`maud`][1] for Rust to do server-side HTML rendering
- Use [`net/http`][2] in Go and [`axum`][3], [`actix`][4], [`rocket`][5], or whatever Rust crate for HTTP handling and routing
- Start with and try to stay with [SQLite][6] and migrate to [PostgreSQL][7], if needed

---

Go and Rust are great languages for systems programming, networking, and high-performance services, but not many folks think to use them for web applications that render HTML and handle form data. Most modern web applications employ the single page application (SPA) model. In these apps, a single web page is served to the user via a REST or GraphQL API (usually in JSON) and rendered in the user agent (usually a browser) by a separate frontend application written in a JavaScript framework like React or Vue.

For the products we're building, we want to keep the user's experience as simple and fast as possible. The SPA model is over-engineered for these kinds of apps. Instead, we implement most of the application on the server and send the UI over the wire in pure HTML, using only the bare minimum of JavaScript and CSS.

## Core Web Application Components

These are the only components we need to build a usable application:

- *Presentation* component to display information and accept user input
- *Persistence* component to durably store data for future access between sessions
- *Session/logic* component to manage access, but also to validate and manipulate data between the Presentation and Persistence components

You might be thinking, *isn’t what you just described simply an SPA*? After all, the components seem to map:

* Presentation -> HTML and CSS
* Persistence -> Database / File System
* Session/Logic -> REST / GraphQL API

You’re right, this could describe an SPA. However, we believe separating the Presentation and Session/Logic components in SPAs adds unnecessary complexity for most web applications. In the end, an SPA results in plain HTML and CSS at render time, but with that unnecessary layer in between. While there are certainly cases where the user experience justifies this approach, we argue that standard web protocols are more often a better choice for several reasons:

1. Simpler tech stacks can be shared more easily across smaller teams, requiring less cognitive load to grok.
2. Removing data transposing from JSON <-> HTML saves compute and loading times.
3. Distribution and development become easier and less fragile since there is no extra transpiling step (these steps are usually far slower than compilation steps).
4. JavaScript is not type-safe (and no, neither is TypeScript).
5. HTML and CSS are well-defined standards. They are thoroughly documented, tested, and can be coded by engineers of any skill level.

## Go and Rust in Web Applications

Let's circle back to Rust and Go. We are backend and systems engineers (and a product manager). All of us (even John) can understand and write enough HTML and CSS to produce highly functional but simple web applications. Our foundational practices are:

1. Render all HTML server-side and send it to the user agent to display information and form inputs.
2. Leverage CDNs for 3rd party CSS (and JavaScript if absolutely necessary).
3. Use standard `<form method="post" action="[route]">` elements to send user input to the server for validation and storage.
4. Generate type-safe template, embed, and static content directly into a single binary executable.

The first three are relatively straightforward, but number four is a big one for us. It allows us to distribute and deploy entire products in *one* file that contains all the HTML, CSS, and images necessary to run the application. Compare that to something like PHP or Python, where we'd have to manage and distribute hundreds or thousands of smaller files.

### Generating Templates in Go with `templ`

While Go has a standard template library package [`html/template`][8], we use [`templ`][10] instead, because it transpiles HTML into Go and then compiles it into the final binary. It allows us to write simple HTML forms like this:

```go
templ RegisterSignInForm() {
	<form method="post" action="/account">
		<p><input type="email" name="email" placeholder="Email" /></p>
		<p><button type="submit">Submit</button></p>
	</form>
}
```

After running `$ templ generate`, we end up with type-safe Go code:

```go
func RegisterSignInForm() templ.Component {   
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {                return templ_7745c5c3_CtxErr    
		}                                
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<form method=\"post\" action=\"/account\"><p><input type=\"email\" name=\"email\" placeholder=\"Email\"></p><p><button type=\"submit\">Submit</button></p></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}
```

Not the prettiest output, but it is generated code; we never read or edit it, so it doesn't matter. Instead, we can focus on writing the template code as simple, easy-to-understand HTML. `templ` supports much more functionality, so I urge you to check out their documentation for more details.

### Generating Templates in Rust with `maud`

Since Rust's standard library is much more focused on systems programming, there isn't much in the way of templates, let alone HTML templates. However, [`maud`][11], is a phenomenal package that offers macros to write HTML-like markup to compile and embed into the final compiled binary. Here's an example:

```rust
pub fn register() -> Markup {
    html! {
        form method="post" action="/register" {
            p {
                label for="name" { "Full Name" }
                input type="text" name="name";
            }
            p {
                label for="email" { "Email" }
                input type="email" name="email";
            }
            p {
                label for="password" { "Password" }
                input type="password" name="password";
            }
            p {
                input type="checkbox" name="tos";
                label for="tos" { "I agree to Limecast's Terms of Service" }
            }
            p {
                button type="submit" { "Register" };
            }
        }
    }
}
```

Although this is not standard HTML like in the Go version, it still gives us a familiar markup that captures the same intent. Since `html!` is a Rust macro, there is no need for a separate code generation step like Go needs with `templ`. Rust will automatically convert the macro into Rust code and compile it when you run `$ cargo build`.

### HTTP Handlers and Routing

With ways to generate and compile HTML server-side in each language, we can use just about any HTTP module to deliver it to the user since `templ` and `maud` offer direct support for writing the results of their respective templates to standard HTTP handlers. Here are examples in each:

```go
func (s *Server) RootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobs, err := ...
		if err != nil {
			templates.Error(err).Render(r.Context(), w)
			return
		}

		templates.Root(jobs).Render(r.Context(), w)
	})
}
```

```rust
pub async fn homepage(
    Extension(pool): Extension<SqlitePool>,
    Path(slug): Path<String>,
    Query(params): Query<HomepageParams>,
) -> Result<Markup, AppError> {
    let podcast: Podcast = ...;
    let episodes: Vec<Episode> = ...;

    Ok(templates::homepage(podcast, episodes))
}
```

### Persistence with SQLite (or PostgreSQL)

All the products we're working today on use SQLite because, well, it works perfectly fine for the size, scale, and load of the application. Both also use just boring old SQL, an easy way to save data to disk in a well-known language. If we ever need to scale up the service and split the database to a dedicated server, we can just migrate SQLite into PostgreSQL and re-deploy. However, we believe SQLite will serve our purposes for a long time. Since it now supports [WAL-mode][12] it can do much higher throughput, if we need it.

## Our Web Apps

We're working on two web app products today. [Apply.coop][13] is a job board for coops, and [Limecast][14] is a podcasting platform.

All of the code for these products, which practice what we preach, are in our [Codeberg](https://codeberg.org/limeleaf) repo. We open-source as much code as we can to be transparent about what we build and how we build it.

[0]: https://templ.guide
[1]: https://maud.lambda.xyz
[2]: https://pkg.go.dev/net/http
[3]: https://docs.rs/axum
[4]: https://docs.rs/actix
[5]: https://docs.rs/rocket
[6]: https://sqlite.org
[7]: https://www.postgresql.org
[8]: https://pkg.go.dev/html/template
[9]: https://limeleaf.coop/blog/why-go-and-rust/
[10]: https://github.com/a-h/templ
[11]: https://maud.lambda.xyz
[12]: https://sqlite.org/wal.html
[13]: https://apply.coop
[14]: https://limecast.com
