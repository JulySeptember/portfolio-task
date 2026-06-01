"use client";

import { useRef, useState } from "react";

import type { UseFormReturn } from "react-hook-form";

import { Textarea } from "@/components/ui/textarea";

import { cn } from "@/lib/utils";

import {
  DESCRIPTION_MAX_LENGTH,
  type TaskFormValues,
} from "../schemas/task-schema";

type Props = {
  form: UseFormReturn<TaskFormValues>;
};

export function TaskDescriptionField({ form }: Props) {
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  const [descriptionLength, setDescriptionLength] = useState(
    form.getValues("description")?.length ?? 0,
  );

  const { ref: descriptionFieldRef, ...descriptionField } =
    form.register("description");

  const isDescriptionTooLong = descriptionLength > DESCRIPTION_MAX_LENGTH;

  const adjustHeight = (target: HTMLTextAreaElement) => {
    target.style.height = "auto";
    target.style.height = `${target.scrollHeight}px`;
  };

  return (
    <div className="min-w-0 space-y-2">
      <Textarea
        placeholder="Add details about this task..."
        {...descriptionField}
        ref={(el) => {
          descriptionFieldRef(el);
          textareaRef.current = el;
        }}
        rows={12}
        onInput={(e) => {
          const target = e.currentTarget;

          adjustHeight(target);

          setDescriptionLength(target.value.length);
        }}
        className={cn(
          `
            min-h-[40vh]
            w-full
            resize-none
            overflow-hidden
            rounded-2xl
            px-5
            py-4
            text-base
            leading-7
            md:text-base
          `,
          isDescriptionTooLong &&
            `
              border-red-500
              bg-red-950/20
              text-red-100
              focus-visible:ring-red-500
            `,
        )}
      />

      <div className="flex items-center justify-between">
        <div>
          {form.formState.errors.description && (
            <p className="text-sm text-red-400">
              {form.formState.errors.description.message}
            </p>
          )}
        </div>

        <p
          className={cn(
            "text-xs",
            isDescriptionTooLong
              ? "font-semibold text-red-400"
              : "text-muted-foreground",
          )}
        >
          {descriptionLength} / {DESCRIPTION_MAX_LENGTH}
        </p>
      </div>
    </div>
  );
}
