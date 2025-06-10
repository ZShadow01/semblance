const { REST, Routes } = require("discord.js");
const path = require("path");
const { walkDir } = require("./src/utils/helpers");

require("dotenv").config();

// Fetch the commands
const commands = [];

walkDir("./src/commands", (filePath, lstat) => {
    if (lstat.isDirectory() || !filePath.endsWith(".js")) return;

    const command = require(path.resolve(filePath));
    if (!("data" in command) || !("execute" in command)) {
        console.log(`Invalid command: '${filePath}'`);
        return;
    }

    // Convert to JSON, ready to be sent to API
    commands.push(command.data.toJSON());
});

// Use REST API to load the commands to Discord
const rest = new REST().setToken(process.env.TOKEN);

(async () => {
    try {
        console.log(`Started refreshing ${commands.length} slash commands`);

        // Upload command information to Discord
        const data = await rest.put(
            Routes.applicationGuildCommands(
                process.env.CLIENT_ID,
                process.env.GUILD_ID
            ),
            { body: commands }
        );

        console.log(`Successfully reloaded ${data.length} slash commands`);
    } catch (err) {
        console.error(err);
    }
})();
