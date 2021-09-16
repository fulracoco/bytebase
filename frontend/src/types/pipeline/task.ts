import { TaskCheckRunId } from "..";
import { Database } from "../database";
import {
  BackupId,
  DatabaseId,
  InstanceId,
  ProjectId,
  TaskId,
  TaskRunId,
} from "../id";
import { Instance, MigrationType } from "../instance";
import { Principal } from "../principal";
import { VCSPushEvent } from "../vcs";
import { Pipeline } from "./pipeline";
import { Stage } from "./stage";

export type TaskType =
  | "bb.task.general"
  | "bb.task.database.create"
  | "bb.task.database.schema.update"
  | "bb.task.database.restore";

export type TaskStatus =
  | "PENDING"
  | "PENDING_APPROVAL"
  | "RUNNING"
  | "DONE"
  | "FAILED"
  | "CANCELED";

export type TaskGeneralPayload = {
  statement: string;
};

export type TaskDatabaseCreatePayload = {
  projectId: ProjectId;
  statement: string;
  databaseName: string;
  characterSet: string;
  collation: string;
};

export type TaskDatabaseSchemaUpdatePayload = {
  migrationType: MigrationType;
  statement: string;
  rollbackStatement: string;
  pushEvent?: VCSPushEvent;
};

export type TaskDatabaseRestorePayload = {
  databaseName: string;
  backupId: BackupId;
};

export type TaskPayload =
  | TaskGeneralPayload
  | TaskDatabaseCreatePayload
  | TaskDatabaseSchemaUpdatePayload
  | TaskDatabaseRestorePayload;

export type Task = {
  id: TaskId;

  // Related fields
  taskRunList: TaskRun[];
  taskCheckRunList: TaskCheckRun[];
  pipeline: Pipeline;
  stage: Stage;

  // Standard fields
  creator: Principal;
  createdTs: number;
  updater: Principal;
  updatedTs: number;

  // Domain specific fields
  name: string;
  status: TaskStatus;
  type: TaskType;
  instance: Instance;
  // Tasks like creating database may not have database.
  database?: Database;
  payload?: TaskPayload;
};

export type TaskCreate = {
  // Domain specific fields
  name: string;
  status: TaskStatus;
  type: TaskType;
  instanceId: InstanceId;
  databaseId?: DatabaseId;
  statement: string;
  rollbackStatement: string;
  databaseName?: string;
  characterSet?: string;
  collation?: string;
  backupId?: BackupId;
  migrationType?: MigrationType;
};

export type TaskStatusPatch = {
  // Domain specific fields
  status: TaskStatus;
  comment?: string;
};

// TaskRun is one run of a particular task
export type TaskRunStatus = "RUNNING" | "DONE" | "FAILED" | "CANCELED";

export type TaskRun = {
  id: TaskRunId;

  // Standard fields
  creator: Principal;
  createdTs: number;
  updater: Principal;
  updatedTs: number;

  // Domain specific fields
  name: string;
  status: TaskRunStatus;
  type: TaskType;
  comment: string;
  payload?: TaskPayload;
};

export type TaskCheckRunStatus = "RUNNING" | "DONE" | "FAILED" | "CANCELED";

export type TaskCheckType = "bb.task-check.database.statement.advise";

export type TaskCheckDatabaseStatementAdvisePayload = {
  statement: string;
};

export type TaskCheckStatus = "SUCCESS" | "WARN" | "ERROR";

export type TaskCheckResult = {
  status: TaskCheckStatus;
  title: string;
  content: string;
};

export type TaskCheckRunResultPayload = {
  resultList: TaskCheckResult[];
};

export type TaskCheckRun = {
  id: TaskCheckRunId;

  // Standard fields
  creator: Principal;
  createdTs: number;
  updater: Principal;
  updatedTs: number;

  // Domain specific fields
  taskId: TaskId;
  name: string;
  status: TaskCheckRunStatus;
  type: TaskCheckType;
  comment: string;
  result: TaskCheckRunResultPayload;
  payload?: TaskPayload;
};
