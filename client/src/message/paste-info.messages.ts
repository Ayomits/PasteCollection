import { codeBlock } from "discord.js";

import type { Paste } from "#api/pastes/pastes.types.js";
import { usersApi } from "#api/users/users.api.js";

export const PasteInfoMessages = {
  validation: {
    title: "Ошибка",
    nan: "Невалидный id",
    nullable: "Паста не найдена",
    internal: "Внутренняя ошибка"
  },
  success: {
    title: (title: string) => `Паста ${title}`,
    fields: async (paste: Paste) => {
      const user = await usersApi.findSignleUser({
        userId: paste.userId,
      });
      return [
        {
          name: "Текст пасты",
          value: codeBlock(paste.paste),
        },
        {
          name: "Автор",
          value: codeBlock(user.data.displayName),
        },
      ];
    },
  },
} as const;
