import type { ListTasksParams } from "../api/list-tasks";

export const taskQueryKeys = {
  all: ["tasks"] as const,

  lists: () => [...taskQueryKeys.all, "list"] as const,

  list: (params?: ListTasksParams) =>
    [
      ...taskQueryKeys.lists(),
      {
        limit: params?.limit ?? 20,
        offset: params?.offset ?? 0,
        status: params?.status ?? null,
        sort: params?.sort ?? "created_at",
        order: params?.order ?? "DESC",
      },
    ] as const,

  details: () => [...taskQueryKeys.all, "detail"] as const,

  detail: (id: number) => [...taskQueryKeys.details(), id] as const,
};
