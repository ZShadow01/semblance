const { Client, GatewayIntentBits } = require("discord.js");
const setupCommands = require("./handlers/command-handler");
const setupEvents = require("./handlers/event-handler");

require("dotenv").config();

// Create Discord client instance
const client = new Client({ intents: [GatewayIntentBits.Guilds] });

// Setup handlers
setupCommands(client, "./commands");
setupEvents(client, "./events");

// Start the client
client.login(process.env.TOKEN);
