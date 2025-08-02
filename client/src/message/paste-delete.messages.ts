import { PasteInfoMessages } from "./paste-info.messages.js";

export const PasteDeleteMessages = {
  success: {
    title: "Успешно удалена паста",
    description: "Вы удалили свою пасту",
  },
  validation: {
    ...PasteInfoMessages.validation,
  },
} as const;
