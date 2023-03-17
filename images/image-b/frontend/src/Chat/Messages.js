import { Box, Tooltip, Typography } from "@mui/joy";
import PropTypes from 'prop-types';

const messagePropType = PropTypes.shape({
    name: PropTypes.string,
    email: PropTypes.string,
    message: PropTypes.string,
    ts: PropTypes.number
})

Messages.propTypes = {
    messages: PropTypes.arrayOf(messagePropType)
}

Message.propTypes = {
    message: messagePropType
}

// Displays all of the given messages
function Messages({ messages }) {
    return (
        <Box height="100%" width="100%" sx={{ overflow: "auto", display: 'flex', flexDirection: 'column-reverse' }}>
            <Box>
            {messages.map((message, i) => {
                return <Message key={i} message={message} />
            })}
            </Box>
        </Box>
    );
}

// Component for a single message
function Message({ message }) {
    const date = new Date(message.ts).toLocaleDateString('en-us', {
        year:"numeric",
        month:"short",
        day:"numeric",
        hour:"numeric",
        minute: "numeric",
        second: "numeric"
    }) 
    return <Box sx={{ width: "100%", padding: "5px 10px 0px 10px", margin: "5px 0px 5px 0px" }}>
        <Box sx={{ display: "flex" }}>
            <Tooltip title={message.email} sx={{ width: 'fit-content' }}>
                <Typography sx={{ color: 'neutral', width: 'fit-content' }}>{message.name}</Typography>
            </Tooltip>
            <Box sx={{ flexGrow: 1, flexShrink: 1 }}></Box>
            <Typography level="body4" sx={{ color: 'neutral', width: 'fit-content' }}>{date}</Typography>
        </Box>
        <Typography sx={{ color: 'neutral.300' }}>{message.message}</Typography>
    </Box>
}

export default Messages;
