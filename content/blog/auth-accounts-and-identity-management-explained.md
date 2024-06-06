+++
title = "Auth, Accounts, and Identity Management Explained"
date = 2024-05-06
draft = true

[taxonomies]
tags = ["Engineering", "Product"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

Most, if not all, products and services you use in your digital life
requires you to identify yourself in order to provide specific features.
You tend not to think about it much because has become so normalized.

<!-- more -->

1. Enter the address of the service into your browser
2. Find the "Sign-In", "Login", "Register", "Sign-Up" links somwhere at
the top right of the page
3. Fill in the Email/Username and Password fields on the form
4. Click the button and voila you're whisked away into your account,
dashboard, home feed, or whatever

Step 3 might be a bit different too depending on how the service allows
you to identify yourself.

1. Traditionally email/username and password
2. Use another account like Google, Facebook, GitHub, etc.
3. Use just and email and they send you a one-time link to that address
to follow

No matter what mechanism is used the result is the same. You have
provided the required information to prove that you are who you say you
are. Some services opt to give you a choice of all 3 to choose from when
you decide to use them.

At first glace they all seem fairly self explainitory, but now you want
to build your own service and build a userbase so you have to decide
which one, some, or all of the above you should offer folks. So which
should you choose, what does each entail to implement, and what are the
trade-offs over one or the other? Since we've been implementing auth
for decades we thought we could unpack some of these questions and cover
the basic building blocks potential users of your service should expect
when using your service.

## Simple Auth

**Go Example**
```go
func RegistraionHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()

  email := r.FormValue("email")
  password := r.FormValue("password")

  if userDB.AlreadyRegistered(email) {
    http.Redirect(w, r, "/register?error=User already exists", http.StatusSeeOther)
    return
  }

  encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

  _ := userDB.Create(email, encryptedPassword)

  http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

func SignInHandler(w http.ResponseWriter, r http.Request) {
  r.ParseForm()

  // Check db for creds

  http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
```

**Rust Example**
```rust
#[derive(Deserialize, Debug)]
struct Credentials {
    email: String,
    password: String,
}

async fn simple_auth_handler(Form(input): Form<Credentials>) {
  // Check db for creds

  Redirect::to("/dashboard")
}
```

## OAuth

## Magic Links
