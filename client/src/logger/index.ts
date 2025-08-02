import { Logger as AyoLogger } from "ayologger";
import type { ILogger } from "discordx";

class Logger implements ILogger {
  private base: AyoLogger;

  constructor() {
    this.base = new AyoLogger();
  }

  error(...args: unknown[]): void {
    this.base.error(args);
  }
  info(...args: unknown[]): void {
    this.base.success(args);
  }
  log(...args: unknown[]): void {
    this.base.info(args);
  }
  warn(...args: unknown[]): void {
    this.base.warn(args);
  }
}

export const createLogger = () => {
  return new Logger();
};
