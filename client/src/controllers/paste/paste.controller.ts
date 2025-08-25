import {
  ApplicationCommandOptionType,
  ApplicationIntegrationType,
  type CommandInteraction,
  type ModalSubmitInteraction,
} from "discord.js";
import {
  Discord,
  ModalComponent,
  Slash,
  SlashGroup,
  SlashOption,
} from "discordx";
import { inject, singleton } from "tsyringe";

import { PasteCreateModalId, PasteUpdateModalId } from "./paste.const.js";
import { PasteService } from "./paste.service.js";

@Discord()
@SlashGroup({ name: "paste", description: "Управление пастами" })
@SlashGroup("paste")
@singleton()
export class PasteController {
  constructor(@inject(PasteService) private pasteService: PasteService) {}

  @Slash({
    name: "info",
    description: "Поиск по пастам",
    integrationTypes: [
      ApplicationIntegrationType.UserInstall,
      ApplicationIntegrationType.GuildInstall,
    ],
  })
  pasteInfoSlash(
    @SlashOption({
      description: "Поиск",
      name: "search",
      required: true,
      type: ApplicationCommandOptionType.String,
      autocomplete: async (interaction) => {
        const value = interaction.options.getFocused();
        return PasteService.pasteIdAutocomplete(interaction, {
          filter: {
            search: value,
          },
        });
      },
    })
    pasteId: string,
    interaction: CommandInteraction
  ) {
    return this.pasteService.infoSlash(interaction, pasteId);
  }

  @Slash({
    name: "create",
    description: "Создать новую пасту",
    integrationTypes: [
      ApplicationIntegrationType.UserInstall,
      ApplicationIntegrationType.GuildInstall,
    ],
  })
  pasteCreateSlash(interaction: CommandInteraction) {
    return this.pasteService.createSlash(interaction);
  }

  @ModalComponent({ id: PasteCreateModalId })
  pasteCreateModal(interaction: ModalSubmitInteraction) {
    return this.pasteService.createModal(interaction);
  }

  @Slash({
    name: "delete",
    description: "Удалить пасту",
    integrationTypes: [
      ApplicationIntegrationType.UserInstall,
      ApplicationIntegrationType.GuildInstall,
    ],
  })
  async pasteDeleteSlash(
    @SlashOption({
      description: "Поиск",
      name: "search",
      required: true,
      type: ApplicationCommandOptionType.String,
      autocomplete: async (interaction) => {
        const value = interaction.options.getFocused();
        return PasteService.pasteIdAutocomplete(interaction, {
          filter: {
            search: value,
            socialId: interaction.user.id,
          },
        });
      },
    })
    pasteId: string,
    interaction: CommandInteraction
  ) {
    return this.pasteService.deleteSlash(interaction, pasteId);
  }

  @ModalComponent({ id: PasteUpdateModalId })
  pasteUpdateModal(interaction: ModalSubmitInteraction) {
    return this.pasteService.updateModal(interaction);
  }

  @Slash({
    name: "update",
    description: "Удалить пасту",
    integrationTypes: [
      ApplicationIntegrationType.UserInstall,
      ApplicationIntegrationType.GuildInstall,
    ],
  })
  async pasteUpdateSlash(
    @SlashOption({
      description: "Поиск",
      name: "search",
      required: true,
      type: ApplicationCommandOptionType.String,
      autocomplete: async (interaction) => {
        const value = interaction.options.getFocused();
        return PasteService.pasteIdAutocomplete(interaction, {
          filter: {
            search: value,
            socialId: interaction.user.id,
          },
        });
      },
    })
    pasteId: string,
    interaction: CommandInteraction
  ) {
    return this.pasteService.updateSlash(interaction, pasteId);
  }
}
