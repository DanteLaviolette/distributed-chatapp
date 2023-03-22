
class MessageNode {
    constructor(val, prev, next) {
        this.val = val
        this.prev = prev
        this.next = next
    }
}

/**
 * LinkedList that supports sorted insertion by the ts value of the inserted
 * object.
 * - When old messages are loaded, they'll likely have O(1) insertion into
 * the head (head is oldest element)
 * - When new messages are loaded, they'll likely have O(1) insertion into
 * the tail (tail is newest element)
 * 
 * Sorted insertion is needed for race conditions (ie. backend sends an old
 * message faster than a new message)
 * - In these cases insertion will be O(n), although this will be closer
 * to O(1) in practice.
 */
export default class SortedMessageList {
    constructor() {
        this.head = null // oldest message
        this.tail = null // newest message
        this.length = 0
        // Expose functions
        this.insertMessageAssumingOld = this.sortedInsertionFromHead
        this.insertMessageAssumingNew = this.sortedInsertionFromTail
        this.getOldestMessage = () => this.head.val
    }

    insertValueBeforeNode(val, listNode) {
        let originalPreviousNode = listNode.prev
        // Insert before node
        listNode.prev = new MessageNode(val, listNode.prev, listNode)
        // Update original previous nodes next
        if (originalPreviousNode) {
            originalPreviousNode.next = listNode.prev
        }
        // Handle head case
        if (listNode === this.head) {
            this.head = listNode.prev
        }
    }

    insertValueAfterNode(val, listNode) {
        let originalNextNode = listNode.next
        // Insert after node
        listNode.next = new MessageNode(val, listNode, listNode.next)
        // Update original next nodes prev
        if (originalNextNode) {
            originalNextNode.prev = listNode.next
        }
        // Handle tail case
        if (listNode === this.tail) {
            this.tail = listNode.next
        }
    }

    /**
     * Inserts the val in sorted order (by val.ts). Begins iteration at the
     * head for faster inserts of old messages.
     * @param {message} val 
     * @returns None
     */
    sortedInsertionFromHead(val) {
        this.length += 1
        // First element edge case
        if (this.head == null) {
            this.head = new MessageNode(val, null, null)
            this.tail = this.head
            return
        }
        // Normal case -- iterate linked list
        let curr = this.head
        while (curr != null) {
            // Found insertion spot
            if (val.ts < curr.val.ts) {
                this.insertValueBeforeNode(val, curr)
                return
            }
            curr = curr.next
        }
        // Must be newest message
        this.insertValueAfterNode(val, this.tail)
    }

    /**
     * Inserts the val in sorted order (by val.ts). Begins iteration at the
     * tail for faster inserts of new messages.
     * @param {message} val 
     * @returns None
     */
    sortedInsertionFromTail(val) {
        this.length += 1
        // First element edge case
        if (this.tail == null) {
            this.tail = new MessageNode(val, null, null)
            this.head = this.tail
            return
        }
        // Normal case -- iterate linked from tail
        let curr = this.tail
        while (curr != null) {
            // Found insertion spot
            if (val.ts > curr.val.ts) {
                this.insertValueAfterNode(val, curr)
                return
            }
            curr = curr.prev
        }
        // Must be new oldest element
        this.insertValueBeforeNode(val, this.head)
    }

    /**
     * 
     * @returns Array representing the linkedList
     */
    toArray() {
        let res = []
        let curr = this.head
        while (curr != null) {
            res.push(curr.val)
            curr = curr.next
        }
        return res
    }
}