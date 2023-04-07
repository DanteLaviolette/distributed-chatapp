import { useEffect, useState } from "react";
import PropTypes from 'prop-types';
import constants from "../constants";
import { Box } from "@mui/system";
import { getAuthJWT } from "../utils/utils";
import authorizedAxios from "../utils/AuthInterceptor";
import Messages from "./Messages";
import MessageBar from "./MessageBar";
import { toast } from "react-toastify";
import axios from "axios";
import { CircularProgress, Typography } from "@mui/joy";

let heartbeatInterval = null;

Chat.propTypes = {
    user: constants.USER_PROP_TYPE,
    setUser: PropTypes.func,
    setUserCount: PropTypes.func
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
const messageWebSocket = (websocket, type, subject, content) => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
        websocket.send(JSON.stringify({
            type, content, subject
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
        messageWebSocket(websocket, "ping", "", "")
    }, 1000);
}

function Chat({ user, setUser, setUserCount }) {
    const [websocket, setWebSocket] = useState(null);
    const [isConnected, setIsConnected] = useState(false);
    const [hasLoadedInitialPage, setHasLoadedInitialPage] = useState(false)
    const [errorOccured, setErrorOccured] = useState(false)
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    // Handle sorting in background
    const [worker, setWorker] = useState(null);

    /**
     * Adds the messages to the messages ensuring no duplicates
     * as well as correct sort order -- using worker
     * @param {Array<String>} newMessages New messages to add
     * @param {boolean} isNewMessage True if these messages are likely new
     * @param {boolean} instantUpdate True if messages should be pushed to UI
     * right after processing. Otherwise, enqueues them to be pushed to UI at
     * some interval.
     */
    function updateMessages(newMessages, isNewMessage, instantUpdate) {
        worker.postMessage({ newMessages, isNewMessage, instantUpdate })
    }
    // Setup websocket & worker on page load
    useEffect(() => {
        // Handle message worker
        const newWorker = new window.Worker('/message-sort-worker.js', { type: "module" })
        setWorker(newWorker);
        // handle websocket
        setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
        // On teardown, close connection & clear heartbeat interval
        return () => {
            clearInterval(heartbeatInterval)
            if (websocket) {
                websocket.close()
            }
            if (worker) {
                worker.terminate()
            }
        }
    }, [])

    // Load initial chat page on load -- ie. when worker is set
    useEffect(() => {
        if (worker) {
            axios.get("/api/messages", { params: { lastTimestamp: Date.now() * constants.MS_TO_NS} }).then(res => {
                if (res.data && res.data.length > 0) {
                    updateMessages(res.data, false, true)
                }
                setHasLoadedInitialPage(true);
            }).catch(() => {
                setErrorOccured(true)
            })
        }
    }, [worker])

    const resetWebSocketValues = () => {
        setIsLoggedIn(false)
        setIsConnected(false)
        setUserCount(null)
    }

    // Attempt to reconnect after waiting 0.5s (to avoid spamming server)
    const attemptReconnect = async () => {
        // Reset values
        resetWebSocketValues()
        websocket.close()
        // Attempt reconnect
        await wait(500)
        setWebSocket(new WebSocket(generateRelativeWebSocketPath("/ws/chat")));
    }

    // Send a message to the websocket
    const sendMessage = (subject, message) => {
        messageWebSocket(websocket, "message", subject, message)
    }

    // Send credentials to the websocket
    const sendAuthentication = () => {
        messageWebSocket(websocket, "auth", "", getAuthJWT())
    }

    // Handle initial websocket connection
    const handleOnConnect = () => {
        setIsConnected(true)
        // Send credentials if logged in
        if (user) {
            sendAuthentication();
        }
        // Initialize heartbeat
        initializeHeartbeat(websocket)
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
        } else if (msg.type === "message_failed") {
            // Notify of failed message
            toast.error("Failed to send message", constants.TOAST_CONFIG)
        } else if (msg.type === "message") {
            // Add message to state -- update instantly if from curr user
            updateMessages([msg], true, user && user.data.email === msg.email)
        } else if (msg.type === "user_count") {
            setUserCount({
                anonymousUsers: msg.anonymousUsers,
                authorizedUsers: msg.authorizedUsers
            })
        }
    }

    // Setup web socket 
    useEffect(() => {
        if (websocket) {
            // If the websocket connects very fast, it will miss the onopen below
            // so we force initial connection handling in this case
            if (isWebSocketOpen(websocket) && !isConnected) {
                handleOnConnect()
            }
            // Handle open
            websocket.onopen = handleOnConnect
            // handle message recipient
            websocket.onmessage = (event) => {
                const msg = JSON.parse(event.data)
                messageHandler(msg)
            }
            // Handle cleanup on error
            websocket.onerror = resetWebSocketValues
            // Handle retry connection on close
            websocket.onclose = attemptReconnect
        }
    }, [websocket, isConnected])

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

    // Display chat screen, error or loading message depending on state
    let chatScreen = <>
        <Messages worker={worker} updateMessages={updateMessages} />
        <MessageBar isLoggedIn={isLoggedIn} sendMessage={sendMessage}/>
        </>
    if (errorOccured) {
        chatScreen = <Box sx={{width: "100%", height: "100%", display: "flex", alignItems: "center", justifyContent: "center"}}>
            <Typography>Something went wrong. Please refresh the page.</Typography>
        </Box>
    } else if (!hasLoadedInitialPage) {
        chatScreen = <Box sx={{width: "100%", height: "100%", display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center"}}>
            <CircularProgress size="md"/>
            <Typography>Loading...</Typography>
    </Box>
    }

    return (
        <Box sx={{ flexGrow: 1, flexShrink: 1, maxHeight: `calc(100vh - ${constants.TOP_BAR_HEIGHT}px - ${constants.MESSAGE_BAR_HEIGHT})` }}>
            {chatScreen}
        </Box>
    );
}

export default Chat;
