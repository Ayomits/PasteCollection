import dt from "dotenv";

export class ConfigService {
  constructor() {
    dt.config();
  }

  public get<T = string>(key: string, default_?: T): T {
    return (process.env[key] ?? default_) as T;
  }

  public getOrThrow<T = string>(key: string, default_?: T): T {
    const existed = this.get<T>(key, default_);
    if (!existed) {
      throw new Error(`Environment key ${key} does not exists`);
    }
    return existed;
  }
}

export const configService = new ConfigService();
