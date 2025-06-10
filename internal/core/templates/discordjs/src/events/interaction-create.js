const { Events, MessageFlags } = require("discord.js");

module.exports = {
    name: Events.InteractionCreate,

    async execute(interaction) {
        // Ignore non-command interactions
        if (!interaction.isChatInputCommand()) return;

        const command = interaction.client.commands.get(
            interaction.commandName
        );

        // Errors if command doesn't exist
        if (!command) {
            console.error(`No command matches ${interaction.commandName}`);
            return;
        }

        try {
            // Execute the command
            await command.execute(interaction);
        } catch (err) {
            // Respond to user even if the command fails
            console.error(err);

            // Use 'followUp()' if the interaction already got a reply
            const reply =
                interaction.replied || interaction.deferred
                    ? interaction.followUp
                    : interaction.reply;

            await reply({
                content: `There was an error while executing ${interaction.commandName}`,
                flags: MessageFlags.Ephemeral,
            });
        }
    },
};
