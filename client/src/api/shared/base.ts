import QueryString from "qs";

export class BaseApi {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  protected getQuery(q: any) {
    return "?" + QueryString.stringify(q, {});
  }
}
