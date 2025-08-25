import "reflect-metadata";

import { dirname, importx } from "@discordx/importer";
import type { Interaction, Message } from "discord.js";
import { GatewayIntentBits } from "discord.js";
import { Client, DIService, tsyringeDependencyRegistryEngine } from "discordx";
import { container } from "tsyringe";

import { configService } from "#/config/config.js";
import { Env } from "#config/env.js";
import { createLogger } from "#logger/index.js";

async function bootstrap() {
  DIService.engine = tsyringeDependencyRegistryEngine.setInjector(container);
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
    silent: Env.AppEnv !== "dev",
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
