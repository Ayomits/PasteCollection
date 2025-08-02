export interface User {
  id: number;
  username: string;
  displayName: string;
  socialId: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserPayload {
  username: string;
  displayName: string;
  socialId: string;
}

export interface UserFilter {
  userId: number;
  username: string;
  displayName: string;
  socialId: string;
  matchAll: boolean;
  strict: boolean;
}

export type UpdateUserPayload = Omit<CreateUserPayload, "socialId">;
