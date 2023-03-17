import { useEffect, useState } from "react";
import PropTypes from 'prop-types';
import constants from "../constants";
import { Box } from "@mui/system";
import { getAuthJWT } from "../utils/utils";
import authorizedAxios from "../utils/AuthInterceptor";
import Messages from "./Messages";
import MessageBar from "./MessageBar";

let heartbeatInterval = null;

Chat.propTypes = {
    user: constants.USER_PROP_TYPE,
    setUser: PropTypes.func
}

const wait = (ms) => new Promise((res) => setTimeout(res, ms));

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

function Chat({ user, setUser }) {
    const [websocket, setWebSocket] = useState(null);
    const [isConnected, setIsConnected] = useState(false)

    const fakeMessages = [
        {
            name: "dante",
            email: "abc@gmail.com",
            message: "1",
            ts: 1679025580872
        },
        {
            name: "dante2",
            email: "abc@gmail.com",
            message: "2",
            ts: 1679025580873
        },
        {
            name: "dante3",
            email: "abc@gmail.com",
            message: "3",
            ts: 1679025580875
        },
        {
            name: "dante",
            email: "abc@gmail.com",
            message: "4",
            ts: 1679025580872
        },
        {
            name: "dante2",
            email: "abc@gmail.com",
            message: "5",
            ts: 1679025580873
        },
        {
            name: "dante3",
            email: "abc@gmail.com",
            message: "6",
            ts: 1679025580875
        },
        {
            name: "dante",
            email: "abc@gmail.com",
            message: "7",
            ts: 1679025580872
        },
        {
            name: "dante2",
            email: "abc@gmail.com",
            message: "8",
            ts: 1679025580873
        },
        {
            name: "dante3",
            email: "abc@gmail.com",
            message: "9",
            ts: 1679025580875
        },
    ]

    // Setup websocket on page load
    useEffect(() => {
        setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
        // On teardown, close connection & clear heartbeat interval
        return () => {
            clearInterval(heartbeatInterval)
            if (websocket) {
                websocket.close()
            }
        }
    }, [])

    // Attempt to reconnect after waiting 0.5s (to avoid spamming server)
    const attemptReconnect = async () => {
        setIsConnected(false)
        await wait(500)
        setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
    }

    const sendAuthentication = () => {
        sendJSON(websocket, {
            type: 'auth',
            content: getAuthJWT()
        })
    }

    const messageHandler = (type, content) => {
        if (type === "refresh") {
            authorizedAxios.get("/api/refresh_credentials").then(() => {
                // Refresh successful, re-auth with socket
                sendAuthentication();
            }).catch(() => {
                // Refresh failed
                setUser(null);
            })
        } else if (type === "signed_in") {
            // TODO: Enable messaging
        }
    }

    // Setup web socket 
    useEffect(() => {
        if (websocket) {
            // Handle open
            websocket.onopen = () => {
                setIsConnected(true)
                // Send credentials if logged in
                if (user) {
                    sendAuthentication();
                }
            }

            // Initialize heartbeat
            initializeHeartbeat(websocket)

            // handle message
            websocket.onmessage = (event) => {
                const msg = JSON.parse(event.data)
                messageHandler(msg.type, msg.content)
            }
            // Handle retry connection on close
            websocket.onerror = attemptReconnect
            websocket.onclose = attemptReconnect
        }
    }, [websocket])

    // Send updated credentials on user change or restart session on logout
    useEffect(() => {
        if (user) {
            sendAuthentication();
        } else {
            if (websocket && websocket && websocket.readyState === WebSocket.OPEN) {
                websocket.close()
            }
        }
    }, [user])

    return (
        <Box sx={{ flexGrow: 1, flexShrink: 1, maxHeight: `calc(100vh - ${constants.TOP_BAR_HEIGHT} - ${constants.MESSAGE_BAR_HEIGHT})` }}>
            {/*<p>Connected: {isConnected ? "true" : "false"}</p>*/}
            <Messages messages={fakeMessages} />
            <MessageBar isLoggedIn={!!user}/>
        </Box>
    );
}

export default Chat;
