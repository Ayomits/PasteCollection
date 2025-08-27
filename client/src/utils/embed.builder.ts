import type { APIEmbed, EmbedData, User } from "discord.js";
import { EmbedBuilder as DjsEmbedBuild } from "discord.js";

import { UsersUtility } from "./user.utility.js";

export class EmbedBuilder extends DjsEmbedBuild {
  constructor(data?: EmbedData | APIEmbed) {
    super(data);

    super.setColor(0x2c2f33);
    super.setTimestamp(Date.now());
  }

  setDefaults(user: User): this {
    const avatar = UsersUtility.getAvatar(user);
    const username = UsersUtility.getUsername(user);

    return this.setThumbnail(avatar).setFooter({
      text: username,
      iconURL: avatar,
    });
  }
}
