import { rest } from "#api/rest.js";
import { BaseApi } from "#api/shared/base.js";
import type { ListResponse } from "#api/shared/index.js";

import type {
  CreatePastePayload,
  Paste,
  PasteFilter,
  PasteQueryParams,
  UpdatePastePayload,
} from "./pastes.types.js";

export class PastesApi extends BaseApi {
  async searchPaste(q: Partial<PasteQueryParams>) {
    return await rest.get<ListResponse<Paste>>(
      `/pastes/search` + this.getQuery(q)
    );
  }

  async findSignlePaste(f: Partial<PasteFilter>) {
    return await rest.get<Paste>(`/pastes` + this.getQuery(f), {});
  }

  async createPaste(p: CreatePastePayload) {
    return await rest.post<Paste, CreatePastePayload>("/pastes", {
      body: p,
    });
  }

  async updatePaste(f: Partial<PasteFilter>, p: UpdatePastePayload) {
    return await rest.put<Paste, UpdatePastePayload>(
      "/pastes" + this.getQuery(f),
      {
        body: p,
      }
    );
  }

  async deletePaste(f: Partial<PasteFilter>) {
    return await rest.delete("/pastes" + this.getQuery(f));
  }
}

export const pastesApi = new PastesApi();
