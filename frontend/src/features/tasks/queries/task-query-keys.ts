export const taskQueryKeys = {
  all: ["tasks"] as const,

  list: () => [...taskQueryKeys.all] as const,

  detail: (id: number) => [...taskQueryKeys.all, id] as const,
};
