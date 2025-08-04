import { codeBlock } from "discord.js";

import { PasteCreateMessages } from "./paste-create.messages.js";

export const PasteUpdateMessages = {
  success: {
    title: "Успех",
    description: "Поздравляю с обновлений пасты!",
  },
  validation: {
    notExists: {
      title: "Ошибка",
      description: "Указанной вами пасты не существует",
      fields: (text: string) => [
        {
          name: "Текст пасты",
          value: codeBlock(text),
        },
      ],
    },
    nan: {
      title: "Ошибка",
      value: "Указанный id не является числом",
    },
    ...PasteCreateMessages.validation,
  },
} as const;
