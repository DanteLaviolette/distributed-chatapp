const { SortedMessageList } = require("./SortedMessageList");

test('SortedMessageList correct sorted order (starting with tail)', () => {
    const expected = [{ ts: -20 }, { ts: -15 }, { ts: -10 },
        { ts: -5 }, { ts: -2 }, { ts: 0 }, { ts: 2 }, { ts: 2 }, { ts: 6 },
        { ts: 10 }, { ts: 12 }, { ts: 15 }, { ts: 18 }, { ts: 50 }
    ]

    let list = new SortedMessageList()
    list.insertMessageAssumingNew({ ts: 0 }) // New tail
    list.insertMessageAssumingNew({ ts: 6 }) // New tail
    list.insertMessageAssumingNew({ ts: 10 }) // New tail
    list.insertMessageAssumingNew({ ts: 15 }) // New tail
    list.insertMessageAssumingNew({ ts: 2 }) // Middle insertion
    list.insertMessageAssumingNew({ ts: 12 }) // Middle insertion
    list.insertMessageAssumingNew({ ts: 18 }) // New tail
    list.insertMessageAssumingNew({ ts: 2 }) // Duplicate
    list.insertMessageAssumingNew({ ts: -5 }) // New head -- starting from tail
    list.insertMessageAssumingNew({ ts: -10 }) // New head -- starting from tail
    list.insertMessageAssumingOld({ ts: -15 }) // New head -- starting from head
    list.insertMessageAssumingOld({ ts: -20 }) // New head
    list.insertMessageAssumingOld({ ts: -2 }) // Middle insertion from head
    list.insertMessageAssumingOld({ ts: 50 }) // New tail -- starting from head
    expect(list.toArray()).toEqual(expected)
});

test('SortedMessageList correct sorted order (starting with head)', () => {
    const expected = [{ ts: -25 }, { ts: -20 }, { ts: -17 }, { ts: -15 },
        { ts: 50 }, { ts: 52 }, { ts: 55 }
    ]

    let list = new SortedMessageList()
    list.insertMessageAssumingOld({ ts: -15 }) // New head
    list.insertMessageAssumingOld({ ts: -20 }) // New head
    list.insertMessageAssumingOld({ ts: -17 }) // Middle insertion
    list.insertMessageAssumingOld({ ts: 50 }) // New tail -- starting from head
    list.insertMessageAssumingNew({ ts: 55 }) // New tail
    list.insertMessageAssumingNew({ ts: 52 }) // Middle insertion
    list.insertMessageAssumingNew({ ts: -25 }) // New tail
    expect(list.toArray()).toEqual(expected)
});