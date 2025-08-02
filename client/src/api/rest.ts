import { LocalCache } from "@ts-fetcher/cache";
import { createRestInstance } from "@ts-fetcher/rest";
import type {
  RequestInterceptor,
  ResponseInterceptor,
} from "@ts-fetcher/types";
import { Logger } from "ayologger";

import { Env } from "#config/env.js";

const logger = new Logger({
  global: {
    template: () => `[API] | {date} | {message}`,
    level: {
      text: "API",
    },
    message: {
      text: "API",
    },
  },
  noBackground: true,
  theme: {
    debug: {
      level: {
        text: "API",
      },
    },
    info: {
      level: {
        text: "API",
      },
    },
    error: {
      level: {
        text: "API",
      },
    },
    warn: {
      level: {
        text: "API",
      },
    },
    success: {
      level: {
        text: "API",
      },
    },
  },
});

const LogRequestInterceptor: RequestInterceptor = (options) => {
  logger.info(`${options.method} | ${options.path} | incoming request`);
  return options;
};

const LogResponseInterceptor: ResponseInterceptor = (res) => {
  const { options, success, cached } = res;
  if (!success) {
    logger.error(`${options.method} | ${options.path} | failure request`);
  } else {
    logger.success(
      `${options.method} | ${options.path} | ${cached ? "CACHED" : "NEW"} | succesful request`
    );
  }
  return res;
};

export const rest = createRestInstance(Env.ApiUrl, {
  defaultRequestOptions: {
    headers: {
      Authorization: Env.ApiToken,
      "Content-Type": "application/json"
    },
  },
  caching: new LocalCache(),
  interceptors: {
    request: [LogRequestInterceptor],
    response: [LogResponseInterceptor],
  },
});
