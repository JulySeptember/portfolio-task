import type { ListTasksParams } from "../api/list-tasks";

export const taskQueryKeys = {
  all: ["tasks"] as const,

  lists: () => [...taskQueryKeys.all, "list"] as const,

  list: (params?: ListTasksParams) =>
    [...taskQueryKeys.lists(), params] as const,

  details: () => [...taskQueryKeys.all, "detail"] as const,

  detail: (id: number) => [...taskQueryKeys.details(), id] as const,
};
