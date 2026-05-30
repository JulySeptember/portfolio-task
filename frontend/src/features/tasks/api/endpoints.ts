export const taskEndpoints = {
  list: () => "/api/v1/tasks",

  detail: (publicId: string) => `/api/v1/tasks/${publicId}`,

  status: (publicId: string) => `/api/v1/tasks/${publicId}/status`,

  create: () => "/api/v1/tasks",

  update: (publicId: string) => `/api/v1/tasks/${publicId}`,

  delete: (publicId: string) => `/api/v1/tasks/${publicId}`,
};
