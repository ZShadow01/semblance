const { Client } = require("discord.js");
const path = require("path");
const { walkDir } = require("../utils/helpers");

const EVENTS_DIR = "./src/events";

/**
 * Load all events in the given events directory (and its subdirectories)
 * @param {Client} client - Discord client instance
 */
function setupEvents(client) {
    walkDir(EVENTS_DIR, (filePath, lstat) => {
        if (lstat.isDirectory() || !filePath.endsWith(".js")) return;

        const event = require(path.resolve(filePath));
        if (!("name" in event) || !("execute" in event)) {
            console.log(`'${filePath}' is not an event`);
            return;
        }

        if (event.once) {
            client.once(event.name, (...args) => event.execute(...args));
        } else {
            client.on(event.name, (...args) => event.execute(...args));
        }
    });
}

module.exports = setupEvents;
