import SortedMessageList from "./SortedMessageList.js";


// Message list
const messages = new SortedMessageList()
// Used to ensure duplicate messages aren't added from paging
const messageIds = new Set()

let messagesChanged = false;

// Send update messages to whoever is listing (ie. UI)
const sendUpdatedMessages = () => {
    if (messagesChanged) {
        messagesChanged = false;
        postMessage({
            messages: messages.toArray()
        });
    }
}

// Send updated messages every 200 ms if needed
// This is to avoid blocking the UI with renders if there is a very high
// message throughput
setInterval(sendUpdatedMessages, 200);

onmessage = (e) => {
    // Get event data
    const { newMessages, isNewMessage, instantUpdate } = e.data;
    // Add all non-existent newMessages to res
    for (let i = 0; i < newMessages.length; i++) {
        if (!messageIds.has(newMessages[i].id)) {
            const cleanedMessage = {
                name: newMessages[i].name,
                email: newMessages[i].email,
                ts: newMessages[i].ts,
                message: newMessages[i].message,
                subject: newMessages[i].subject,
                id: newMessages[i].id
            }
            // Sorted-insertion w/ optimizations based on assumption
            // that messages are new or old
            if (isNewMessage) {
                messages.insertMessageAssumingNew(cleanedMessage)
            } else {
                messages.insertMessageAssumingOld(cleanedMessage)
            }
            messageIds.add(newMessages[i].id)
        }
    }
    messagesChanged = true;
    if (instantUpdate) {
        sendUpdatedMessages()
    }
}