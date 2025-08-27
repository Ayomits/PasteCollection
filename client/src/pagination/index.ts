import type { LiteralEnum } from "@ts-fetcher/types";
import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  type CacheType,
  type CollectedInteraction,
  type Interaction,
  type InteractionCollectorOptions,
  type InteractionEditReplyOptions,
  type InteractionReplyOptions,
  type RepliableInteraction,
} from "discord.js";

const PaginationButtonType = {
  Next: "next",
  Prev: "prev",
} as const;

const nextCustomId = "@ayo-next";
const prevCustomId = "@ayo-prev";

type PaginationButtonType = LiteralEnum<typeof PaginationButtonType>;

interface DynamicPaginationButton {
  label?: string;
  emoji?: string;
  style: ButtonStyle;
  type: PaginationButtonType;
}

interface DynamicPaginationOptions<T> {
  fetchNext: (prev: T) => Promise<T> | T;
  canFetchNext: (prev: T) => boolean;
  fetchInitial: () => T | Promise<T>;
  buildMessage: (
    data: T,
    pageNumber: number,
  ) =>
    | Promise<InteractionEditReplyOptions | InteractionReplyOptions>
    | InteractionReplyOptions
    | InteractionEditReplyOptions;
  buttons?: DynamicPaginationButton[];
}

export class DynamicPagination<T> {
  private config: DynamicPaginationOptions<T>;
  private messages: (InteractionEditReplyOptions | InteractionReplyOptions)[] =
    [];
  private data: T[] = [];

  constructor(options: DynamicPaginationOptions<T>) {
    this.config = options;
  }

  public async send(
    interaction: Interaction,
    options?: InteractionCollectorOptions<CollectedInteraction, CacheType> & {
      callback?: (interaction: Interaction) => unknown | Promise<unknown>;
    },
  ) {
    if (interaction.isRepliable()) {
      const { buildMessage, fetchInitial, canFetchNext } = this.config;

      const initial = await fetchInitial();
      const payload = await buildMessage(initial, 0);

      this.data.push(initial);
      this.messages.push(payload);

      const mergeComponents = (
        payload: InteractionEditReplyOptions | InteractionReplyOptions,
        canNext: boolean,
      ) => {
        const rows = [];

        if (payload.components?.length) {
          rows.push(...payload.components);
        }

        rows.push(
          this.buildButtons({
            prev: true,
            next: canNext,
          }),
        );

        return rows;
      };

      // @ts-expect-error trust me
      const reply = await this.sendReply(interaction, {
        ...payload,
        components: [...mergeComponents(payload, !canFetchNext(initial))],
      });

      //@ts-expect-error trust me
      const collector = reply.createMessageComponentCollector({
        filter: (i, c) =>
          i.message.id === reply.id &&
          (options?.filter ? options?.filter(i, c) : true),
        ...options,
      });

      let pageNumber = 0;

      collector.on("collect", async (interaction) => {
        const customId = interaction.customId;
        let data: T;
        let canNext: boolean = false;
        let payload: InteractionEditReplyOptions | InteractionReplyOptions;

        if (customId !== nextCustomId && customId !== prevCustomId) {
          return options?.callback?.(interaction);
        }

        await interaction.deferUpdate();

        if (customId === nextCustomId) {
          data = await this.config.fetchNext(this.data[pageNumber]!);
          pageNumber++;
          canNext = canFetchNext(data);
          payload = await buildMessage(data, pageNumber);
          this.data.push(data!);
          this.messages.push(payload);
        }

        if (customId === prevCustomId) {
          pageNumber--;
          data = this.data[pageNumber]!;
          canNext = canFetchNext(data);
          payload = this.messages[pageNumber]!;
        }

        // @ts-expect-error trust me
        return interaction.editReply({
          ...payload!,
          components: mergeComponents(payload!, !canNext),
        });
      });
    }
  }

  private buildButtons(disables?: { next?: boolean; prev?: boolean }) {
    const { buttons = [] } = this.config;

    if (buttons.length > 0) {
      return new ActionRowBuilder<ButtonBuilder>().addComponents(
        ...buttons.slice(0, 5).map((btn) => {
          const builder = new ButtonBuilder();

          if (btn.label) {
            builder.setLabel(btn.label);
          }

          if (btn.emoji) {
            builder.setEmoji(btn.emoji);
          }

          if (!btn.label && !btn.emoji) {
            throw new Error("Label or emoji required");
          }

          if (btn.style) {
            builder.setStyle(btn.style);
          }

          if (btn.type === "next") {
            builder.setCustomId(nextCustomId);
            builder.setDisabled(disables?.next);
          }

          if (btn.type === "prev") {
            builder.setCustomId(prevCustomId);
            builder.setDisabled(disables?.prev);
          }

          return builder;
        }),
      );
    }

    const defaultNextButton = new ButtonBuilder()
      .setCustomId(nextCustomId)
      .setEmoji("➡️")
      .setDisabled(disables?.next)
      .setStyle(ButtonStyle.Secondary);

    const defaultPrevButton = new ButtonBuilder()
      .setCustomId(prevCustomId)
      .setEmoji("⬅️")
      .setDisabled(disables?.prev)
      .setStyle(ButtonStyle.Secondary);

    return new ActionRowBuilder<ButtonBuilder>().addComponents(
      defaultPrevButton,
      defaultNextButton,
    );
  }

  private async sendReply(
    interaction: RepliableInteraction,
    options: InteractionReplyOptions,
  ) {
    const tryEdit = () => {
      try {
        return interaction.editReply(options as InteractionEditReplyOptions);
      } catch {
        return interaction.followUp(options);
      }
    };

    if (interaction.deferred || interaction.replied) {
      return await tryEdit();
    }
    return await interaction.reply(options);
  }
}
