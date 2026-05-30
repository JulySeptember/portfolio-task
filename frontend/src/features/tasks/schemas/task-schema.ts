import { z } from "zod";

// =========================
// constants
// =========================

export const TITLE_MAX_LENGTH = 255;

export const DESCRIPTION_MAX_LENGTH = 5000;

// =========================
// enums
// =========================

export const taskStatusSchema = z.enum(["TODO", "DOING", "DONE"]);

export type TaskStatus = z.infer<typeof taskStatusSchema>;

// =========================
// API Task Schema
// =========================

export const taskSchema = z
  .object({
    publicId: z.string().optional(), // API は public_id を返す場合あり
    public_id: z.string().optional(),

    title: z.string().min(1, "Title is required").max(255, "Title too long"),

    description: z.string().nullable().optional(),

    status: taskStatusSchema,

    dueDate: z.string().nullable().optional(),
    due_date: z.string().nullable().optional(),

    createdAt: z.string().optional(),
    created_at: z.string().optional(),

    updatedAt: z.string().optional(),
    updated_at: z.string().optional(),
  })
  .transform((task) => ({
    publicId: task.publicId ?? task.public_id ?? "",
    title: task.title,
    description: task.description ?? null,
    status: task.status,
    dueDate: task.dueDate ?? task.due_date ?? null,
    createdAt: task.createdAt ?? task.created_at ?? "",
    updatedAt: task.updatedAt ?? task.updated_at ?? "",
  }));

export type Task = z.infer<typeof taskSchema>;

// =========================
// Task List Response
// =========================

export const taskListResponseSchema = z.object({
  items: z.array(taskSchema),
  count: z.number(),
  limit: z.number(),
  offset: z.number(),
});

export type TaskListResponse = z.infer<typeof taskListResponseSchema>;

// =========================
// Form Schema
// =========================

export const taskFormSchema = z.object({
  title: z
    .string()
    .trim()
    .min(1, "Title is required")
    .max(TITLE_MAX_LENGTH, `${TITLE_MAX_LENGTH} characters max`),
  description: z
    .string()
    .trim()
    .max(DESCRIPTION_MAX_LENGTH, `${DESCRIPTION_MAX_LENGTH} characters max`)
    .optional(),

  status: taskStatusSchema,

  due_date: z.string().nullable().optional(),
});

export type TaskFormInput = z.input<typeof taskFormSchema>;

export type TaskFormValues = z.output<typeof taskFormSchema>;
