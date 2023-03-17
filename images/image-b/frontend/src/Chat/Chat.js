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
Sends a JSON w/ type & content to the websocket
*/
const messageWebSocket = (websocket, type, content) => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
        websocket.send(JSON.stringify({
            type, content
        }))
    }
}

// Returns true if the websocket is open, false otherwise
const isWebSocketOpen = (websocket) => {
    return websocket && websocket.readyState === WebSocket.OPEN
}

/*
Closes previous heartbeat and opens a new one for the websocket.
*/
const initializeHeartbeat = (websocket) => {
    if (heartbeatInterval !== null) {
        clearInterval(heartbeatInterval)
    }
    heartbeatInterval = window.setInterval(function () {
        messageWebSocket(websocket, "ping", "")
    }, 1000);
}

function Chat({ user, setUser }) {
    const [websocket, setWebSocket] = useState(null);
    const [isConnected, setIsConnected] = useState(false);
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [messages, setMessages] = useState([]);

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
        setIsLoggedIn(false)
        setIsConnected(false)
        await wait(500)
        setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
    }

    // Send a message to the websocket
    const sendMessage = (message) => {
        messageWebSocket(websocket, "message", message)
    }

    // Send credentials to the websocket
    const sendAuthentication = () => {
        messageWebSocket(websocket, "auth", getAuthJWT())
    }

    // Handle initial websocket connection
    const handleOnConnect = () => {
        setIsConnected(true)
        // Send credentials if logged in
        if (user) {
            sendAuthentication();
        }
    }

    // Main websocket message handler
    const messageHandler = (msg) => {
        if (msg.type === "refresh") {
            authorizedAxios.get("/api/refresh_credentials").then(() => {
                // Refresh successful, re-auth with socket
                sendAuthentication();
            }).catch(() => {
                // Refresh failed
                setUser(null);
            })
        } else if (msg.type === "signed_in") {
            // Enable messaging
            setIsLoggedIn(true)
        } else if (msg.type === "message") {
            // Received message, update array & sort by timestamp
            setMessages((messages) => {
                return [...messages, {
                    name: msg.name,
                    email: msg.email,
                    ts: msg.ts,
                    message: msg.message
                }].sort((a, b) => a.ts - b.ts)
            })
        }
    }

    // Setup web socket 
    useEffect(() => {
        if (websocket) {
            // If the websocket connects very fast, it will miss the onopen below
            // so we force initial connection handling in this case
            if (isWebSocketOpen(websocket)) {
                handleOnConnect()
            }
            // Handle open
            websocket.onopen = handleOnConnect
            // Initialize heartbeat
            initializeHeartbeat(websocket)
            // handle message recipient
            websocket.onmessage = (event) => {
                const msg = JSON.parse(event.data)
                messageHandler(msg)
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
            <Messages messages={messages} />
            <MessageBar isLoggedIn={isLoggedIn} sendMessage={sendMessage}/>
        </Box>
    );
}

export default Chat;
