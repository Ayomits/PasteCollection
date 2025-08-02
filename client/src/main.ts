import "reflect-metadata";

import { dirname, importx } from "@discordx/importer";
import type { Interaction, Message } from "discord.js";
import { GatewayIntentBits } from "discord.js";
import { Client, DIService, tsyringeDependencyRegistryEngine } from "discordx";
import { container } from "tsyringe";

import { configService } from "#/config/config.js";
import { createLogger } from "#logger/index.js";

async function bootstrap() {
  DIService.engine = tsyringeDependencyRegistryEngine.setInjector(container);
  const APP_ENV = configService.get("APP_ENV");
  const logger = createLogger();
  const client = new Client({
    intents: [
      GatewayIntentBits.MessageContent,
      GatewayIntentBits.GuildMembers,
      GatewayIntentBits.Guilds,
      GatewayIntentBits.GuildMessages,
    ],
    logger: logger,
    simpleCommand: {
      prefix: "!",
    },
    silent: APP_ENV !== "dev",
    botGuilds:
      APP_ENV === "dev"
        ? ["1265957323193716788", "1391117548036165752", "1369790516627247255"]
        : undefined,
  });

  client.once("ready", async () => {
    async function initCommands(__retries = 0) {
      if (__retries < 3) {
        try {
          await client.initApplicationCommands();
        } catch (err) {
          await client.clearApplicationCommands();
          await initCommands(__retries + 1);
          logger.error(err);
        }
      }
    }
    await initCommands();
  });

  client.on("interactionCreate", (interaction: Interaction) => {
    try {
      void client.executeInteraction(interaction);
    } catch (err) {
      logger.error(err);
    }
  });

  client.on("messageCreate", (message: Message) => {
    try {
      void client.executeCommand(message);
    } catch (err) {
      logger.error(err);
    }
  });

  await importx(`${dirname(import.meta.url)}/**/*.{ts,js}`);

  await client.login(configService.get("DISCORD_TOKEN")).then(() => {
    logger.log("Successfully logged in");
  });
}

bootstrap();
