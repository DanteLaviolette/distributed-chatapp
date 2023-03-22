const { SortedMessageList } = require("./SortedMessageList");


test('SortedMessageList correct sorted order (Only head inserts)', () => {
    const expected = [{ ts: -5 }, { ts: 0 }, { ts: 2 }, { ts: 5 }]
    // All possible state changes on head inserts
    let list = new SortedMessageList()
    list.insertMessageAssumingOld({ ts: 0 }) // New head/tail
    list.insertMessageAssumingOld({ ts: -5 }) // New head
    list.insertMessageAssumingOld({ ts: 5 }) // New tail
    list.insertMessageAssumingOld({ ts: 2 }) // Middle insert
    expect(list.toArray()).toEqual(expected)
});

test('SortedMessageList correct sorted order (Only tail inserts)', () => {
    const expected = [{ ts: -5 }, { ts: 0 }, { ts: 2 }, { ts: 5 }]
    // All possible state changes on tail inserts
    let list = new SortedMessageList()
    list.insertMessageAssumingNew({ ts: 0 }) // New head/tail
    list.insertMessageAssumingNew({ ts: -5 }) // New head
    list.insertMessageAssumingNew({ ts: 5 }) // New tail
    list.insertMessageAssumingNew({ ts: 2 }) // Middle insert
    expect(list.toArray()).toEqual(expected)
});


test('SortedMessageList correct sorted order (both insert types)', () => {
    const expected = [{ ts: -10 }, { ts: -5 }, { ts: -2 }, { ts: 0 },
        { ts: 2 }, { ts: 2 }, { ts: 15 }, { ts: 15 }, { ts: 18 }, { ts: 50 }
    ]
    // Tests all possible state changes on both head & tail
    let list = new SortedMessageList()
    list.insertMessageAssumingNew({ ts: 0 }) // New head/tail
    list.insertMessageAssumingNew({ ts: 15 }) // New tail
    list.insertMessageAssumingNew({ ts: 2 }) // Middle insertion
    list.insertMessageAssumingNew({ ts: 18 }) // New tail
    list.insertMessageAssumingNew({ ts: 2 }) // Duplicate
    list.insertMessageAssumingNew({ ts: -5 }) // New head -- from tail
    list.insertMessageAssumingOld({ ts: -2 }) // Middle insertion from head
    list.insertMessageAssumingOld({ ts: 15 }) // Duplicate
    list.insertMessageAssumingOld({ ts: 50 }) // New tail -- from head
    list.insertMessageAssumingOld({ ts: -10 }) // New head -- from head
    expect(list.toArray()).toEqual(expected)
});