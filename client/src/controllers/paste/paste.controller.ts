import {
  ApplicationCommandOptionType,
  type CommandInteraction,
  type ModalSubmitInteraction,
} from "discord.js";
import { Discord, ModalComponent, Slash, SlashOption } from "discordx";
import { inject, singleton } from "tsyringe";

import { PasteCreateModalId, PasteUpdateModalId } from "./paste.const.js";
import { PasteService } from "./paste.service.js";

@Discord()
@singleton()
export class PasteController {
  constructor(@inject(PasteService) private pasteService: PasteService) {}

  @Slash({
    name: "paste-info",
    description: "Поиск по пастам",
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
    name: "paste-create",
    description: "Создать новую пасту",
  })
  pasteCreateSlash(interaction: CommandInteraction) {
    return this.pasteService.createSlash(interaction);
  }

  @ModalComponent({ id: PasteCreateModalId })
  pasteCreateModal(interaction: ModalSubmitInteraction) {
    return this.pasteService.createModal(interaction);
  }

  @Slash({
    name: "paste-delete",
    description: "Удалить пасту",
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
    name: "paste-update",
    description: "Удалить пасту",
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
