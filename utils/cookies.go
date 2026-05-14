package utils

/*

how should I use cookies?

during login, create session (via cookiestore)
when using .save it saves on client via cookie

purpose of cookiestore is all data stored in cookie, not on server side


when user accesses other sites, we need to validate them
middleware is best place to validate (intercept other handlers)

middleware to check if cookie matches that of user?



session_key ("mtg_app_session") <-- name of cookie
[loggedInUserKey] = user.ID		<-- value in the cookie
									(loggedInUserKey is "user_id")
This lets us retrieve cookie from request, check the value as its constant,
then validate the users id w/ DB comparison.
Also while not sharing users data in cookie.



*/
