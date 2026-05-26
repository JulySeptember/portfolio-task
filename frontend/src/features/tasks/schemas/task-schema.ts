import { z } from "zod";

// =========================
// enums
// =========================

export const taskStatusSchema = z.enum(["TODO", "DOING", "DONE"]);

// =========================
// api task schema
// backend response shape
// =========================

export const taskSchema = z
  .object({
    id: z.number().int().positive(),

    user_id: z.number().int().positive(),

    title: z
      .string()
      .trim()
      .min(1, "Title is required")
      .max(255, "Title is too long"),

    description: z.string().max(5000, "Description is too long"),

    status: taskStatusSchema,

    due_date: z.iso.datetime().nullable(),

    created_at: z.iso.datetime(),

    updated_at: z.iso.datetime(),
  })
  .transform((data) => ({
    id: data.id,

    userId: data.user_id,

    title: data.title,

    description: data.description,

    status: data.status,

    dueDate: data.due_date,

    createdAt: data.created_at,

    updatedAt: data.updated_at,
  }));

// =========================
// task list
// =========================

export const taskListSchema = z.object({
  count: z.number().int().nonnegative(),

  items: z.array(taskSchema),

  limit: z.number().int().positive(),

  offset: z.number().int().nonnegative(),
});

// =========================
// create/update request
// backend request shape
// =========================

export const taskRequestSchema = z.object({
  title: z
    .string()
    .trim()
    .min(1, "Title is required")
    .max(255, "Title is too long"),

  description: z.string().max(5000, "Description is too long"),

  status: taskStatusSchema,

  due_date: z.iso.datetime().nullable(),
});

// =========================
// frontend types
// camelCase
// =========================

export type TaskStatus = z.infer<typeof taskStatusSchema>;

export type Task = z.infer<typeof taskSchema>;

export type TaskListResponse = z.infer<typeof taskListSchema>;

export type TaskRequest = z.infer<typeof taskRequestSchema>;
