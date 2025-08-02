import type { GuildMember, User } from "discord.js";

export type UsersUtilityAccept = User | GuildMember;

export class UsersUtility {
  static getAvatar(user: UsersUtilityAccept) {
    return user.displayAvatarURL() ?? user.avatarURL() ?? undefined;
  }

  static getUsername(user: UsersUtilityAccept) {
    const transformedUser = "user" in user ? user.user : user;
    const nick = "nickname" in user ? user.nickname : null;
    const globalName = "globalName" in user ? user.globalName : null;
    const username =
      "username" in transformedUser ? transformedUser.username : null;

    return nick ?? user.displayName ?? globalName ?? username;
  }
}
