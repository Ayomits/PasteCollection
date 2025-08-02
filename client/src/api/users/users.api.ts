import type { ApiResponse, HttpMethodType } from "@ts-fetcher/types";
import type { User as DiscordUser } from "discord.js";

import { rest } from "#api/rest.js";
import { BaseApi } from "#api/shared/base.js";
import { UsersUtility } from "#utils/user.utility.js";

import type {
  CreateUserPayload,
  UpdateUserPayload,
  User,
  UserFilter,
} from "./users.types.js";

export class UsersApi extends BaseApi {
  async findOrCreate(
    usr: DiscordUser
  ): Promise<ApiResponse<User, unknown, HttpMethodType>> {
    const existed = await this.findSignleUser({
      socialId: usr.id,
    });
    if (!existed || !existed.success) {
      const newUsr = await this.createUser({
        displayName: UsersUtility.getUsername(usr),
        username: UsersUtility.getUsername(usr),
        socialId: usr.id,
      });
      return newUsr;
    }
    return existed;
  }

  async findSignleUser(f: Partial<UserFilter>) {
    return await rest.get<User>(`/users` + this.getQuery(f));
  }

  async createUser(p: CreateUserPayload) {
    return await rest.post<User, CreateUserPayload>("/users", {
      body: p,
    });
  }

  async updateUser(f: Partial<UserFilter>, p: UpdateUserPayload) {
    return await rest.put<User, UpdateUserPayload>(
      "/users" + this.getQuery(f),
      {
        body: p,
      }
    );
  }

  async deleteUser(f: Partial<UserFilter>) {
    return await rest.delete("/users" + this.getQuery(f));
  }
}

export const usersApi = new UsersApi();
