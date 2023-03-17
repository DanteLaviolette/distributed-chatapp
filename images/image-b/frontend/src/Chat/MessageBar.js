import { Box, Button, Divider, Input } from "@mui/joy";
import SendIcon from '@mui/icons-material/Send';
import PropTypes from 'prop-types';
import { Fragment, useRef, useState } from "react";

MessageBar.propTypes = {
    sendMessage: PropTypes.func,
    isLoggedIn: PropTypes.bool,
}

// Message bar that displays an input and send button
// Input/button will be disabled if not logged in
function MessageBar({ sendMessage, isLoggedIn }) {
    const [subject, setSubject] = useState("")
    const [message, setMessage] = useState("")
    const messageInput = useRef()

    const handleMessageSend = (subject, message) => {
        // Send message & clear input
        sendMessage(subject, message)
        setSubject("")
        setMessage("")
    }

    const subjectInput = <Fragment><Input
        variant="plain"
        value={subject}
        onChange={e => setSubject(e.target.value)}
        placeholder={"Subject"}
        readOnly={!isLoggedIn}
        inputProps={{ maxLength: 12 }}
        onFocus={(e) => {
            // Enable focus styling of outer input on focus
            messageInput.current.classList.add("Joy-focused")
        }}
        onBlur={(e) => {
            // Disable focus styling of outer input on focus
            messageInput.current.classList.remove("Joy-focused")
        }}
        onKeyDown={(ev) => {
            // Handle sending message on enter
            if (ev.key === 'Enter') {
                ev.preventDefault();
                // Move focus to the next input
                messageInput.current.children[1].focus()
            }
        }}
        sx={{ width: "120px", '&:hover': { bgcolor: 'transparent' }, bgcolor: 'transparent', "--Input-focusedThickness": "0px" }} />
        <Divider orientation="vertical" />
    </Fragment>

    return (
        <Box height="50px" width="100%" sx={{ display: "inline-flex", padding: "5px 10px 5px 10px", marginBottom: "5px" }}>
            <Input value={message}
                ref={messageInput}
                placeholder={isLoggedIn ? "Input your message here!" : "Sign in to send a message"}
                readOnly={!isLoggedIn}
                autoFocus={true}
                onChange={e => setMessage(e.target.value)}
                onKeyDown={(ev) => {
                    // Handle sending message on enter
                    if (ev.key === 'Enter') {
                        ev.preventDefault();
                        handleMessageSend(subject, message)
                    }
                }}
                startDecorator={subjectInput}
                sx={{ flexGrow: 1, flexShrink: 1, marginRight: "10px", marginLeft: "10px", paddingLeft: "0px" }} />
            <Button size="sm" color="info"
                onClick={() => handleMessageSend(subject, message)}
                disabled={message === "" || !isLoggedIn}
            ><SendIcon /></Button>
        </Box>
    );
}

export default MessageBar;
