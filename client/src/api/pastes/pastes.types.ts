import type { Pagination } from "#api/shared/index.js";

export interface CreatePastePayload {
  title: string;
  paste: string;
  userId: number;
}

export type UpdatePastePayload = Omit<CreatePastePayload, "userId">;

export interface Paste {
  id: number
  title: string;
  paste: string;
  userId: number;

  createdAt: string;
  updatedAt: string;
}

export interface PasteFilter {
  search: string;
  userId: number;
  strict: boolean;
  pasteId: number;
  socialId: string
}

export interface PasteQueryParams {
  pagination: Partial<Pagination>;
  filter: Partial<PasteFilter>;
}
