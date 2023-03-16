import jwtDecode from "jwt-decode";
import constants from '../constants'

/*
Returns the data of the auth token JWT if the user is signed in, null otherwise.
*/
export const getSignedInUser = () => {
    const auth = getAuthJWT();
    if (auth) {
        return jwtDecode(auth);
    }
    return null;
}

export const getAuthJWT = () => {
    return localStorage.getItem(constants.AUTH_HEADER);
}

export const setAuthJWT = (authToken) => {
    localStorage.setItem(constants.AUTH_HEADER, authToken);
}

export const deleteAuthJWT = () => {
    localStorage.removeItem(constants.AUTH_HEADER);
}

/*
Returns true if the password is valid, false otherwise.
*/
export const validatePassword = (password) => {
    return password.length >= 8;
}