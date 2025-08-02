import dt from "dotenv";

export class ConfigService {
  private entries: dt.DotenvParseOutput;

  constructor(options?: dt.DotenvConfigOptions) {
    const config = dt.config(options);
    if (config.error || !config) {
      throw new Error(`Failded to parse config \n ${options}`);
    }
    this.entries = config.parsed!;
  }

  public get<T = string>(key: string, default_?: T): T {
    const existed = this.entries[key] ?? default_;
    return existed as T;
  }

  public getOrThrow<T = string>(key: string): T {
    const existed = this.get<T>(key);
    if (!existed) {
      throw new Error(`Key ${key} does not exists`);
    }
    return existed;
  }
}

export const configService = new ConfigService();
