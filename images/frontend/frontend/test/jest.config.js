const config = {
    testEnvironment: 'jest-environment-node',
    transform: {
        "^.+\\.[t|j]sx?$": "babel-jest"
    }
};

module.exports = config;