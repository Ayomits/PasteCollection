import { codeBlock } from "discord.js";

export const PasteCreateMessages = {
  success: {
    title: "Успех",
    description: "Поздравляю с созданием вашей пасты!",
  },
  validation: {
    unique: {
      title: "Паста с таким названием уже существует",
      fields: (text: string) => [
        { name: "Текст пасты", value: codeBlock(text) },
      ],
    },
    internal: {
      title: "Произошла внутренняя ошибка",
      description: "Что-то пошло не так...",
    },
  },
} as const;
