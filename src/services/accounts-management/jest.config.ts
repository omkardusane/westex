export default {
    preset: 'ts-jest',
    testEnvironment: 'node',
    testPathIgnorePatterns: ['./dist/'],
    collectCoverage: false,
    moduleFileExtensions: ['ts', 'js', 'json', 'node'],
    roots: ['./src', './tests'],
    setupFilesAfterEnv: ['./jest.setup.ts'], // Optional for custom setups
};
