import { useEffect, useState } from "react";
import constants from "../constants";
import { Box } from "@mui/system";

let heartbeatInterval = null;

Chat.propTypes = {
    user: constants.USER_PROP_TYPE
}

/*
Returns a ws url with the given path (ie. /xyz)
*/
const generateRelativeWebSocketPath = (path) => {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    return `${protocol}://${window.location.host}${path}`
}

/*
Sends a JSON to the websocket
*/
const sendJSON = (websocket, json) => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
        websocket.send(JSON.stringify(json))
    }
}

/*
Closes previous heartbeat and opens a new one for the websocket.
*/
const initializeHeartbeat = (websocket) => {
    if (heartbeatInterval !== null) {
        clearInterval(heartbeatInterval)
    }
    heartbeatInterval = window.setInterval(function () {
        sendJSON(websocket, {
            type: "ping"
        })
    }, 1000);
}

function Chat() {
    const [websocket, setWebSocket] = useState(null);

    // Setup websocket on page load
    useEffect(() => {
        setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
        // On teardown, close connection & clear heartbeat interval
        return () => {
            clearInterval(heartbeatInterval)
            websocket.close()
        }
    }, [])

    // Setup web socket 
    useEffect(() => {
        if (websocket) {
            console.log(websocket.readyState)
            initializeHeartbeat(websocket)
            websocket.onopen = () => {
                console.log('opened')
            }

            websocket.onmessage = (event) => {
                const msg = event.data
                console.log(msg)
            }

            websocket.onerror = () => {
                // Attempt to reconnect
                setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
            }
            websocket.onclose = () => {
                // Attempt to reconnect
                setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
            }
        }
    }, [websocket])

    return (
        <Box sx={{ flexGrow: 1, flexShrink: 1 }}>
            <p>Hello World!</p>
        </Box>
    );
}

export default Chat;
