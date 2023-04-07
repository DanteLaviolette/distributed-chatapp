import { Box, Button, CircularProgress, FormControl, FormLabel, Input, MenuItem, MenuList, Modal, ModalDialog, Typography } from "@mui/joy";
import PropTypes from 'prop-types';
import { styled } from '@mui/joy/styles';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import ClickAwayListener from '@mui/base/ClickAwayListener';
import PopperUnstyled from '@mui/base/PopperUnstyled';
import { useState } from "react";
import Cookies from "js-cookie";
import { Stack } from "@mui/system";
import constants from "../constants";
import { deleteAuthJWT, validatePassword } from "../utils/utils";
import authorizedAxios from "../utils/AuthInterceptor";

ProfileMenu.propTypes = {
    user: constants.USER_PROP_TYPE,
    setUser: PropTypes.func
}

ChangePasswordModal.propTypes = {
    setUser: PropTypes.func,
    isModalOpen: PropTypes.bool,
    setIsModalOpen: PropTypes.func
}

const Popup = styled(PopperUnstyled)({
    zIndex: 1000,
});

/*
Button displaying users name, that shows a menu allowing them to logout
or change their password.
Note: I created the same thing for another project, so popup jsx
is basically the same as:
- https://github.com/EECS4481Project/frontend/blob/main/src/agent/dashboard/Dashboard.js
*/
function ProfileMenu(props) {
    const [isChangePasswordOpen, setIsChangePasswordOpen] = useState(false)
    const [anchorEl, setAnchorEl] = useState(null);
    const open = Boolean(anchorEl)

    const closePopup = () => {
        setAnchorEl(null)
    }

    const handleLogout = () => {
        // Post to login endpoint
        authorizedAxios.post('/api/logout').then((res) => {
            // Set user to null & close popup
            props.setUser(null)
            deleteAuthJWT()
            closePopup()
        }).catch((err) => {
            // Error occurred, manually logout
            Cookies.remove("refresh")
            deleteAuthJWT()
            props.setUser(null)
            closePopup()
        })
    }

    return (
        <Box>
            <Button variant="plain" color="neutral" onClick={e => setAnchorEl(e.currentTarget)}>
                <Typography level="body1">{props.user.data.name}</Typography>
                <MoreVertIcon />
            </Button>
            <Popup open={open} anchorEl={anchorEl} disablePortal>
                <ClickAwayListener onClickAway={closePopup}>
                    <MenuList variant="outlined">
                        <MenuItem onClick={() => setIsChangePasswordOpen(true)}>Change Password</MenuItem>
                        <MenuItem onClick={handleLogout}>Logout</MenuItem>
                    </MenuList>
                </ClickAwayListener>
            </Popup>
            <ChangePasswordModal isModalOpen={isChangePasswordOpen} setIsModalOpen={setIsChangePasswordOpen} />
        </Box>
    );
}

/*
Modal that displays a form allowing the user to change their password.
*/
function ChangePasswordModal(props) {
    const [isLoading, setIsLoading] = useState(false)

    const [errorMessage, setErrorMessage] = useState("")

    // Form values
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")


    const changePassword = () => {
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
        authorizedAxios.post('/api/change_password', {
            password
        }).then(res => {
            // Success -- close modal
            props.setIsModalOpen(false)
        }).catch(err => {
            if (err.response.status === 400) {
                setErrorMessage("Invalid Password")
            } else if (err.response.status === 401) {
                // Unauthorized -- user is now signed out
                props.setUser(null)
                props.setIsModalOpen(false)
            }
        }).finally(() => {
            setIsLoading(false)
        })
    }


    return <div>
        <Modal open={props.isModalOpen} onClose={() => props.setIsModalOpen(false)}>
            <ModalDialog
                sx={{ maxWidth: 500, minHeight: 'fit-content', overflow: 'scroll' }}
            >
                <Typography level="h5">
                    Change Password
                </Typography>
                {isLoading && <Box sx={{ display: 'flex', justifyContent: 'center' }}><CircularProgress /></Box>}
                {!isLoading && <div>
                    <Typography level="body2">
                        {constants.PASSWORD_REQUIREMENTS}
                    </Typography>
                    <form
                        onSubmit={(event) => {
                            event.preventDefault();
                            changePassword()
                        }}
                    >
                        <Stack spacing={1}>
                            <FormControl>
                                <FormLabel>New Password</FormLabel>
                                <Input type="password" autoFocus required value={password} onChange={e => setPassword(e.target.value)} />
                            </FormControl>
                            <FormControl>
                                <FormLabel>Confirm New Password</FormLabel>
                                <Input type="password" required value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} />
                            </FormControl>
                            {errorMessage !== "" && <Typography color="danger" level="body3">{errorMessage}</Typography>}
                            <Button type="submit">Change Password</Button>
                        </Stack>
                    </form>
                </div>}
            </ModalDialog>
        </Modal>
    </div>
}

export default ProfileMenu;
