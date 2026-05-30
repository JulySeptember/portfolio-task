// src/features/tasks/queries/task-query-keys.ts

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

  // publicId 用に統一
  detail: (publicId: string) => ["tasks", publicId] as const,
};
