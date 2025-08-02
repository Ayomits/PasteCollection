import type { APIEmbed, EmbedData } from "discord.js";
import { EmbedBuilder as DjsEmbedBuild } from "discord.js";

export class EmbedBuilder extends DjsEmbedBuild {
  constructor(data?: EmbedData | APIEmbed) {
    super(data);

    super.setColor(0x2c2f33);
    super.setTimestamp(Date.now());
  }
}
