import { Box, Tooltip, Typography } from "@mui/joy";
import axios from "axios";
import PropTypes from 'prop-types';
import { useEffect, useRef, useState } from "react";
import { toast } from "react-toastify";
import useMeasure from 'react-use-measure'
import constants from "../constants";

const messagePropType = PropTypes.shape({
    name: PropTypes.string,
    email: PropTypes.string,
    message: PropTypes.string,
    subject: PropTypes.string,
    ts: PropTypes.number
})

Messages.propTypes = {
    messages: PropTypes.arrayOf(messagePropType),
    updateMessages: PropTypes.func
}

Message.propTypes = {
    message: messagePropType
}

// Displays all of the given messages
function Messages({ messages, updateMessages }) {
    const [isLoadingMessages, setIsLoadingMessages] = useState(false)
    const [loadedAllMessages, setLoadedAllMessages] = useState(false)

    const messageBox = useRef()
    const [messageBoxChild, bounds] = useMeasure()
    const [previousMessageBoxChildHeight, setPreviousMessageBoxChildHeight] = useState(null)

    // Persist scroll position on messages added
    useEffect(() => {
        if (messageBox.current && bounds && previousMessageBoxChildHeight) {
            // Move scroll to previous position
            messageBox.current.scrollTop += (bounds.height - previousMessageBoxChildHeight)
            setIsLoadingMessages(false);
        }
        // Keep track of previous height
        if (bounds) {
            setPreviousMessageBoxChildHeight(bounds.height)
        }
    }, [bounds, previousMessageBoxChildHeight])

    function handleScroll(ev) {
        // When the user scrolls to the top load more messages if some
        // are available and we aren't already loading them
        if (ev.target.children[0].getBoundingClientRect().y === constants.TOP_BAR_HEIGHT
                && !isLoadingMessages && !loadedAllMessages) {
            // Set is loading
            setIsLoadingMessages(true)
            // Get messages
            const lastTimestamp = messages[0].ts
            axios.get("/api/messages", { params: { lastTimestamp } }).then(res => {
                // Handle case where all messages have been received
                if (res.data === null || res.data.length === 0) {
                    setLoadedAllMessages(true)
                    setIsLoadingMessages(false)
                } else {
                    // Store messages in chat
                    updateMessages(res.data)
                }
            }).catch(() => {
                toast.error("Failed to load previous messages. Try again later.", constants.TOAST_CONFIG)
                setIsLoadingMessages(false);
            });
        }
    }

    return (
        <Box ref={messageBox} onScroll={handleScroll} height="100%" width="100%" sx={{ overflow: "auto", display: 'flex', flexDirection: 'column-reverse' }}>
            <Box ref={messageBoxChild}>
                {(loadedAllMessages || isLoadingMessages) && <Box sx={{ display: "flex", alignItems: "center", justifyContent: "center" }}>
                    <Typography level="body3">{isLoadingMessages ? "Loading Previous Messages..." : "Loaded All Messages"}</Typography>
                    </Box>}
                {messages.map((message, i) => {
                    return <Message key={i} message={message} />
                })}
            </Box>
        </Box>
    );
}

// Component for a single message
function Message({ message }) {
    const ts = Math.floor(message.ts / constants.MS_TO_NS)
    const date = new Date(ts).toLocaleDateString('en-us', {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "numeric",
        minute: "numeric",
        second: "numeric"
    })
    return <Box sx={{ width: "100%", padding: "5px 10px 0px 10px", margin: "5px 0px 5px 0px" }}>
        <Box sx={{ display: "flex" }}>
            <Tooltip title={message.email} sx={{ width: 'fit-content' }}>
                <Typography fontWeight="bold" sx={{ color: 'neutral', width: 'fit-content' }}>{message.name}</Typography>
            </Tooltip>
            <Box sx={{ flexGrow: 1, flexShrink: 1 }}></Box>
            <Typography level="body4" sx={{ color: 'neutral', width: 'fit-content' }}>{date}</Typography>
        </Box>
        {message.subject !== "" && <Typography sx={{ color: 'neutral' }}>{message.subject}</Typography>}
        <Typography sx={{ color: 'neutral.300' }}>{message.message}</Typography>
    </Box>
}

export default Messages;
