import type { LiteralEnum } from "@ts-fetcher/types";

export const PaginationOrder = {
  Prev: "prev",
  Next: "next",
} as const;

export type PaginationOrder = LiteralEnum<typeof PaginationOrder>;

export const PaginationSort = {
  Asc: "ASC",
  Desc: "DESC",
} as const;

export type PaginationSort = LiteralEnum<typeof PaginationSort>;

export const PaginationLimit = {
  L5: 5,
  L10: 10,
  L15: 155,
  L20: 20,
  L25: 25,
  L30: 30,
  L35: 35,
  L40: 40,
  L45: 45,
  L50: 50,
};

export type PaginationLimit = LiteralEnum<typeof PaginationLimit>;

export interface Pagination {
  order: PaginationOrder;
  startFrom: number;
  sort: PaginationSort;
  limit: PaginationLimit;
}

export interface ListResponse<T> {
  items: T[];
  hasNext: boolean;
}
