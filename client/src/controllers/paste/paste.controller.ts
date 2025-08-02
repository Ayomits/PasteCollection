import {
  ActionRowBuilder,
  ApplicationCommandOptionType,
  type AutocompleteInteraction,
  type CommandInteraction,
  ModalBuilder,
  type ModalSubmitInteraction,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";
import { Discord, ModalComponent, Slash, SlashOption } from "discordx";

import { pastesApi } from "#api/pastes/pastes.api.js";
import { usersApi } from "#api/users/users.api.js";
import { PasteCreateMessages } from "#message/paste-create.messages.js";
import { PasteDeleteMessages } from "#message/paste-delete.messages.js";
import { PasteInfoMessages } from "#message/paste-info.messages.js";
import { EmbedBuilder } from "#utils/embed.builder.js";
import { UsersUtility } from "#utils/user.utility.js";

@Discord()
export class PasteController {
  @Slash({
    name: "paste-info",
    description: "Поиск по пастам",
  })
  async pasteInfo(
    @SlashOption({
      description: "Поиск",
      name: "search",
      required: true,
      type: ApplicationCommandOptionType.String,
      autocomplete: async function (interaction: AutocompleteInteraction) {
        try {
          const value = interaction.options.getFocused();
          const entries = await pastesApi.searchPaste({
            filter: {
              search: value,
            },
            pagination: {
              limit: 25,
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
      },
    })
    pasteId: string,
    interaction: CommandInteraction
  ) {
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

  @Slash({
    name: "paste-create",
    description: "Создать новую пасту",
  })
  async pasteCreate(interaction: CommandInteraction) {
    const title = new ActionRowBuilder<TextInputBuilder>().addComponents(
      new TextInputBuilder()
        .setCustomId("title")
        .setLabel("Название пасты")
        .setPlaceholder("О любви")
        .setMinLength(1)
        .setMaxLength(32)
        .setRequired(true)
        .setStyle(TextInputStyle.Short)
    );
    const paste = new ActionRowBuilder<TextInputBuilder>().addComponents(
      new TextInputBuilder()
        .setCustomId("paste")
        .setLabel("Текст пасты")
        .setPlaceholder("Имел отец слепого сына...")
        .setMinLength(1)
        .setMaxLength(2096)
        .setRequired(true)
        .setStyle(TextInputStyle.Paragraph)
    );
    const modal = new ModalBuilder()
      .setTitle("Создать новую пасту")
      .setCustomId("create-paste")
      .setComponents(title, paste);
    return interaction.showModal(modal);
  }

  @ModalComponent({ id: "create-paste" })
  async pasteCreateModal(interaction: ModalSubmitInteraction) {
    await interaction.deferReply();
    const embed = new EmbedBuilder();

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

  @Slash({
    name: "paste-delete",
    description: "Создать новую пасту",
  })
  async pasteDelete(
    @SlashOption({
      description: "Поиск",
      name: "search",
      required: true,
      type: ApplicationCommandOptionType.String,
      autocomplete: async function (interaction: AutocompleteInteraction) {
        try {
          const value = interaction.options.getFocused();
          const entries = await pastesApi.searchPaste({
            filter: {
              search: value,
              socialId: interaction.user.id,
            },
            pagination: {
              limit: 10,
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
      },
    })
    pasteId: string,
    interaction: CommandInteraction
  ) {
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
}
