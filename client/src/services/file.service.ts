import { QUERY_FILES_KEY, QUERY_FILE_KEY } from "@/constant/query.constant";
import { apiClient } from "@/lib/httpCommon";
import { FileEntity, fileEntitySchema } from "@/types/file.type";
import { UseQueryOptions, useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";

const getFiles = async () => {
  const res = await apiClient.get("/files");
  return fileEntitySchema.parse(res.data);
};

export const useGetFiles = (
  options?: UseQueryOptions<
    FileEntity,
    AxiosError,
    FileEntity,
    readonly [string]
  >
) => {
  return useQuery({
    queryKey: [QUERY_FILES_KEY],
    queryFn: getFiles,
    ...options,
  });
};

export const useQueryFile = (filePath: string) => {
  const getCustomer = async (customerId: string) => {
    const res = await apiClient.get(`/files/${customerId}`);
    return fileEntitySchema.parse(res.data);
  };
  return useQuery<FileEntity>([QUERY_FILE_KEY, filePath], () =>
    getCustomer(filePath)
  );
};
