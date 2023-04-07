import { Box, Typography } from "@mui/joy";
import axios from "axios";
import PropTypes from 'prop-types';
import { memo, useEffect, useState } from "react";
import { toast } from "react-toastify";
import InfiniteScroll from 'react-infinite-scroll-component';
import constants from "../constants";

const messagePropType = PropTypes.shape({
    name: PropTypes.string,
    email: PropTypes.string,
    message: PropTypes.string,
    subject: PropTypes.string,
    ts: PropTypes.number
})

Messages.propTypes = {
    worker: PropTypes.instanceOf(window.Worker),
    updateMessages: PropTypes.func
}

Message.propTypes = {
    message: messagePropType
}

const MemoMessage = memo(Message)

// Displays all of the given messages using paging
function Messages({ updateMessages, worker }) {
    const [loadedAllMessages, setLoadedAllMessages] = useState(false)
    const [messages, setMessages] = useState([])

    useEffect(() => {
        if (worker) {
            worker.onmessage = (e) => {
                setMessages(e.data.messages)
            };
        }
    });

    const messageBoxId = "messageBox"

    // Load the next page of messages
    function loadMoreMessages() {
        // Get oldest message timestamp
        if (messages.length === 0) {
            return
        }
        const oldestMessage = messages[0]
        const lastTimestamp = oldestMessage.ts
        // Get older messages
        axios.get("/api/messages", { params: { lastTimestamp } }).then(res => {
            // Handle case where all messages have been received
            if (res.data === null || res.data.length === 0) {
                setLoadedAllMessages(true)
            } else {
                // Store messages in chat
                updateMessages(res.data)
            }
        }).catch(() => {
            toast.error("Failed to load previous messages. Try again later.", constants.TOAST_CONFIG)
        });
    }

    return (
        <Box id={messageBoxId} height="100%" width="100%" sx={{ overflow: "auto", display: 'flex', flexDirection: 'column-reverse' }}>
            <InfiniteScroll
                dataLength={messages.length}
                next={loadMoreMessages}
                scrollableTarget={messageBoxId}
                inverse={true}
                hasMore={!loadedAllMessages}
                loader={messages.length === 0 &&
                <Box sx={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                    <Typography level="body3">
                        No messages. Start the conversation 👋
                    </Typography>
                </Box>}
            >
                {messages.map(val => <MemoMessage key={val.id} message={val.message} subject={val.subject} name={val.name} email={val.email} timestamp={val.ts} />)}
            </InfiniteScroll>
        </Box>
    );
}

// Component for a single message
function Message({ message, subject, name, email, timestamp }) {
    const ts = Math.floor(timestamp / constants.MS_TO_NS)
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
            <Typography fontWeight="bold" sx={{ color: 'neutral', width: 'fit-content' }}>{name + " | " + email}</Typography>
            <Box sx={{ flexGrow: 1, flexShrink: 1 }}></Box>
            <Typography level="body4" sx={{ color: 'neutral', width: 'fit-content' }}>{date}</Typography>
        </Box>
        {subject !== "" && <Typography sx={{ color: 'neutral' }}>{subject}</Typography>}
        <Typography sx={{ color: 'neutral.300' }}>{message}</Typography>
    </Box>
}

export default Messages;
