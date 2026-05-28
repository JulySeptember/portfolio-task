// src/features/tasks/api/endpoints.ts

const TASKS_API_BASE = "/api/v1/tasks";

export const taskEndpoints = {
  list: () => TASKS_API_BASE,

  detail: (id: number) => `${TASKS_API_BASE}/${id}`,

  status: (id: number) => `${TASKS_API_BASE}/${id}/status`,
};
