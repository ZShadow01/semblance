const { Events } = require("discord.js");

module.exports = {
    name: Events.ClientReady, // or "ready"
    once: true,

    async execute(client) {
        console.log(client.user.username);
        console.log(client.user.id);
        console.log("Logged in");
    },
};
