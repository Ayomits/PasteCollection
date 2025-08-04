import {
  ActionRowBuilder,
  type AutocompleteInteraction,
  type CommandInteraction,
  ModalBuilder,
  type ModalSubmitInteraction,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";
import { injectable } from "tsyringe";

import { pastesApi } from "#api/pastes/pastes.api.js";
import type { PasteQueryParams } from "#api/pastes/pastes.types.js";
import { PaginationLimit } from "#api/shared/index.js";
import { usersApi } from "#api/users/users.api.js";
import { PasteCreateMessages } from "#message/paste-create.messages.js";
import { PasteDeleteMessages } from "#message/paste-delete.messages.js";
import { PasteInfoMessages } from "#message/paste-info.messages.js";
import { PasteUpdateMessages } from "#message/paste-update.messages.js";
import { EmbedBuilder } from "#utils/embed.builder.js";
import { UsersUtility } from "#utils/user.utility.js";

import {
  MaxPasteTextLength,
  MaxPasteTitleLength,
  PasteCreateModalId,
  PasteUpdateModalId,
} from "./paste.const.js";

@injectable()
export class PasteService {
  async infoSlash(interaction: CommandInteraction, pasteId: string) {
    await interaction.deferReply();
    const usrname = UsersUtility.getUsername(interaction.user);
    const avatar = UsersUtility.getAvatar(interaction.user);
    const embed = new EmbedBuilder()
      .setThumbnail(avatar)
      .setFooter({ text: usrname, iconURL: avatar });
    const numPasteId = Number(pasteId);
    if (Number.isNaN(numPasteId)) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteInfoMessages.validation.title)
            .setDescription(PasteInfoMessages.validation.nan),
        ],
      });
    }
    const paste = await pastesApi.findSignlePaste({
      pasteId: Number(pasteId),
    });
    if (!paste || !paste.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteInfoMessages.validation.title)
            .setDescription(PasteInfoMessages.validation.nullable),
        ],
      });
    }

    return interaction.editReply({
      embeds: [
        embed
          .setTitle(PasteInfoMessages.success.title(paste.data.title))
          .setFields(await PasteInfoMessages.success.fields(paste.data)),
      ],
    });
  }

  createSlash(interaction: CommandInteraction) {
    const title = new ActionRowBuilder<TextInputBuilder>().addComponents(
      new TextInputBuilder()
        .setCustomId("title")
        .setLabel("Название пасты")
        .setPlaceholder("О любви")
        .setMinLength(1)
        .setMaxLength(MaxPasteTitleLength)
        .setRequired(true)
        .setStyle(TextInputStyle.Short)
    );
    const paste = new ActionRowBuilder<TextInputBuilder>().addComponents(
      new TextInputBuilder()
        .setCustomId("paste")
        .setLabel("Текст пасты")
        .setPlaceholder("Имел отец слепого сына...")
        .setMinLength(1)
        .setMaxLength(MaxPasteTextLength)
        .setRequired(true)
        .setStyle(TextInputStyle.Paragraph)
    );
    const modal = new ModalBuilder()
      .setTitle("Создать новую пасту")
      .setCustomId(PasteCreateModalId)
      .setComponents(title, paste);
    return interaction.showModal(modal);
  }

  async createModal(interaction: ModalSubmitInteraction) {
    await interaction.deferReply();
    const usrname = UsersUtility.getUsername(interaction.user);
    const avatar = UsersUtility.getAvatar(interaction.user);
    const embed = new EmbedBuilder()
      .setThumbnail(avatar)
      .setFooter({ text: usrname, iconURL: avatar });

    const [title, text] = [
      interaction.fields.getTextInputValue("title"),
      interaction.fields.getTextInputValue("paste"),
    ];

    const user = await usersApi.findOrCreate(interaction.user);

    if (!user || !user.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteCreateMessages.validation.internal.title)
            .setDescription(
              PasteCreateMessages.validation.internal.description
            ),
        ],
      });
    }

    const existed = await pastesApi.findSignlePaste({
      search: title.trim(),
      strict: true,
    });

    if (existed.success && existed.data) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteCreateMessages.validation.unique.title)
            .setFields(PasteCreateMessages.validation.unique.fields(text)),
        ],
      });
    }

    const paste = await pastesApi.createPaste({
      title: title.trim(),
      paste: text,
      userId: user.data.id,
    });

    if (!paste || !paste.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteCreateMessages.validation.internal.title)
            .setDescription(
              PasteCreateMessages.validation.internal.description
            ),
        ],
      });
    }

    return interaction.editReply({
      embeds: [
        embed
          .setTitle(PasteCreateMessages.success.title)
          .setDescription(PasteCreateMessages.success.description),
      ],
    });
  }

  async updateSlash(interaction: CommandInteraction, pasteId: string) {
    const usrname = UsersUtility.getUsername(interaction.user);
    const avatar = UsersUtility.getAvatar(interaction.user);
    const embed = new EmbedBuilder()
      .setThumbnail(avatar)
      .setFooter({ text: usrname, iconURL: avatar });
    const numPasteId = Number(pasteId);
    if (Number.isNaN(numPasteId)) {
      return interaction.reply({
        embeds: [
          embed
            .setTitle(PasteUpdateMessages.validation.nan.title)
            .setDescription(PasteUpdateMessages.validation.nan.value),
        ],
      });
    }

    const existed = await pastesApi.findSignlePaste({
      pasteId: numPasteId,
    });

    if (!existed.success) {
      return interaction.reply({
        embeds: [
          embed
            .setTitle(PasteUpdateMessages.validation.notExists.title)
            .setDescription(
              PasteUpdateMessages.validation.notExists.description
            ),
        ],
      });
    }

    const modal = new ModalBuilder()
      .setCustomId(PasteUpdateModalId)
      .setTitle("Обновление пасты");
    const pasteTitleField =
      new ActionRowBuilder<TextInputBuilder>().addComponents(
        new TextInputBuilder()
          .setCustomId("title")
          .setLabel("Название")
          .setValue(existed.data.title)
          .setMinLength(1)
          .setMaxLength(MaxPasteTitleLength)
          .setStyle(TextInputStyle.Short)
      );
    const pasteTextField =
      new ActionRowBuilder<TextInputBuilder>().addComponents(
        new TextInputBuilder()
          .setCustomId("paste")
          .setLabel("Текст пасты")
          .setValue(existed.data.paste)
          .setMinLength(1)
          .setMaxLength(MaxPasteTextLength)
          .setStyle(TextInputStyle.Paragraph)
      );

    const pasteIdField = new ActionRowBuilder<TextInputBuilder>().addComponents(
      new TextInputBuilder()
        .setCustomId("pasteid")
        .setLabel("Текст пасты")
        .setValue(existed.data.id.toString())
        .setPlaceholder("Если ты это видишь - верни как было")
        .setMinLength(1)
        .setStyle(TextInputStyle.Paragraph)
    );

    return interaction.showModal(
      modal.addComponents(pasteTitleField, pasteTextField, pasteIdField)
    );
  }

  async updateModal(interaction: ModalSubmitInteraction) {
    await interaction.deferReply();
    const [title, paste, pasteId] = [
      interaction.fields.getTextInputValue("title"),
      interaction.fields.getTextInputValue("paste"),
      interaction.fields.getTextInputValue("pasteid"),
    ];

    const usrname = UsersUtility.getUsername(interaction.user);
    const avatar = UsersUtility.getAvatar(interaction.user);
    const embed = new EmbedBuilder()
      .setThumbnail(avatar)
      .setFooter({ text: usrname, iconURL: avatar });
    const numPasteId = Number(pasteId);
    if (Number.isNaN(numPasteId)) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteUpdateMessages.validation.nan.title)
            .setDescription(PasteUpdateMessages.validation.nan.value),
        ],
      });
    }

    const existed = await pastesApi.findSignlePaste({
      pasteId: numPasteId,
    });

    if (!existed.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteUpdateMessages.validation.notExists.title)
            .setFields(PasteUpdateMessages.validation.notExists.fields(paste)),
        ],
      });
    }

    const newExisted = await pastesApi.findSignlePaste({
      search: title.trim(),
      strict: true,
    });

    if (newExisted.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteUpdateMessages.validation.unique.title)
            .setFields(PasteUpdateMessages.validation.unique.fields(paste)),
        ],
      });
    }

    const updated = await pastesApi.updatePaste(
      {
        pasteId: existed.data.id,
      },
      {
        title,
        paste,
      }
    );

    if (!updated.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteUpdateMessages.validation.internal.title)
            .setFields(PasteUpdateMessages.validation.internal.fields(paste)),
        ],
      });
    }

    return interaction.editReply({
      embeds: [
        embed
          .setTitle(PasteUpdateMessages.success.title)
          .setDescription(PasteUpdateMessages.success.description),
      ],
    });
  }

  async deleteSlash(interaction: CommandInteraction, pasteId: string) {
    await interaction.deferReply();
    const usrname = UsersUtility.getUsername(interaction.user);
    const avatar = UsersUtility.getAvatar(interaction.user);
    const embed = new EmbedBuilder()
      .setThumbnail(avatar)
      .setFooter({ text: usrname, iconURL: avatar });
    const numPasteId = Number(pasteId);
    if (Number.isNaN(numPasteId)) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteDeleteMessages.validation.title)
            .setDescription(PasteDeleteMessages.validation.nan),
        ],
      });
    }
    const paste = await pastesApi.findSignlePaste({
      pasteId: Number(pasteId),
    });
    if (!paste || !paste.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteDeleteMessages.validation.title)
            .setDescription(PasteDeleteMessages.validation.nullable),
        ],
      });
    }

    const deleted = await pastesApi.deletePaste({
      pasteId: Number(pasteId),
    });

    if (!deleted || !deleted.success) {
      return interaction.editReply({
        embeds: [
          embed
            .setTitle(PasteDeleteMessages.validation.title)
            .setDescription(PasteDeleteMessages.validation.nullable),
        ],
      });
    }

    return interaction.editReply({
      embeds: [
        embed
          .setTitle(PasteDeleteMessages.success.title)
          .setDescription(PasteDeleteMessages.success.description),
      ],
    });
  }

  static async pasteIdAutocomplete(
    interaction: AutocompleteInteraction,
    query: Partial<PasteQueryParams>
  ) {
    try {
      const entries = await pastesApi.searchPaste({
        ...query,
        pagination: {
          ...query.pagination,
          limit: PaginationLimit.L25,
        },
      });
      return interaction.respond(
        entries.data?.items?.map((item) => ({
          name: item.title,
          value: item.id.toString(),
        })) ?? []
      );
    } catch {
      return interaction.respond([]);
    }
  }
}
