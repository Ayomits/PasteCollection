import type { LiteralEnum } from "@ts-fetcher/types";

import { configService } from "./config.js";

export const AppEnv = {
  Dev: "dev",
  Prod: "prod",
} as const;

export type AppEnv = LiteralEnum<typeof AppEnv>;

export const Env = {
  DiscordToken: configService.getOrThrow("DISCORD_TOKEN"),
  AppEnv: configService.getOrThrow("APP_ENV") as AppEnv,
  ApiUrl: configService.getOrThrow("API_URL"),
  ApiToken: configService.getOrThrow("SECRET_API_TOKEN"),
} as const;
