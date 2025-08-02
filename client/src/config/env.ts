import { configService } from "./config.js";

export const Env = {
  DiscordToken: configService.getOrThrow("DISCORD_TOKEN"),
  AppEnv: configService.getOrThrow("APP_ENV"),
  ApiUrl: configService.getOrThrow("API_URL"),
  ApiToken: configService.getOrThrow("API_TOKEN"),
} as const;
