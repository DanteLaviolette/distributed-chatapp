import { Box, Button, CircularProgress, FormControl, FormLabel, Input, Modal, ModalDialog, Typography } from "@mui/joy";
import { Stack } from "@mui/system";
import axios from "axios";
import PropTypes from 'prop-types';
import { useState } from "react";
import { getSignedInUser } from "../utils";

SignIn.propTypes = {
    setUser: PropTypes.func
}

/*
Login button that displays a login form as a modal.
Calls props.setUser upon successful login.
*/
function SignIn(props) {
    // Form related values
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [errorMessage, setErrorMessage] = useState("")
    const [isLoading, setIsLoading] = useState(false)
    // Input values
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")

    // Attempts to register user, either setting the errorMessage
    // Or calling props.setUser and exiting
    const register = () => {
        setErrorMessage("")
        setIsLoading(true)
        // Attempt to register
        axios.post("/api/login", {
            "email": email,
            "password": password,
        }).then(res => {
            props.setUser(getSignedInUser())
            setIsModalOpen(false);
        }).catch(err => {
            // Error occurred -- set error message
            if (err.response.status === 400) {
                setErrorMessage("Invalid credentials.")
            } else {
                setErrorMessage("Login Failed. Try again later.")
            }
            setIsLoading(false)
        })
    }

    return (
        <div>
            <Button size="sm" sx={{ marginRight: "10px" }} variant="plain" onClick={() => setIsModalOpen(true)}>Sign In</Button>
            <Modal open={isModalOpen} onClose={() => setIsModalOpen(false)}>
                <ModalDialog
                    sx={{ maxWidth: 500, minHeight: 'fit-content', overflow: 'scroll' }}
                >
                    <Typography level="h5">
                        Sign In
                    </Typography>
                    {isLoading && <Box sx={{ display: 'flex', justifyContent: 'center' }}><CircularProgress /></Box>}
                    {!isLoading && <div>
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
                                    <FormLabel>Password</FormLabel>
                                    <Input type="password" required value={password} onChange={e => setPassword(e.target.value)} />
                                </FormControl>
                                {errorMessage !== "" && <Typography color="danger" level="body3">{errorMessage}</Typography>}
                                <Button type="submit">Sign In</Button>
                            </Stack>
                        </form>
                    </div>}
                </ModalDialog>
            </Modal>
        </div>
    );
}

export default SignIn;
