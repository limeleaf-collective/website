+++
title = "WebDev Tech"
date = 2025-02-26
draft = false
authors = ["Blain Smith"]

[taxonomies]
tags = ["Engineering"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

At Limeleaf we chose to specialize in Go and Rust for our clients and we've covered that topic as to why we made that choice. However, since we're getting into building our own products we need to be able to develop applications for the web too.

<!-- more -->

## TL;DR

- Use [`templ`][0] in Go and [`maud`][1] for Rust to do server-side HTML rendering
- Use [`net/http`][2] in Go and [`axum`][3], [`actix`][4], [`rocket`][5], or whatever Rust crate for HTTP handling and routing
- Start with and try to stay with [SQLite][6] and migrate to [PostgreSQL][7] if needed

---

Go and Rust are great systems programming languages for networking and high performance services, but not many folks think to use them traditional web applications that render HTML and handle form data. Most modern web applications subscribe to the single page application (SPA) model where there exists some backend REST or GraphQL API and a completely separate frontend application written in JavaScript with some popular framework like React or Vue. 

We think for the products we're building and for the open web these choices are overkill.

## Core Web Application Components

If we break down what a web application really is to it's core components we end up with:

- Presentation component to display information and accept user input
- Session/Logic component to manage access, validate, and manipulate data between the Presentation and Persistence component
- Persistence component to durably store data on a disk for future access between sessions, failures, and restarts.

Now, you might be thinking that I just described exactly what the modern web applications are doing:

- Presentation -> React / Vue SPA
- Session/Logic -> REST / GraphQL API
- Persistence -> Database / File System

While you are correct we think that splitting up the Presentation and Session/Logic components is needlessly complicated for a large class of web applications since at the end of the day the SPA choice ultimately ends up being plain old HTML and CSS, but with an extra step. Sure, there are times where the UX might warrant the need for such a technology choice, but we'd argue that using standard web protocols is a better choice for a few reasons:

1. Simpler tech stacks can be shared easier across smaller teams requiring less cognitive load to grok
2. Removing intermediate data transposing from JSON <-> HTML saves compute and loading times.
3. Distribution and development become easier and less fragile since there is no extra transpiling step given that these steps are usually much slower than compilation steps.
4. JavaScript is not type safe and, no, neither is TypeScript.
5. Vanilla HTML and CSS are well defined standards, thoroughly documented, tested, and can be learned by all walks of life.

## Go and Rust in Web Applications

Now, let's circle back to our choice to specialize. We can expand our uses of Go and Rust for building web applications. We, as mostly backend and systems engineers (and a product manager), can all certainly learn and understand HTML and CSS enough to produce high quality web applications. The concrete ideas we practice are:

1. Render HTML server-side and send it to the browser to display information and form inputs.
2. Leverage CDNs for 3rd part external CSS (and JavaScript if absolutely necessary).
3. Use standard `<form method="post" action="*>` support in browsers to send user input to the server for validation and storage.
4. Generate type safe template and embed and static content directly into the compiled binary

The first 3 are relatively straightforward, but number 4 is a big one for us that we leverage heavily. This allows us to distribute and deploy a single binary file with all of the HTML, CSS, and images the application needs embedded directly into that same file and for the size of web applications we're building this is more than acceptable.

### Generating Templates in Go with `templ`

While Go does have its standard library package [`html/template`][8] we instead use `templ` which is a much more powerful choice to write HTML that gets directly transpiled into Go which then gets compiled into the final binary.  We can write simple HTML forms like:

```go
templ RegisterSignInForm() {
	<form method="post" action="/account">
		<p><input type="email" name="email" placeholder="Email" /></p>                   <p><button type="submit">Submit</button></p>
	</form>
}
```

Then after running `$ templ generate` we end up with type-safe Go code:

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

Sure, it may not look like the prettiest output, but it is generated code so it doesn't matter. We can focus on writing the template code as regular HTML. `templ` supports much more functionality so I urge you to go and checkout their documentation for more details.

### Generate Templates in Rust with `maud`

Since Rust's standard library is much more focused on systems programming there isn't much in the way of templates, let alone HTML templates. However, `maud`, is a phenomenal package that offers macros to write HTML-like markup to compile and embed into the final compiled binary. 

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

Even though this is not HTML like the Go version, we still have a familiar markup that captures the same intent. Since `html!` above is a Rust macro there is no need for a separate code generation step like Go needs with `templ`. Rust will automatically convert the macro into Rust code and compiling it once you run `$ cargo build`.

### HTTP Handlers and Routing

Now that we have a way to generate and compile HTML server-side in each language we can use any number of HTTP modules to deliver them to the browser since `templ` and `maud` offer direct support for writing the results of their respective templates to standard HTTP handlers.

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

Everything we've built so far use SQLite because, well, it works perfectly fine for the size, scale, and load of the application. Both choices are also just boring SQL and that is entirely the point here. We need a way to save data to disk with a well-known language and SQL does just that. If there comes a time where we need to scale up the service and split the database to a dedicated server then we can just migrate SQLite into PostgreSQL and re-deploy. However, SQLite should serve our purposes for a very long time and with the addition of running it in WAL-mode it allows for much higher throughput which serves our needs just fine as well.

## Our Web Apps

You can find all of our web applications that practice what we preach over on our [Codeberg](https://codeberg.org/limeleaf) page. Everything we do we try to keep open source to be transparent about what we build and how we build it.

[0]: https://templ.guide
[1]: https://maud.lambda.xyz
[2]: https://pkg.go.dev/net/http
[3]: https://docs.rs/axum
[4]: https://docs.rs/actix
[5]: https://docs.rs/rocket
[6]: https://sqlite.org
[7]: https://www.postgresql.org
[8]: https://pkg.go.dev/html/template
