import { useEffect, useState } from "react";
import constants from "../constants";
import { Box } from "@mui/system";

Chat.propTypes = {
    user: constants.USER_PROP_TYPE
}

const generateRelativeWebSocketPath = (path) => {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    return `${protocol}://${window.location.host}${path}`
}

function Chat() {

    useEffect(() => {
        let websocket = new WebSocket(generateRelativeWebSocketPath("/ws/chat"));
        websocket.onopen = () => {
            console.log("open!")
            websocket.send("hey")
        }
    })

  return (
    <Box sx={{flexGrow: 1, flexShrink: 1}}>
      <p>Hello World!</p>
    </Box>
  );
}

export default Chat;
