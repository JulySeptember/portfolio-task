// src/app/(protected)/tasks/[id]/page.tsx

"use client";

import { useEffect, useState } from "react";

import Link from "next/link";

import { useParams, useRouter } from "next/navigation";

import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";

import { apiClient } from "@/lib/api/client";

import { taskSchema, type Task } from "@/features/tasks/schemas/task-schema";

import { TaskEditor } from "@/features/tasks/components/task-editor";

export default function TaskDetailPage() {
  const params = useParams();

  const router = useRouter();

  const [task, setTask] = useState<Task | null>(null);

  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const publicIdRaw = params.id;

    // invalid
    if (!publicIdRaw || Array.isArray(publicIdRaw)) {
      router.replace("/tasks");

      return;
    }

    async function fetchTask() {
      try {
        // publicId で取得
        const data = await apiClient(`/api/v1/tasks/public/${publicIdRaw}`);

        setTask(taskSchema.parse(data));
      } catch (error) {
        console.error(error);

        router.replace("/tasks");
      } finally {
        setIsLoading(false);
      }
    }

    fetchTask();
  }, [params.id, router]);

  if (isLoading || !task) {
    return (
      <div
        className="
          flex
          min-h-screen
          items-center
          justify-center
        "
      >
        <p
          className="
            text-sm
            text-muted-foreground
          "
        >
          Loading...
        </p>
      </div>
    );
  }

  return (
    <div
      className="
        mx-auto
        w-full
        max-w-5xl
        px-4
        py-6
        sm:px-6
        lg:px-8
      "
    >
      <div className="space-y-10">
        {/* back */}

        <div className="space-y-5">
          <Button
            asChild
            variant="ghost"
            className="
              h-auto
              px-0
              text-sm
            "
          >
            <Link href="/tasks">
              <ArrowLeft
                className="
                  mr-2
                  h-4
                  w-4
                "
              />
              Back to Tasks
            </Link>
          </Button>
        </div>

        {/* editor */}

        <TaskEditor mode="edit" task={task} autoResizeDescription />
      </div>
    </div>
  );
}
