module.exports = {
  preset: "ts-jest",
  moduleFileExtensions: ["js", "json", "ts", "tsx"],
  testMatch: ["**/?(*.)+(spec|test).[tj]s?(x)"],
  testEnvironment: "node",
  roots: ["<rootDir>"],
  moduleNameMapper: {
    "^@domains/(.*)$": "<rootDir>/src/domains/$1",
    "^@frameworks/(.*)$": "<rootDir>/src/frameworks/$1"
  }
}
