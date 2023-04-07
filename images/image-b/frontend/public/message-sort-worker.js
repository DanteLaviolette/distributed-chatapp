import SortedMessageList from "./SortedMessageList.js";


// Message list
const messages = new SortedMessageList()
// Used to ensure duplicate messages aren't added from paging
const messageIds = new Set()

onmessage = (e) => {
    // Get event data
    const { newMessages, isNewMessage } = e.data;
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
    postMessage({
        messages: messages.toArray()
    });
}