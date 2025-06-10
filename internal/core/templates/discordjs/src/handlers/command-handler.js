const { Collection, Client } = require("discord.js");
const path = require("path");
const { walkDir } = require("../utils/helpers");

// Directory that holds all command files
const COMMANDS_DIR = "./src/commands";

/**
 * Load all commands in the given commands directory (and its subdirectories)
 * @param {Client} client - Discord client instance
 */
function setupCommands(client) {
    client.commands = new Collection();

    walkDir(COMMANDS_DIR, (filePath, lstat) => {
        if (lstat.isDirectory() || !filePath.endsWith(".js")) return;

        const command = require(path.resolve(filePath));
        if (!("data" in command) || !("execute" in command)) {
            console.log(`'${filePath}' is not a command`);
            return;
        }

        client.commands.set(command.data.name, command);
    });
}

module.exports = setupCommands;
