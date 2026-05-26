import { Badge } from "@/components/ui/badge";

import type { TaskStatus } from "../schemas/task-schema";

type Props = {
  status: TaskStatus;
};

export function TaskStatusBadge({ status }: Props) {
  return <Badge variant="outline">{status}</Badge>;
}
