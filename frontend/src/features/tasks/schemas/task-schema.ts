// src/features/tasks/schemas/task-schema.ts

import { z } from "zod";

// =========================
// enums
// =========================

export const taskStatusSchema = z.enum(["TODO", "DOING", "DONE"]);

// =========================
// api task schema
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

    due_date: z.string().datetime().nullable(),

    created_at: z.string().datetime(),

    updated_at: z.string().datetime(),
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
// task list response
// =========================

export const taskListResponseSchema = z.object({
  count: z.number().int().nonnegative(),

  items: z.array(taskSchema),

  limit: z.number().int().positive(),

  offset: z.number().int().nonnegative(),
});

// =========================
// form schema
// =========================

export const taskFormSchema = z
  .object({
    title: z
      .string()
      .trim()
      .min(1, "Title is required")
      .max(255, "Title is too long"),

    description: z.string().max(5000, "Description is too long"),

    status: taskStatusSchema,

    due_date: z
      .string()
      .optional()
      .refine(
        (value) => {
          if (!value) return true;

          return !Number.isNaN(new Date(value).getTime());
        },
        {
          message: "Invalid date",
        },
      ),
  })
  .transform((data) => ({
    title: data.title,

    description: data.description,

    status: data.status,

    due_date: data.due_date ? new Date(data.due_date).toISOString() : undefined,
  }));

// =========================
// request schema
// =========================

export const taskRequestSchema = z.object({
  title: z
    .string()
    .trim()
    .min(1, "Title is required")
    .max(255, "Title is too long"),

  description: z.string().max(5000, "Description is too long"),

  status: taskStatusSchema,

  due_date: z.string().datetime().nullable(),
});

// =========================
// types
// =========================

export type TaskStatus = z.infer<typeof taskStatusSchema>;

export type Task = z.infer<typeof taskSchema>;

export type TaskListResponse = z.infer<typeof taskListResponseSchema>;

export type TaskFormInput = z.input<typeof taskFormSchema>;

export type TaskFormValues = z.output<typeof taskFormSchema>;

export type TaskRequest = z.infer<typeof taskRequestSchema>;
