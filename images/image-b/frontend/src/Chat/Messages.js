import { Box, Tooltip, Typography } from "@mui/joy";
import axios from "axios";
import PropTypes from 'prop-types';
import { useState } from "react";
import { toast } from "react-toastify";
import InfiniteScroll from 'react-infinite-scroll-component';
import constants from "../constants";
import SortedMessageList from "../utils/SortedMessageList";

const messagePropType = PropTypes.shape({
    name: PropTypes.string,
    email: PropTypes.string,
    message: PropTypes.string,
    subject: PropTypes.string,
    ts: PropTypes.number
})

Messages.propTypes = {
    messages: PropTypes.instanceOf(SortedMessageList),
    updateMessages: PropTypes.func
}

Message.propTypes = {
    message: messagePropType
}

// Displays all of the given messages using paging
function Messages({ messages, updateMessages }) {
    const [loadedAllMessages, setLoadedAllMessages] = useState(false)

    const messageBoxId = "messageBox"

    // Load the next page of messages
    function loadMoreMessages() {
        // Get oldest message timestamp
        const oldestMessage = messages.getOldestMessage()
        if (oldestMessage == null) {
            return
        }
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

    // Convert messages linked list to array of elements
    const messagesToElements = () => {
        let arr = []
        let curr = messages.head
        while (curr != null) {
            arr.push( <Message key={curr.val.id} message={curr.val} />)
            curr = curr.next
        }
        return arr
    }

    return (
        <Box id={messageBoxId} height="100%" width="100%" sx={{ overflow: "auto", display: 'flex', flexDirection: 'column-reverse' }}>
            <InfiniteScroll
                dataLength={messages.length}
                next={loadMoreMessages}
                scrollableTarget={messageBoxId}
                inverse={true}
                hasMore={!loadedAllMessages}
                loader={<Typography level="body3">Loading messages...</Typography>}
            >
                {messagesToElements()}
            </InfiniteScroll>
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
