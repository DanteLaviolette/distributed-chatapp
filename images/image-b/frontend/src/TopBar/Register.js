import { Box, Button, CircularProgress, FormControl, FormLabel, Input, Modal, ModalDialog, Typography } from "@mui/joy";
import { Stack } from "@mui/system";
import axios from "axios";
import PropTypes from 'prop-types';
import { useState } from "react";
import constants from "../constants";
import { getSignedInUser, validatePassword } from "../utils/utils";
import authorizedAxios from "../utils/AuthInterceptor";

Register.propTypes = {
    setUser: PropTypes.func
}

/*
Registration button that opens a modal for the user to register on click.
Calls props.setUser upon successful registration & login.
*/
function Register(props) {
    // Form related values
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [errorMessage, setErrorMessage] = useState("")
    const [isLoading, setIsLoading] = useState(false)
    // Input values
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const [firstName, setFirstName] = useState("")
    const [lastName, setLastName] = useState("")

    // Attempts to register user, either setting the errorMessage
    // Or calling props.setUser and exiting
    const register = () => {
        setErrorMessage("")
        // Ensure password is valid
        if (password !== confirmPassword) {
            setErrorMessage("Passwords not equal")
            return
        } else if (!validatePassword(password)) {
            setErrorMessage("Passwords must be 8 characters")
            return
        }
        setIsLoading(true)
        // Attempt to register
        axios.post("/api/register", {
            "email": email,
            "firstName": firstName,
            "lastName": lastName,
            "password": password,
        }).then(res => {
            // Success -- attempt to login
            authorizedAxios.post("/api/login", {
                "email": email,
                "password": password
            }).finally(() => {
                // Using finally, as in the case that logging in failed
                // The user can manually sign in later
                // Set user & close modal
                props.setUser(getSignedInUser())
                setIsModalOpen(false);
            })
        }).catch(err => {
            // Error occurred -- set error message
            if (err.response.status === 409) {
                setErrorMessage("Email already registered")
            } else if (err.response.status === 400) {
                setErrorMessage(err.response.data)
            } else {
                setErrorMessage("Registration Failed. Try again later.")
            }
            setIsLoading(false)
        })
    }

    return (
        <div>
            <Button size="sm" color="neutral" variant="plain" onClick={() => setIsModalOpen(true)}>Register</Button>
            <Modal open={isModalOpen} onClose={() => setIsModalOpen(false)}>
                <ModalDialog
                    sx={{ maxWidth: 500, minHeight: 'fit-content', overflow: 'scroll' }}
                >
                    <Typography level="h5">
                        Register New Account
                    </Typography>
                    {isLoading && <Box sx={{ display: 'flex', justifyContent: 'center' }}><CircularProgress /></Box>}
                    {!isLoading && <div>
                        <Typography level="body2">
                            {constants.PASSWORD_REQUIREMENTS}
                        </Typography>
                        <form
                            onSubmit={(event) => {
                                event.preventDefault();
                                register()
                            }}
                        >
                            <Stack spacing={1}>
                                <FormControl>
                                    <FormLabel>Email</FormLabel>
                                    <Input type="email" autoFocus required value={email} onChange={e => setEmail(e.target.value)} />
                                </FormControl>
                                <FormControl>
                                    <FormLabel>First Name</FormLabel>
                                    <Input required value={firstName} onChange={e => setFirstName(e.target.value)} />
                                </FormControl>
                                <FormControl>
                                    <FormLabel>Last Name</FormLabel>
                                    <Input required value={lastName} onChange={e => setLastName(e.target.value)} />
                                </FormControl>
                                <FormControl>
                                    <FormLabel>Password</FormLabel>
                                    <Input type="password" required value={password} onChange={e => setPassword(e.target.value)} />
                                </FormControl>
                                <FormControl>
                                    <FormLabel>Confirm Password</FormLabel>
                                    <Input type="password" required value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} />
                                </FormControl>
                                {errorMessage !== "" && <Typography color="danger" level="body3">{errorMessage}</Typography>}
                                <Button type="submit">Register</Button>
                            </Stack>
                        </form>
                    </div>}
                </ModalDialog>
            </Modal>
        </div>
    );
}

export default Register;
