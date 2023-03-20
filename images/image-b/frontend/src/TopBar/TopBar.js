import { Box, CircularProgress, Typography } from "@mui/joy";
import PropTypes from 'prop-types';
import constants from "../constants";
import ProfileMenu from "./ProfileMenu";
import Register from "./Register";
import SignIn from "./SignIn";

TopBar.propTypes = {
    user: constants.USER_PROP_TYPE,
    setUser: PropTypes.func,
    userCount: PropTypes.shape({
        anonymousUsers: PropTypes.number,
        authorizedUsers: PropTypes.number
    })
}

// Top bar that shows the application name, user count, along with a profile menu
// if signed in, or a sign in & register button otherwise.
function TopBar(props) {
    return (
        <Box height={constants.TOP_BAR_HEIGHT + "px"} width="100%" sx={{
            backgroundColor: 'background.level1',
            padding: "0px 10px 0px 10px", display: 'inline-flex', alignItems: 'center'
        }}>
            <Typography level="h4">Chat App</Typography>
            <Box sx={{display: 'flex', flexDirection: 'column', marginLeft: '10px'}}>
                <Typography level="body4">
                    Online Anonymous Users: {!props.userCount ? "" : props.userCount.anonymousUsers}
                </Typography>
                <Typography level="body4">
                    Online Registered Users: {!props.userCount ? "" : props.userCount.authorizedUsers}
                </Typography>
            </Box>
            {!props.userCount && <CircularProgress sx={{marginLeft: "10px"}} size="sm"/>}
            <Box flexGrow={10} flexShrink={10} />
            {!props.user && <Box sx={{ display: 'inline-flex' }}>
                <SignIn setUser={props.setUser} />
                <Register setUser={props.setUser} />
            </Box>}
            {props.user && <ProfileMenu user={props.user} setUser={props.setUser} />}
        </Box>
    );
}

export default TopBar;
