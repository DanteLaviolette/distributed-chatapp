import { Box, Button, Input } from "@mui/joy";
import SendIcon from '@mui/icons-material/Send';
import PropTypes from 'prop-types';
import { useState } from "react";

MessageBar.propTypes = {
    sendMessage: PropTypes.func,
    isLoggedIn: PropTypes.bool,
}

// Message bar that displays an input and send button
// Input/button will be disabled if not logged in
function MessageBar({ sendMessage, isLoggedIn }) {
    const [message, setMessage] = useState("")

    const handleMessageSend = (message) => {
        // Send message & clear input
        sendMessage(message)
        setMessage("")
    }

    return (
        <Box height="50px" width="100%" sx={{display: "inline-flex", padding: "5px 10px 5px 10px", marginBottom: "5px"}}>
            <Input value={message}
                placeholder={isLoggedIn ? "Input your message here!" : "Sign in to send a message"}
                readOnly={!isLoggedIn}
                onChange={e => setMessage(e.target.value)}
                onKeyDown={(ev) => {
                    // Handle sending message on enter
                    if (ev.key === 'Enter') {
                      ev.preventDefault();
                      handleMessageSend(message)
                    }
                  }}
                sx={{flexGrow: 1, flexShrink: 1, marginRight: "10px"}}>
            </Input>
            <Button size="sm" color="info"
                onClick={() => handleMessageSend(message)}
                disabled={message === "" || !isLoggedIn}
            ><SendIcon/></Button>
        </Box>
    );
}

export default MessageBar;
