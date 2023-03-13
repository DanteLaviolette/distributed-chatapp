import Cookies from 'js-cookie'
import jwt_decode from "jwt-decode";

/*
Returns the data of the auth token JWT if the user is signed in, null otherwise.
*/
export const getSignedInUser = () => {
    const auth_cookie = Cookies.get("auth");
    if (auth_cookie) {
        return jwt_decode(auth_cookie);
    }
    return null;
}