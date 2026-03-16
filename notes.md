Why did you told me to have a apiconfig struct ?
to be used alongside handlers?
what is this technique?


That apiConfig pattern is a handy little wizard trick called a dependency-carrying handler (or “struct with dependencies + methods as handlers”).

Why apiConfig?
You need your handlers to access shared stuff:

DB connection / sqlc Queries
config values (like PLATFORM)
maybe metrics, loggers, etc.
You could make all those globals, but that’s hard to test and reason about. Instead, you:

type apiConfig struct {
    db       *database.Queries
    platform string
    // other shared stuff...
}

Then your handlers become methods:

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
    user, err := cfg.db.CreateUser(r.Context(), ...)
    // use cfg.platform, etc.
}

And in main:

cfg := &apiConfig{
    db:       dbQueries,
    platform: os.Getenv("PLATFORM"),
}

http.HandleFunc("/api/users", cfg.handlerUsersCreate)

What’s this technique called?
Common names you’ll hear:

“Handlers as methods on a config struct”
“Dependency injection via struct”
“Passing dependencies through a context struct”
In essence, it’s a lightweight form of dependency injection: instead of hardcoding globals, you inject what your handlers need through apiConfig. That keeps your code more modular and testable.

If you’d like, I can show you how you’d test one of these handlers by building a fake apiConfig.



https://stackoverflow.com/questions/17306358/removing-fields-from-struct-or-hiding-them-in-json-response


why signingkey is given in []byte() type in the signingkey function
In computer science, cryptographic algorithms don't actually "know" what a string is. A string is a high-level concept meant for humans to read. Under the hood, a string is just a sequence of characters like "s", "e", "c", "r", "e", "t".

Mathematical operations like the ones used in HS256 (the signing method you're using) need to work with the raw, underlying bits and bytes of data to perform their calculations.

In Go:

A string is a read-only slice of bytes, but it's treated as a distinct type.
A []byte (byte slice) is the standard way to represent raw binary data.
The golang-jwt library requires a []byte because it's going to perform bitwise operations on that key to sign or verify the token. It doesn't care about the "letters" in your secret; it only cares about the raw binary values.

It's similar to how a scale doesn't care if you're weighing "gold," "feathers," or "lead"—it only cares about the raw mass. Converting to []byte is like stripping away the label so the algorithm can "weigh" the raw data.

######################################################

Handling nullstring when sqlc generates a parameter when a column can be NULL.
var email sql.NullString

if params.Email == "" {
    email = sql.NullString{
        Valid: false, // will write NULL
    }
} else {
    email = sql.NullString{
        String: params.Email,
        Valid:  true,
    }
}

err := cfg.db.UpdateEmail(r.Context(), database.UpdateEmailParams{
    Email: email,
    ID:    userId,
})
