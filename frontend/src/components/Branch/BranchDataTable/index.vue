<template>
  <NDataTable
    v-bind="$attrs"
    :loading="!ready"
    :columns="columns"
    :data="branches"
    :row-key="rowKey"
    :pagination="pagination"
    :paginate-single-page="false"
  />
</template>

<script lang="ts" setup>
import { NDataTable, DataTableColumn, PaginationProps } from "naive-ui";
import { computed, h, ref } from "vue";
import { useI18n } from "vue-i18n";
import BranchBaseline from "@/components/Branch/BranchBaseline.vue";
import { Branch } from "@/types/proto/v1/branch_service";
import AsyncBranchNameCell from "./cells/AsyncBranchNameCell.vue";
import BranchNameCell from "./cells/BranchNameCell.vue";
import BranchUpdatedCell from "./cells/BranchUpdatedCell.vue";

defineOptions({
  name: "BranchDataTable",
});

const props = defineProps<{
  branches: Branch[];
  ready?: boolean;
  showParentBranchColumn?: boolean;
}>();

const { t } = useI18n();

const pagination = ref<PaginationProps>({
  pageSize: 5,
});

const columns = computed(() => {
  const columns: (DataTableColumn<Branch> & { hide?: boolean })[] = [
    {
      title: t("common.branch"),
      key: "branchName",
      render: (row: Branch) => {
        return h(BranchNameCell, {
          branch: row,
        });
      },
    },
    {
      title: t("schema-designer.parent-branch"),
      key: "branchName",
      width: 256,
      hide: !props.showParentBranchColumn,
      render: (row: Branch) => {
        if (!row.parentBranch) return null;
        return h(AsyncBranchNameCell, {
          name: row.parentBranch,
        });
      },
    },
    {
      title: t("schema-designer.baseline-version"),
      key: "baselineVersion",
      width: 300,
      render: (row: Branch) => {
        return h(BranchBaseline, {
          branch: row,
          showInstanceIcon: true,
        });
      },
    },
    {
      title: t("common.updated"),
      key: "updated",
      width: 200,
      render: (row: Branch) => {
        return h(BranchUpdatedCell, {
          branch: row,
        });
      },
    },
  ];

  return columns.filter((column) => !column.hide);
});

const rowKey = (row: Branch) => {
  return row.name;
};
</script>
