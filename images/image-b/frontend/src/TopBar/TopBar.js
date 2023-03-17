import { Box, Typography } from "@mui/joy";
import PropTypes from 'prop-types';
import constants from "../constants";
import ProfileMenu from "./ProfileMenu";
import Register from "./Register";
import SignIn from "./SignIn";

TopBar.propTypes = {
    user: constants.USER_PROP_TYPE,
    setUser: PropTypes.func
}

function TopBar(props) {
    return (
        <Box height={constants.TOP_BAR_HEIGHT} width="100%" sx={{
            backgroundColor: 'background.level1',
            padding: "0px 10px 0px 10px", display: 'inline-flex', alignItems: 'center'
        }}>
            <Typography level="h4">Chat App</Typography>
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
